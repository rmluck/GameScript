// NFL standings logic

package standings

import (
	"fmt"
	"sort"

	"gamescript/internal/database"
)


type NFLTeamRecord struct {
	TeamID              int
	TeamCity            string
	TeamName            string
	TeamAbbr            string
	Conference          string
	Division            string
	Wins                int
	Losses              int
	Ties                int
	HomeWins            int
	HomeLosses          int
	HomeTies            int
	AwayWins            int
	AwayLosses          int
	AwayTies            int
	DivisionWins        int
	DivisionLosses      int
	DivisionTies        int
	ConferenceWins      int
	ConferenceLosses    int
	ConferenceTies      int
	PointsFor           int
	PointsAgainst       int
	WinPct              float64
	ConferenceGamesBack float64
	DivisionGamesBack   float64
	StrengthOfSchedule  float64
	StrengthOfVictory   float64
	LogoURL             string
	TeamPrimaryColor    string
	TeamSecondaryColor  string
}

type NFLStandings struct {
	AFC        NFLConferenceStandings
	NFC        NFLConferenceStandings
	DraftOrder []NFLDraftPick
}

type NFLConferenceStandings struct {
	Divisions    map[string][]NFLTeamRecord // Keyed by division name
	PlayoffSeeds []NFLPlayoffSeed
}

type NFLPlayoffSeed struct {
	Seed             int
	Team             NFLTeamRecord
	IsDivisionWinner bool
}

type NFLDraftPick struct {
	Pick int
	Team NFLTeamRecord
}

type NFLGameResult struct {
	GameID     int
	HomeTeamID int
	AwayTeamID int
	HomeScore  int
	AwayScore  int
	Week       int
}

func CalculateNFLStandings(db *database.DB, scenarioID int, seasonID int) (*NFLStandings, error) {
	// Get all teams for the season
	teams, err := getNFLTeams(db, seasonID)
	if err != nil {
		return nil, fmt.Errorf("error getting teams: %w", err)
	}

	// Get all game results for the scenario
	games, err := getNFLGameResults(db, scenarioID, seasonID)
	if err != nil {
		return nil, fmt.Errorf("error getting game results: %w", err)
	}

	// Calculate team records
	records := calculateNFLTeamRecords(teams, games)

	// Calculate strength metrics
	calculateNFLStrengthMetrics(records, games)

	// Separate by conference
	afcTeams := filterByNFLConference(records, "AFC")
	nfcTeams := filterByNFLConference(records, "NFC")

	// Calculate playoff seeds for each conference
	afcStandings := calculateNFLConferenceStandings(afcTeams, games)
	nfcStandings := calculateNFLConferenceStandings(nfcTeams, games)

	// Calculate draft order
	draftOrder := calculateNFLDraftOrder(records, afcStandings, nfcStandings)

	return &NFLStandings{
		AFC:        afcStandings,
		NFC:        nfcStandings,
		DraftOrder: draftOrder,
	}, nil
}

func getNFLTeams(db *database.DB, seasonID int) ([]NFLTeamRecord, error) {
	query := `
		SELECT
			id, city, name, abbreviation, conference, division, logo_url, primary_color, secondary_color
		FROM teams
		WHERE season_id = $1
		ORDER BY conference, division, name
	`

	rows, err := db.Query(query, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []NFLTeamRecord
	for rows.Next() {
		var team NFLTeamRecord
		err := rows.Scan(
			&team.TeamID,
			&team.TeamCity,
			&team.TeamName,
			&team.TeamAbbr,
			&team.Conference,
			&team.Division,
			&team.LogoURL,
			&team.TeamPrimaryColor,
			&team.TeamSecondaryColor,
		)
		if err != nil {
			continue
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func getNFLGameResults(db *database.DB, scenarioID int, seasonID int) ([]NFLGameResult, error) {
	query := `
        SELECT
            game.id, game.home_team_id, game.away_team_id, game.week,
            game.home_score AS actual_home_score,
            game.away_score AS actual_away_score,
            game.status,
            pick.picked_team_id,
            pick.predicted_home_score,
            pick.predicted_away_score
        FROM games game
        LEFT JOIN picks pick ON game.id = pick.game_id AND pick.scenario_id = $1
        WHERE game.season_id = $2
        ORDER BY game.week, game.start_time
    `

	rows, err := db.Query(query, scenarioID, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []NFLGameResult
	for rows.Next() {
		var game NFLGameResult
		var actualHomeScore, actualAwayScore *int
		var status string
		var pickedTeamID, predictedHomeScore, predictedAwayScore *int

		err := rows.Scan(
			&game.GameID,
			&game.HomeTeamID,
			&game.AwayTeamID,
			&game.Week,
			&actualHomeScore,
			&actualAwayScore,
			&status,
			&pickedTeamID,
			&predictedHomeScore,
			&predictedAwayScore,
		)
		if err != nil {
			continue
		}

		// Priority 1: User has made a pick
		if pickedTeamID != nil {
			// If user provided predicted scores, use those
			if predictedHomeScore != nil && predictedAwayScore != nil {
				game.HomeScore = *predictedHomeScore
				game.AwayScore = *predictedAwayScore
			} else {
				// No predicted scores, use dummy scores based on picked winner
				if *pickedTeamID == game.HomeTeamID {
					game.HomeScore = 1
					game.AwayScore = 0
				} else if *pickedTeamID == game.AwayTeamID {
					game.HomeScore = 0
					game.AwayScore = 1
				} else if *pickedTeamID == 0 {
					// Tie picked
					game.HomeScore = 0
					game.AwayScore = 0
				} else {
					// Invalide picked team ID
					continue
				}
			}
			games = append(games, game)
			continue
		}

		// Priority 2: No pick, but game is final
		if status == "final" && actualHomeScore != nil && actualAwayScore != nil {
			game.HomeScore = *actualHomeScore
			game.AwayScore = *actualAwayScore
			games = append(games, game)
			continue
		}

		// Priority 3: No pick and game not final - skip this game
		// (don't add to games slice)
	}

	return games, nil
}

func calculateNFLTeamRecords(teams []NFLTeamRecord, games []NFLGameResult) []NFLTeamRecord {
	// Initialize record map for easy lookup
	recordMap := make(map[int]*NFLTeamRecord)
	for i := range teams {
		recordMap[teams[i].TeamID] = &teams[i]
	}

	// Process each game to update team records
	for _, game := range games {
		homeTeam := recordMap[game.HomeTeamID]
		awayTeam := recordMap[game.AwayTeamID]

		if homeTeam == nil || awayTeam == nil {
			continue
		}

		isDivisionGame := homeTeam.Division == awayTeam.Division
		isConferenceGame := homeTeam.Conference == awayTeam.Conference

		// Determine winner
		if game.HomeScore > game.AwayScore {
			// Home team wins
			homeTeam.Wins++
			homeTeam.HomeWins++
			awayTeam.Losses++
			awayTeam.AwayLosses++

			if isDivisionGame {
				homeTeam.DivisionWins++
				awayTeam.DivisionLosses++
			}

			if isConferenceGame {
				homeTeam.ConferenceWins++
				awayTeam.ConferenceLosses++
			}
		} else if game.AwayScore > game.HomeScore {
			// Away team wins
			awayTeam.Wins++
			awayTeam.AwayWins++
			homeTeam.Losses++
			homeTeam.HomeLosses++

			if isDivisionGame {
				awayTeam.DivisionWins++
				homeTeam.DivisionLosses++
			}

			if isConferenceGame {
				awayTeam.ConferenceWins++
				homeTeam.ConferenceLosses++
			}
		} else {
			// Tie
			homeTeam.Ties++
			homeTeam.HomeTies++
			awayTeam.Ties++
			awayTeam.AwayTies++

			if isDivisionGame {
				homeTeam.DivisionTies++
				awayTeam.DivisionTies++
			}

			if isConferenceGame {
				homeTeam.ConferenceTies++
				awayTeam.ConferenceTies++
			}
		}

		// Update points for/against
		homeTeam.PointsFor += game.HomeScore
		homeTeam.PointsAgainst += game.AwayScore
		awayTeam.PointsFor += game.AwayScore
		awayTeam.PointsAgainst += game.HomeScore
	}

	// Calculate win percentages
	for _, team := range recordMap {
		team.WinPct = calculateNFLWinPct(team.Wins, team.Losses, team.Ties)
	}

	return teams
}

func calculateNFLStrengthMetrics(teams []NFLTeamRecord, games []NFLGameResult) {
	// Initialize team map for easy lookup
	teamMap := make(map[int]*NFLTeamRecord)
	for i := range teams {
		teamMap[teams[i].TeamID] = &teams[i]
	}

	// Calculate strength of schedule and strength of victory
	for i := range teams {
		team := &teams[i]

		var opponentTotalWins, opponentTotalLosses, opponentTotalTies int
		var opponentCount int
		var defeatedOpponentWins, defeatedOpponentLosses, defeatedOpponentTies int
		var defeatedCount int

		for _, game := range games {
			var opponentID int
			var teamWon bool

			if game.HomeTeamID == team.TeamID {
				opponentID = game.AwayTeamID
				teamWon = game.HomeScore > game.AwayScore
			} else if game.AwayTeamID == team.TeamID {
				opponentID = game.HomeTeamID
				teamWon = game.AwayScore > game.HomeScore
			} else {
				continue
			}

			opponent := teamMap[opponentID]
			if opponent == nil {
				continue
			}

			// Strength of schedule: all opponents
			opponentTotalWins += opponent.Wins
			opponentTotalLosses += opponent.Losses
			opponentTotalTies += opponent.Ties
			opponentCount++

			// Strength of victory: only defeated opponents
			if teamWon {
				defeatedOpponentWins += opponent.Wins
				defeatedOpponentLosses += opponent.Losses
				defeatedOpponentTies += opponent.Ties
				defeatedCount++
			}
		}

		// Calculate averages
		if opponentCount > 0 {
			team.StrengthOfSchedule = calculateNFLWinPct(opponentTotalWins, opponentTotalLosses, opponentTotalTies)
		}
		if defeatedCount > 0 {
			team.StrengthOfVictory = calculateNFLWinPct(defeatedOpponentWins, defeatedOpponentLosses, defeatedOpponentTies)
		}
	}
}

func filterByNFLConference(teams []NFLTeamRecord, conference string) []NFLTeamRecord {
	var filtered []NFLTeamRecord
	for _, team := range teams {
		if team.Conference == conference {
			filtered = append(filtered, team)
		}
	}
	return filtered
}

func calculateNFLConferenceStandings(teams []NFLTeamRecord, games []NFLGameResult) NFLConferenceStandings {
	// Group teams by division
	divisions := make(map[string][]NFLTeamRecord)
	for _, team := range teams {
		divisions[team.Division] = append(divisions[team.Division], team)
	}

	// Determine division winners
	divisionWinners := []NFLTeamRecord{}
	nonWinners := []NFLTeamRecord{}
	for divName, divTeams := range divisions {
		// Sort division teams with tiebreakers
		sortedDiv := applyNFLDivisionTiebreakers(divTeams, games)

		// Calculate division games back
		divLeader := sortedDiv[0]
		for i := range sortedDiv {
			sortedDiv[i].DivisionGamesBack = calculateNFLGamesBack(divLeader, sortedDiv[i])
		}

		// First place is division winner
		divisionWinners = append(divisionWinners, sortedDiv[0])

		// Rest are wild card candidates
		nonWinners = append(nonWinners, sortedDiv[1:]...)

		// Update divisions map with sorted teams
		divisions[divName] = sortedDiv
	}

	// Rank division winners (seeds 1-4)
	divisionWinners = applyNFLConferenceTiebreakers(divisionWinners, games, true)

	// Rank non-division winners (seeds 5-16)
	nonWinners = applyNFLConferenceTiebreakers(nonWinners, games, false)

	// Create playoff seeds
	playoffSeeds := []NFLPlayoffSeed{}
	for i, team := range divisionWinners {
		playoffSeeds = append(playoffSeeds, NFLPlayoffSeed{
			Seed:             i + 1,
			Team:             team,
			IsDivisionWinner: true,
		})
	}
	for i, team := range nonWinners {
		playoffSeeds = append(playoffSeeds, NFLPlayoffSeed{
			Seed:             i + 5,
			Team:             team,
			IsDivisionWinner: false,
		})
	}

	// Calculate conference games back
	conferenceLeader := playoffSeeds[0].Team
	for i := range playoffSeeds {
		playoffSeeds[i].Team.ConferenceGamesBack = calculateNFLGamesBack(conferenceLeader, playoffSeeds[i].Team)
	}

	return NFLConferenceStandings{
		Divisions:    divisions,
		PlayoffSeeds: playoffSeeds,
	}
}

func applyNFLDivisionTiebreakers(teams []NFLTeamRecord, games []NFLGameResult) []NFLTeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// Group by win percentage (rounded to avoid floating point issues)
	pctGroups := make(map[string][]NFLTeamRecord)
	for _, team := range teams {
		key := fmt.Sprintf("%.6f", team.WinPct)
		pctGroups[key] = append(pctGroups[key], team)
	}

	var result []NFLTeamRecord
	for _, group := range pctGroups {
		if len(group) == 1 {
			result = append(result, group[0])
		} else if len(group) == 2 {
			result = append(result, resolveNFLTwoTeamDivisionTie(group, games)...)
		} else {
			result = append(result, resolveNFLMultiTeamDivisionTie(group, games)...)
		}
	}

	// Sort by win percentage
	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct
	})

	return result
}

func resolveNFLTwoTeamDivisionTie(teams []NFLTeamRecord, games []NFLGameResult) []NFLTeamRecord {
	a, b := teams[0], teams[1]

	// Step 1: Head-to-head
	h2h := compareNFLHeadToHead(teams, games)
	if len(h2h) == 1 {
		if h2h[0].TeamID == a.TeamID {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 2: Division win percentage
	aDivPct := calculateNFLWinPct(a.DivisionWins, a.DivisionLosses, a.DivisionTies)
	bDivPct := calculateNFLWinPct(b.DivisionWins, b.DivisionLosses, b.DivisionTies)
	if aDivPct != bDivPct {
		if aDivPct > bDivPct {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 3: Common games
	commonResults := compareNFLCommonGames(teams, games, 0)
	if len(commonResults) == 1 {
		if commonResults[0].TeamID == a.TeamID {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 4: Conference win percentage
	aConfPct := calculateNFLWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
	bConfPct := calculateNFLWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
	if aConfPct != bConfPct {
		if aConfPct > bConfPct {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 5: Strength of victory
	if a.StrengthOfVictory != b.StrengthOfVictory {
		if a.StrengthOfVictory > b.StrengthOfVictory {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 6: Strength of schedule
	if a.StrengthOfSchedule != b.StrengthOfSchedule {
		if a.StrengthOfSchedule > b.StrengthOfSchedule {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 7: Point differential
	aDiff := a.PointsFor - a.PointsAgainst
	bDiff := b.PointsFor - b.PointsAgainst
	if aDiff != bDiff {
		if aDiff > bDiff {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 8: Points scored
	if a.PointsFor != b.PointsFor {
		if a.PointsFor > b.PointsFor {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 9: Points allowed (fewer is better)
	if a.PointsAgainst != b.PointsAgainst {
		if a.PointsAgainst < b.PointsAgainst {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Random drawing - use TeamID for consistency
	if a.TeamID < b.TeamID {
		return []NFLTeamRecord{a, b}
	}
	return []NFLTeamRecord{b, a}
}

func resolveNFLMultiTeamDivisionTie(teams []NFLTeamRecord, games []NFLGameResult) []NFLTeamRecord {
	// Step 1: Head-to-head (best win pct in games among tied teams)
	h2hWinner := findNFLHeadToHeadWinner(teams, games)
	if h2hWinner != nil {
		remaining := removeNFLTeam(teams, h2hWinner.TeamID)
		result := []NFLTeamRecord{*h2hWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 2: Division win percentage
	divWinner := findBestNFLDivisionRecord(teams)
	if divWinner != nil {
		remaining := removeNFLTeam(teams, divWinner.TeamID)
		result := []NFLTeamRecord{*divWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 3: Common games
	commonWinner := findBestNFLCommonGamesRecord(teams, games, 0)
	if commonWinner != nil {
		remaining := removeNFLTeam(teams, commonWinner.TeamID)
		result := []NFLTeamRecord{*commonWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 4: Conference win percentage
	confWinner := findBestNFLConferenceRecord(teams)
	if confWinner != nil {
		remaining := removeNFLTeam(teams, confWinner.TeamID)
		result := []NFLTeamRecord{*confWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 5: Strength of victory
	sovWinner := findBestNFLStrengthOfVictory(teams)
	if sovWinner != nil {
		remaining := removeNFLTeam(teams, sovWinner.TeamID)
		result := []NFLTeamRecord{*sovWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 6: Strength of schedule
	sosWinner := findBestNFLStrengthOfSchedule(teams)
	if sosWinner != nil {
		remaining := removeNFLTeam(teams, sosWinner.TeamID)
		result := []NFLTeamRecord{*sosWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 7: Point differential
	pdWinner := findBestNFLPointDifferential(teams)
	if pdWinner != nil {
		remaining := removeNFLTeam(teams, pdWinner.TeamID)
		result := []NFLTeamRecord{*pdWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 8: Points scored
	psWinner := findBestNFLPointsScored(teams)
	if psWinner != nil {
		remaining := removeNFLTeam(teams, psWinner.TeamID)
		result := []NFLTeamRecord{*psWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 9: Points allowed (fewer is better)
	paWinner := findBestNFLPointsAllowed(teams)
	if paWinner != nil {
		remaining := removeNFLTeam(teams, paWinner.TeamID)
		result := []NFLTeamRecord{*paWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Random drawing - use TeamID for consistency
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].TeamID < teams[j].TeamID
	})
	return teams
}

func applyNFLConferenceTiebreakers(teams []NFLTeamRecord, games []NFLGameResult, areDivisionWinners bool) []NFLTeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// Group by win percentage
	pctGroups := make(map[string][]NFLTeamRecord)
	for _, team := range teams {
		key := fmt.Sprintf("%.6f", team.WinPct)
		pctGroups[key] = append(pctGroups[key], team)
	}

	var result []NFLTeamRecord
	for _, group := range pctGroups {
		if len(group) == 1 {
			result = append(result, group[0])
		} else if len(group) == 2 {
			resolved := resolveNFLTwoTeamConferenceTie(group, games)
			result = append(result, resolved...)
		} else {
			resolved := resolveNFLMultiTeamConferenceTie(group, games)
			result = append(result, resolved...)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct
	})

	return result
}

func resolveNFLTwoTeamConferenceTie(teams []NFLTeamRecord, games []NFLGameResult) []NFLTeamRecord {
	a, b := teams[0], teams[1]

	// Step 1: Division winner if from same division
	if a.Division == b.Division {
		return resolveNFLTwoTeamDivisionTie(teams, games)
	}

	// Step 2: Head-to-head
	h2h := compareNFLHeadToHead(teams, games)
	if len(h2h) == 1 {
		if h2h[0].TeamID == a.TeamID {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 3: Conference win percentage
	aConfPct := calculateNFLWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
	bConfPct := calculateNFLWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
	if aConfPct != bConfPct {
		if aConfPct > bConfPct {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 4: Common games (minimum of 4)
	commonResults := compareNFLCommonGames(teams, games, 4)
	if len(commonResults) == 1 {
		if commonResults[0].TeamID == a.TeamID {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 5: Strength of victory
	if a.StrengthOfVictory != b.StrengthOfVictory {
		if a.StrengthOfVictory > b.StrengthOfVictory {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 6: Strength of schedule
	if a.StrengthOfSchedule != b.StrengthOfSchedule {
		if a.StrengthOfSchedule > b.StrengthOfSchedule {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 7: Point differential
	aDiff := a.PointsFor - a.PointsAgainst
	bDiff := b.PointsFor - b.PointsAgainst
	if aDiff != bDiff {
		if aDiff > bDiff {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 8: Points scored
	if a.PointsFor != b.PointsFor {
		if a.PointsFor > b.PointsFor {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 9: Points allowed (fewer is better)
	if a.PointsAgainst != b.PointsAgainst {
		if a.PointsAgainst < b.PointsAgainst {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Random drawing - use TeamID for consistency
	if a.TeamID < b.TeamID {
		return []NFLTeamRecord{a, b}
	}
	return []NFLTeamRecord{b, a}
}

func resolveNFLMultiTeamConferenceTie(teams []NFLTeamRecord, games []NFLGameResult) []NFLTeamRecord {
	if len(teams) == 1 {
		return teams
	}
	if len(teams) == 2 {
		return resolveNFLTwoTeamConferenceTie(teams, games)
	}

	// Step 1: Division winner if all from same division
	allSameDivision := true
	firstDivision := teams[0].Division
	for _, team := range teams[1:] {
		if team.Division != firstDivision {
			allSameDivision = false
			break
		}
	}
	if allSameDivision {
		return resolveNFLMultiTeamDivisionTie(teams, games)
	}

	// Step 2: Apply division tiebreaker to get best from each division
	divGroups := make(map[string][]NFLTeamRecord)

	for _, team := range teams {
		divGroups[team.Division] = append(divGroups[team.Division], team)
	}

	var filtered []NFLTeamRecord
	for _, divTeams := range divGroups {
		if len(divTeams) == 1 {
			filtered = append(filtered, divTeams[0])
		} else {
			sorted := applyNFLDivisionTiebreakers(divTeams, games)
			filtered = append(filtered, sorted[0])
		}
	}

	if len(filtered) == 1 {
		winner := filtered[0]
		remaining := removeNFLTeam(teams, winner.TeamID)
		result := []NFLTeamRecord{winner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 3: Head-to-head sweep
	sweepWinner := checkNFLHeadToHeadSweep(filtered, games)
	if sweepWinner != nil {
		remaining := removeNFLTeam(teams, sweepWinner.TeamID)
		result := []NFLTeamRecord{*sweepWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 4: Conference win percentage
	confWinner := findBestNFLConferenceRecord(filtered)
	if confWinner != nil {
		remaining := removeNFLTeam(teams, confWinner.TeamID)
		result := []NFLTeamRecord{*confWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 5: Common games (minimum of 4)
	commonWinner := findBestNFLCommonGamesRecord(filtered, games, 4)
	if commonWinner != nil {
		remaining := removeNFLTeam(teams, commonWinner.TeamID)
		result := []NFLTeamRecord{*commonWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 6: Strength of victory
	if sovWinner := findBestNFLStrengthOfVictory(filtered); sovWinner != nil {
		remaining := removeNFLTeam(teams, sovWinner.TeamID)
		result := []NFLTeamRecord{*sovWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 7: Strength of schedule
	if sosWinner := findBestNFLStrengthOfSchedule(filtered); sosWinner != nil {
		remaining := removeNFLTeam(teams, sosWinner.TeamID)
		result := []NFLTeamRecord{*sosWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 8: Point differential
	if pdWinner := findBestNFLPointDifferential(filtered); pdWinner != nil {
		remaining := removeNFLTeam(teams, pdWinner.TeamID)
		result := []NFLTeamRecord{*pdWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 9: Points scored
	if psWinner := findBestNFLPointsScored(filtered); psWinner != nil {
		remaining := removeNFLTeam(teams, psWinner.TeamID)
		result := []NFLTeamRecord{*psWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Step 10: Points allowed (fewer is better)
	if paWinner := findBestNFLPointsAllowed(filtered); paWinner != nil {
		remaining := removeNFLTeam(teams, paWinner.TeamID)
		result := []NFLTeamRecord{*paWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Random drawing - use TeamID for consistency
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TeamID < filtered[j].TeamID
	})

	winner := filtered[0]
	remaining := removeNFLTeam(teams, winner.TeamID)
	result := []NFLTeamRecord{winner}
	if len(remaining) > 0 {
		result = append(result, resolveNFLMultiTeamConferenceTie(remaining, games)...)
	}
	return result
}

func findNFLHeadToHeadWinner(teams []NFLTeamRecord, games []NFLGameResult) *NFLTeamRecord {
	type h2hRecord struct {
		teamID int
		wins   int
		losses int
		ties   int
	}

	records := make(map[int]*h2hRecord)
	teamIDs := make(map[int]bool)
	for _, team := range teams {
		records[team.TeamID] = &h2hRecord{teamID: team.TeamID}
		teamIDs[team.TeamID] = true
	}

	for _, game := range games {
		if !teamIDs[game.HomeTeamID] || !teamIDs[game.AwayTeamID] {
			continue
		}

		if game.HomeScore > game.AwayScore {
			records[game.HomeTeamID].wins++
			records[game.AwayTeamID].losses++
		} else if game.AwayScore > game.HomeScore {
			records[game.AwayTeamID].wins++
			records[game.HomeTeamID].losses++
		} else {
			records[game.HomeTeamID].ties++
			records[game.AwayTeamID].ties++
		}
	}

	var bestPct float64 = -2.0
	var bestTeamID int
	var tie bool

	for _, rec := range records {
		pct := calculateNFLWinPct(rec.wins, rec.losses, rec.ties)
		if pct > bestPct {
			bestPct = pct
			bestTeamID = rec.teamID
			tie = false
		} else if pct == bestPct {
			tie = true
		}
	}

	if tie {
		return nil
	}

	for i := range teams {
		if teams[i].TeamID == bestTeamID {
			return &teams[i]
		}
	}
	return nil
}

func findBestNFLDivisionRecord(teams []NFLTeamRecord) *NFLTeamRecord {
	var bestPct float64 = -2.0
	var bestTeam *NFLTeamRecord
	var tie bool

	for i := range teams {
		pct := calculateNFLWinPct(teams[i].DivisionWins, teams[i].DivisionLosses, teams[i].DivisionTies)
		if pct > bestPct {
			bestPct = pct
			bestTeam = &teams[i]
			tie = false
		} else if pct == bestPct {
			tie = true
		}
	}

	if tie {
		return nil
	}
	return bestTeam
}

func findBestNFLConferenceRecord(teams []NFLTeamRecord) *NFLTeamRecord {
	var bestPct float64 = -2.0
	var bestTeam *NFLTeamRecord
	var tie bool

	for i := range teams {
		pct := calculateNFLWinPct(teams[i].ConferenceWins, teams[i].ConferenceLosses, teams[i].ConferenceTies)
		if pct > bestPct {
			bestPct = pct
			bestTeam = &teams[i]
			tie = false
		} else if pct == bestPct {
			tie = true
		}
	}

	if tie {
		return nil
	}
	return bestTeam
}

func findBestNFLCommonGamesRecord(teams []NFLTeamRecord, games []NFLGameResult, minGames int) *NFLTeamRecord {
	teamIDs := make(map[int]bool)
	for _, team := range teams {
		teamIDs[team.TeamID] = true
	}

	// Find common opponents
	opponentCounts := make(map[int]int)
	for _, team := range teams {
		for _, game := range games {
			var oppID int
			if game.HomeTeamID == team.TeamID && !teamIDs[game.AwayTeamID] {
				oppID = game.AwayTeamID
			} else if game.AwayTeamID == team.TeamID && !teamIDs[game.HomeTeamID] {
				oppID = game.HomeTeamID
			} else {
				continue
			}
			opponentCounts[oppID]++
		}
	}

	var commonOpponents []int
	for oppID, count := range opponentCounts {
		if count == len(teams) {
			commonOpponents = append(commonOpponents, oppID)
		}
	}

	if len(commonOpponents) < minGames {
		return nil
	}

	type record struct {
		teamID int
		wins   int
		losses int
		ties   int
	}

	records := make(map[int]*record)
	for _, team := range teams {
		records[team.TeamID] = &record{teamID: team.TeamID}
	}

	for _, game := range games {
		var teamID, oppID int
		var isCommon bool

		if teamIDs[game.HomeTeamID] {
			teamID = game.HomeTeamID
			oppID = game.AwayTeamID
		} else if teamIDs[game.AwayTeamID] {
			teamID = game.AwayTeamID
			oppID = game.HomeTeamID
		} else {
			continue
		}

		for _, commonOpp := range commonOpponents {
			if oppID == commonOpp {
				isCommon = true
				break
			}
		}

		if !isCommon {
			continue
		}

		rec := records[teamID]
		if game.HomeTeamID == teamID {
			if game.HomeScore > game.AwayScore {
				rec.wins++
			} else if game.AwayScore > game.HomeScore {
				rec.losses++
			} else {
				rec.ties++
			}
		} else {
			if game.AwayScore > game.HomeScore {
				rec.wins++
			} else if game.HomeScore > game.AwayScore {
				rec.losses++
			} else {
				rec.ties++
			}
		}
	}

	var bestPct float64 = -2.0
	var bestTeamID int
	var tie bool

	for _, rec := range records {
		pct := calculateNFLWinPct(rec.wins, rec.losses, rec.ties)
		if pct > bestPct {
			bestPct = pct
			bestTeamID = rec.teamID
			tie = false
		} else if pct == bestPct {
			tie = true
		}
	}

	if tie {
		return nil
	}

	for i := range teams {
		if teams[i].TeamID == bestTeamID {
			return &teams[i]
		}
	}
	return nil
}

func findBestNFLStrengthOfVictory(teams []NFLTeamRecord) *NFLTeamRecord {
	var best float64 = -2.0
	var bestTeam *NFLTeamRecord
	var tie bool

	for i := range teams {
		if teams[i].StrengthOfVictory > best {
			best = teams[i].StrengthOfVictory
			bestTeam = &teams[i]
			tie = false
		} else if teams[i].StrengthOfVictory == best {
			tie = true
		}
	}

	if tie {
		return nil
	}
	return bestTeam
}

func findBestNFLStrengthOfSchedule(teams []NFLTeamRecord) *NFLTeamRecord {
	var best float64 = -2.0
	var bestTeam *NFLTeamRecord
	var tie bool

	for i := range teams {
		if teams[i].StrengthOfSchedule > best {
			best = teams[i].StrengthOfSchedule
			bestTeam = &teams[i]
			tie = false
		} else if teams[i].StrengthOfSchedule == best {
			tie = true
		}
	}

	if tie {
		return nil
	}
	return bestTeam
}

func findBestNFLPointDifferential(teams []NFLTeamRecord) *NFLTeamRecord {
	var bestDiff int = -1 << 31
	var bestTeam *NFLTeamRecord
	var tie bool

	for i := range teams {
		diff := teams[i].PointsFor - teams[i].PointsAgainst
		if diff > bestDiff {
			bestDiff = diff
			bestTeam = &teams[i]
			tie = false
		} else if diff == bestDiff {
			tie = true
		}
	}

	if tie {
		return nil
	}
	return bestTeam
}

func findBestNFLPointsScored(teams []NFLTeamRecord) *NFLTeamRecord {
	var best int = -1
	var bestTeam *NFLTeamRecord
	var tie bool

	for i := range teams {
		if teams[i].PointsFor > best {
			best = teams[i].PointsFor
			bestTeam = &teams[i]
			tie = false
		} else if teams[i].PointsFor == best {
			tie = true
		}
	}

	if tie {
		return nil
	}
	return bestTeam
}

func findBestNFLPointsAllowed(teams []NFLTeamRecord) *NFLTeamRecord {
	var best int = 1 << 31
	var bestTeam *NFLTeamRecord
	var tie bool

	for i := range teams {
		if teams[i].PointsAgainst < best {
			best = teams[i].PointsAgainst
			bestTeam = &teams[i]
			tie = false
		} else if teams[i].PointsAgainst == best {
			tie = true
		}
	}

	if tie {
		return nil
	}
	return bestTeam
}

func resolveNFLMultiTeamTie(teams []NFLTeamRecord, games []NFLGameResult, twoTeamCompare func(NFLTeamRecord, NFLTeamRecord) int) []NFLTeamRecord {
	// Try head-to-head sweep first
	sweepWinner := checkNFLHeadToHeadSweep(teams, games)
	if sweepWinner != nil {
		remaining := []NFLTeamRecord{}
		for _, team := range teams {
			if team.TeamID != sweepWinner.TeamID {
				remaining = append(remaining, team)
			}
		}

		result := []NFLTeamRecord{*sweepWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNFLMultiTeamTie(remaining, games, twoTeamCompare)...)
		}
		return result
	}

	// Fall back to sorting with two-team comparisons
	sorted := make([]NFLTeamRecord, len(teams))
	copy(sorted, teams)

	sort.Slice(sorted, func(i, j int) bool {
		return twoTeamCompare(sorted[i], sorted[j]) < 0
	})

	return sorted
}

func compareNFLHeadToHead(teams []NFLTeamRecord, games []NFLGameResult) []NFLTeamRecord {
	if len(teams) != 2 {
		return teams
	}

	teamA, teamB := teams[0], teams[1]
	aWins, bWins := 0, 0

	// Find games between these two teams
	for _, game := range games {
		isMatchup := (game.HomeTeamID == teamA.TeamID && game.AwayTeamID == teamB.TeamID) || (game.HomeTeamID == teamB.TeamID && game.AwayTeamID == teamA.TeamID)

		if !isMatchup {
			continue
		}

		// Determine winner
		var winner int
		if game.HomeScore > game.AwayScore {
			winner = game.HomeTeamID
		} else if game.AwayScore > game.HomeScore {
			winner = game.AwayTeamID
		} else {
			continue
		}

		// Track wins
		if winner == teamA.TeamID {
			aWins++
		} else if winner == teamB.TeamID {
			bWins++
		}
	}

	// Return winner or 0 if tied/no games
	if aWins > bWins {
		return []NFLTeamRecord{teamA}
	} else if bWins > aWins {
		return []NFLTeamRecord{teamB}
	}

	return teams
}

func checkNFLHeadToHeadSweep(teams []NFLTeamRecord, games []NFLGameResult) *NFLTeamRecord {
	for _, candidate := range teams {
		hasDefeatedAll := true

		for _, opponent := range teams {
			if candidate.TeamID == opponent.TeamID {
				continue
			}

			defeatedOpponent := false
			for _, game := range games {
				isMatchup := (game.HomeTeamID == candidate.TeamID && game.AwayTeamID == opponent.TeamID) || (game.HomeTeamID == opponent.TeamID && game.AwayTeamID == candidate.TeamID)

				if !isMatchup {
					continue
				}

				var winner int
				if game.HomeScore > game.AwayScore {
					winner = game.HomeTeamID
				} else if game.AwayScore > game.HomeScore {
					winner = game.AwayTeamID
				} else {
					continue
				}

				if winner == candidate.TeamID {
					defeatedOpponent = true
					break
				}
			}

			if !defeatedOpponent {
				hasDefeatedAll = false
				break
			}
		}

		if hasDefeatedAll {
			return &candidate
		}
	}

	return nil
}

func compareNFLCommonGames(teams []NFLTeamRecord, games []NFLGameResult, minCommonGames int) []NFLTeamRecord {
	if len(teams) != 2 {
		return teams
	}

	teamA, teamB := teams[0], teams[1]

	// Find common opponents (teams both have played)
	aOpponents := make(map[int]bool)
	bOpponents := make(map[int]bool)

	for _, game := range games {
		if game.HomeTeamID == teamA.TeamID {
			aOpponents[game.AwayTeamID] = true
		} else if game.AwayTeamID == teamA.TeamID {
			aOpponents[game.HomeTeamID] = true
		}

		if game.HomeTeamID == teamB.TeamID {
			bOpponents[game.AwayTeamID] = true
		} else if game.AwayTeamID == teamB.TeamID {
			bOpponents[game.HomeTeamID] = true
		}
	}

	commonOpponents := []int{}
	for oppID := range aOpponents {
		if bOpponents[oppID] {
			commonOpponents = append(commonOpponents, oppID)
		}
	}

	// Check minimum common games requirement
	if len(commonOpponents) < minCommonGames {
		return teams
	}

	// Calculate records against common opponents
	aWins, aLosses, aTies := 0, 0, 0
	bWins, bLosses, bTies := 0, 0, 0

	for _, game := range games {
		// isCommon := false
		for _, oppID := range commonOpponents {
			if (game.HomeTeamID == teamA.TeamID && game.AwayTeamID == oppID) || (game.AwayTeamID == teamA.TeamID && game.HomeTeamID == oppID) {
				// isCommon = true

				if game.HomeScore > game.AwayScore {
					if game.HomeTeamID == teamA.TeamID {
						aWins++
					} else {
						aLosses++
					}
				} else if game.AwayScore > game.HomeScore {
					if game.AwayTeamID == teamA.TeamID {
						aWins++
					} else {
						aLosses++
					}
				} else {
					aTies++
				}
				break
			}

			if (game.HomeTeamID == teamB.TeamID && game.AwayTeamID == oppID) || (game.AwayTeamID == teamB.TeamID && game.HomeTeamID == oppID) {
				// isCommon = true

				if game.HomeScore > game.AwayScore {
					if game.HomeTeamID == teamB.TeamID {
						bWins++
					} else {
						bLosses++
					}
				} else if game.AwayScore > game.HomeScore {
					if game.AwayTeamID == teamB.TeamID {
						bWins++
					} else {
						bLosses++
					}
				} else {
					bTies++
				}
				break
			}
		}
	}

	aPct := calculateNFLWinPct(aWins, aLosses, aTies)
	bPct := calculateNFLWinPct(bWins, bLosses, bTies)

	if aPct > bPct {
		return []NFLTeamRecord{teamA}
	} else if bPct > aPct {
		return []NFLTeamRecord{teamB}
	}

	return teams
}

func removeNFLTeam(teams []NFLTeamRecord, teamID int) []NFLTeamRecord {
	var result []NFLTeamRecord
	for _, team := range teams {
		if team.TeamID != teamID {
			result = append(result, team)
		}
	}
	return result
}

func calculateNFLDraftOrder(allTeams []NFLTeamRecord, afc NFLConferenceStandings, nfc NFLConferenceStandings) []NFLDraftPick {
	// Get playoff teams (first 7 from each conference)
	playoffTeamIDs := make(map[int]bool)
	for _, seed := range afc.PlayoffSeeds {
		if seed.Seed <= 7 {
			playoffTeamIDs[seed.Team.TeamID] = true
		}
	}
	for _, seed := range nfc.PlayoffSeeds {
		if seed.Seed <= 7 {
			playoffTeamIDs[seed.Team.TeamID] = true
		}
	}

	// Get non-playoff teams
	var nonPlayoffTeams []NFLTeamRecord
	for _, team := range allTeams {
		if !playoffTeamIDs[team.TeamID] {
			nonPlayoffTeams = append(nonPlayoffTeams, team)
		}
	}

	// Sort non-playoff teams by record (worst to best)
	sortedNonPlayoff := applyNFLDraftOrderTiebreakers(nonPlayoffTeams, []NFLGameResult{}, afc, nfc)

	// Build draft order
	draftOrder := []NFLDraftPick{}
	pickNum := 1

	// Picks 1-18: Non-playoff teams
	for i := len(sortedNonPlayoff) - 1; i >= 0; i-- {
		team := sortedNonPlayoff[i]
		draftOrder = append(draftOrder, NFLDraftPick{
			Pick: pickNum,
			Team: team,
		})
		pickNum++
	}

	// Picks 19-32: Playoff teams (ordered by seed, worst to best)
	// TODO: NEEDS TO BE UPDATED AS PLAYOFFS PROGRESS
	playoffTeams := make([]NFLPlayoffSeed, 0)

	// Combine both conferences' playoff teams and sort by seed (descending)
	for _, seed := range afc.PlayoffSeeds {
		if seed.Seed <= 7 {
			playoffTeams = append(playoffTeams, seed)
		}
	}
	for _, seed := range nfc.PlayoffSeeds {
		if seed.Seed <= 7 {
			playoffTeams = append(playoffTeams, seed)
		}
	}

	// Sort by seed (7 seed picks before 6 seed, etc.)
	sort.Slice(playoffTeams, func(i, j int) bool {
		return playoffTeams[i].Seed > playoffTeams[j].Seed
	})

	// Add playoff teams to draft order
	for _, seed := range playoffTeams {
		draftOrder = append(draftOrder, NFLDraftPick{
			Pick: pickNum,
			Team: seed.Team,
		})
		pickNum++
	}

	return draftOrder
}

func applyNFLDraftOrderTiebreakers(teams []NFLTeamRecord, games []NFLGameResult, afc NFLConferenceStandings, nfc NFLConferenceStandings) []NFLTeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// Group teams by win percentage
	pctGroups := make(map[float64][]NFLTeamRecord)
	for _, team := range teams {
		pctGroups[team.WinPct] = append(pctGroups[team.WinPct], team)
	}

	var result []NFLTeamRecord

	// Sort each group
	for _, group := range pctGroups {
		if len(group) == 1 {
			result = append(result, group[0])
		} else if len(group) == 2 {
			sorted := resolveNFLTwoTeamDraftTie(group, games, afc, nfc)
			result = append(result, sorted...)
		} else {
			sorted := resolveNFLMultiTeamDraftTie(group, games, afc, nfc)
			result = append(result, sorted...)
		}
	}

	// Sort groups by win percentage
	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct
	})

	return result
}

func resolveNFLTwoTeamDraftTie(teams []NFLTeamRecord, games []NFLGameResult, afc NFLConferenceStandings, nfc NFLConferenceStandings) []NFLTeamRecord {
	a, b := teams[0], teams[1]

	// Step 1: Strength of schedule (worse gets earlier pick)
	if a.StrengthOfSchedule != b.StrengthOfSchedule {
		if a.StrengthOfSchedule < b.StrengthOfSchedule {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 2: Division rank if same division (lower rank = earlier pick)
	if a.Division == b.Division {
		var divTeams []NFLTeamRecord
		conf := a.Conference
		div := a.Division

		if conf == "AFC" {
			divTeams = afc.Divisions[div]
		} else {
			divTeams = nfc.Divisions[div]
		}

		// Find ranks
		aRank, bRank := -1, -1
		for i, team := range divTeams {
			if team.TeamID == a.TeamID {
				aRank = i + 1
			}
			if team.TeamID == b.TeamID {
				bRank = i + 1
			}
		}

		if aRank != bRank && aRank != -1 && bRank != -1 {
			if aRank > bRank { // Higher rank number = worse record = earlier pick
				return []NFLTeamRecord{a, b}
			}
			return []NFLTeamRecord{b, a}
		}
	}

	// Step 3: Conference rank if same conference (but different division)
	if a.Conference == b.Conference {
		var confSeeds []NFLPlayoffSeed
		if a.Conference == "AFC" {
			confSeeds = afc.PlayoffSeeds
		} else {
			confSeeds = nfc.PlayoffSeeds
		}

		aRank, bRank := -1, -1
		for _, seed := range confSeeds {
			if seed.Team.TeamID == a.TeamID {
				aRank = seed.Seed
			}
			if seed.Team.TeamID == b.TeamID {
				bRank = seed.Seed
			}
		}

		if aRank != bRank && aRank != -1 && bRank != -1 {
			if aRank > bRank {
				return []NFLTeamRecord{a, b}
			}
			return []NFLTeamRecord{b, a}
		}
	}

	// Inter-conference tiebreakers (reversed for draft order)

	// Step 1: Head-to-head (if applicable) (loser gets earlier pick)
	h2hResult := compareNFLHeadToHead(teams, games)
	if len(h2hResult) == 1 {
		if h2hResult[0].TeamID == a.TeamID {
			return []NFLTeamRecord{b, a}
		}
		return []NFLTeamRecord{a, b}
	}

	// Step 2: Worst win percentage in common games (minimum of 4)
	commonResults := compareNFLCommonGames(teams, games, 4)
	if len(commonResults) == 1 {
		if commonResults[0].TeamID == a.TeamID {
			return []NFLTeamRecord{b, a}
		}
		return []NFLTeamRecord{a, b}
	}

	// Step 3: Strength of victory (lower gets earlier pick)
	if a.StrengthOfVictory != b.StrengthOfVictory {
		if a.StrengthOfVictory < b.StrengthOfVictory {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Step 4: Worst point differential
	aDiff := a.PointsFor - a.PointsAgainst
	bDiff := b.PointsFor - b.PointsAgainst
	if aDiff != bDiff {
		if aDiff < bDiff {
			return []NFLTeamRecord{a, b}
		}
		return []NFLTeamRecord{b, a}
	}

	// Random drawing - use TeamID for consistency
	if a.TeamID < b.TeamID {
		return []NFLTeamRecord{a, b}
	}
	return []NFLTeamRecord{b, a}
}

func resolveNFLMultiTeamDraftTie(teams []NFLTeamRecord, games []NFLGameResult, afc NFLConferenceStandings, nfc NFLConferenceStandings) []NFLTeamRecord {
	if len(teams) == 1 {
		return teams
	}
	if len(teams) == 2 {
		return resolveNFLTwoTeamDraftTie(teams, games, afc, nfc)
	}

	// Check if all same division
	allSameDivision := true
	firstDiv := teams[0].Division
	for _, team := range teams {
		if team.Division != firstDiv {
			allSameDivision = false
			break
		}
	}
	if allSameDivision {
		conf := teams[0].Conference
		div := teams[0].Division

		var divTeams []NFLTeamRecord
		if conf == "AFC" {
			divTeams = afc.Divisions[div]
		} else {
			divTeams = nfc.Divisions[div]
		}

		// Sort by division rank
		var result []NFLTeamRecord
		for _, divTeam := range divTeams {
			for _, team := range teams {
				if team.TeamID == divTeam.TeamID {
					result = append(result, team)
					break
				}
			}
		}

		// Reverse for draft order
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}

		return result
	}

	// Check if all same conference
	allSameConference := true
	firstConference := teams[0].Conference
	for _, team := range teams[1:] {
		if team.Conference != firstConference {
			allSameConference = false
			break
		}
	}

	// If all same conference, use conference rank
	if allSameConference {
		var confSeeds []NFLPlayoffSeed
		if firstConference == "AFC" {
			confSeeds = afc.PlayoffSeeds
		} else {
			confSeeds = nfc.PlayoffSeeds
		}

		// Sort by conference rank
		sort.Slice(teams, func(i, j int) bool {
			iSeed, jSeed := -1, -1
			for _, seed := range confSeeds {
				if seed.Team.TeamID == teams[i].TeamID {
					iSeed = seed.Seed
				}
				if seed.Team.TeamID == teams[j].TeamID {
					jSeed = seed.Seed
				}
			}
			return iSeed > jSeed // Higher seed number = worse record = earlier pick
		})
		return teams
	}

	// Inter-conference
	// Group by division
	divisionGroups := make(map[string][]NFLTeamRecord)
	for _, team := range teams {
		key := team.Conference + " " + team.Division
		divisionGroups[key] = append(divisionGroups[key], team)
	}

	// Determine lowest-ranked team in each division
	var divRepresentatives []NFLTeamRecord
	for _, divTeams := range divisionGroups {
		if len(divTeams) == 1 {
			divRepresentatives = append(divRepresentatives, divTeams[0])
		} else {
			sorted := applyNFLDivisionTiebreakers(divTeams, games)
			divRepresentatives = append(divRepresentatives, sorted[len(sorted)-1])
		}
	}

	// If only 2 teams remain, use two-team inter-conference tiebreakers
	if len(divRepresentatives) == 2 {
		return resolveNFLTwoTeamDraftTie(divRepresentatives, games, afc, nfc)
	}

	// Group remaining teams by conference
	confGroups := make(map[string][]NFLTeamRecord)
	for _, team := range divRepresentatives {
		confGroups[team.Conference] = append(confGroups[team.Conference], team)
	}

	// Determine lowest-ranked team in each conference
	var finalRepresentatives []NFLTeamRecord
	for conf, confTeams := range confGroups {
		if len(confTeams) == 1 {
			finalRepresentatives = append(finalRepresentatives, confTeams[0])
		} else {
			var confSeeds []NFLPlayoffSeed
			if conf == "AFC" {
				confSeeds = afc.PlayoffSeeds
			} else {
				confSeeds = nfc.PlayoffSeeds
			}

			// Find lowest-ranked team
			lowestRanked := confTeams[0]
			lowestSeed := -1
			for _, seed := range confSeeds {
				if seed.Team.TeamID == lowestRanked.TeamID {
					lowestSeed = seed.Seed
					break
				}
			}

			for _, team := range confTeams[1:] {
				for _, seed := range confSeeds {
					if seed.Team.TeamID == team.TeamID {
						if seed.Seed > lowestSeed {
							lowestRanked = team
						}
						break
					}
				}
			}
			finalRepresentatives = append(finalRepresentatives, lowestRanked)
		}
	}

	// Now should have max 2 inter-conference teams
	if len(finalRepresentatives) <= 2 {
		return resolveNFLTwoTeamDraftTie(finalRepresentatives, games, afc, nfc)
	}

	// Fallback: sort by strength of schedule
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].StrengthOfSchedule < teams[j].StrengthOfSchedule
	})
	return teams
}

func calculateNFLWinPct(wins, losses, ties int) float64 {
	total := wins + losses + ties
	if total == 0 {
		return -1.0 // Treat 0-0 as worse than any actual record
	}
	return (float64(wins) + 0.5*float64(ties)) / float64(total)
}

func calculateNFLGamesBack(leader NFLTeamRecord, team NFLTeamRecord) float64 {
	if leader.TeamID == team.TeamID {
		return 0.0
	}

	// Formula: [(Leader Wins - Team Wins) + (Team Losses - Leader Losses)] / 2
	leaderWins := float64(leader.Wins) + (0.5 * float64(leader.Ties))
	teamWins := float64(team.Wins) + (0.5 * float64(team.Ties))
	leaderLosses := float64(leader.Losses) + (0.5 * float64(leader.Ties))
	teamLosses := float64(team.Losses) + (0.5 * float64(team.Ties))

	return ((leaderWins - teamWins) + (teamLosses - leaderLosses)) / 2.0
}
