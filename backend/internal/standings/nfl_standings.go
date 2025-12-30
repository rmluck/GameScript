package standings

import (
	"fmt"
	"log"
	"sort"

	"gamescript/internal/database"
)

type TeamRecord struct {
	TeamID       		int
	TeamCity   			string
	TeamName     		string
	TeamAbbr	 		string
	Conference   		string
	Division     		string
	Wins         		int
	Losses       		int
	Ties         		int
	HomeWins	 		int
	HomeLosses 			int
	HomeTies    		int
	AwayWins	 		int
	AwayLosses 			int
	AwayTies    		int
	DivisionWins 		int
	DivisionLosses 		int
	DivisionTies   		int
	ConferenceWins 		int
	ConferenceLosses 	int
	ConferenceTies   	int
	PointsFor    		int
	PointsAgainst 		int
	WinPct       		float64
	ConferenceGamesBack	float64
	DivisionGamesBack	float64
	StrengthOfSchedule	float64
	StrengthOfVictory	float64
	LogoURL	 			string
	TeamPrimaryColor	string
	TeamSecondaryColor	string
}

type NFLStandings struct {
	AFC 		ConferenceStandings
	NFC 		ConferenceStandings
	DraftOrder 	[]DraftPick
}

type ConferenceStandings struct {
	Divisions 		map[string][]TeamRecord // Keyed by division name
	PlayoffSeeds 	[]PlayoffSeed
}

type PlayoffSeed struct {
	Seed				int
	Team				TeamRecord
	IsDivisionWinner	bool
}

type DraftPick struct {
	Pick 		int
	Team		TeamRecord
	Reason 		string // e.g., "Non-playoff", "Wild card loss", etc.
}

type GameResult struct {
	GameID			int
	HomeTeamID		int
	AwayTeamID		int
	HomeScore		int
	AwayScore		int
	Week 			int
}

func CalculateNFLStandings(db *database.DB, scenarioID int, seasonID int) (*NFLStandings, error) {
	// Get all teams for the season
	teams, err := getTeams(db, seasonID)
	if err != nil {
		return nil, fmt.Errorf("error getting teams: %w", err)
	}

	// Get all game results for the scenario (actual + picks)
	games, err := getGameResults(db, scenarioID, seasonID)
	if err != nil {
		return nil, fmt.Errorf("error getting game rsults: %w", err)
	}

	// Calculate team records based on game results
	records := calculateTeamRecords(teams, games)

	// Calculate strength metrics for all teams
	calculateStrengthMetrics(records, games)

	// Separate by conference
	afcTeams := filterByConference(records, "AFC")
	nfcTeams := filterByConference(records, "NFC")

	// Calculate playoff seeds for each conference
	afcStandings := calculateConferenceStandings(afcTeams, games)
	nfcStandings := calculateConferenceStandings(nfcTeams, games)

	// Calculate draft order
	draftOrder := calculateDraftOrder(records, afcStandings, nfcStandings)

	return &NFLStandings{
		AFC: afcStandings,
		NFC: nfcStandings,
		DraftOrder: draftOrder,
	}, nil
}

func getTeams(db *database.DB, seasonID int) ([]TeamRecord, error) {
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

	var teams []TeamRecord
	for rows.Next() {
		var team TeamRecord
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

func getGameResults(db *database.DB, scenarioID int, seasonID int) ([]GameResult, error) {
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

    var games []GameResult
    for rows.Next() {
        var game GameResult
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

        // Priority 1: User has made a pick - always use the pick
        if pickedTeamID != nil {
            // If user provided predicted scores, use those
            if predictedHomeScore != nil && predictedAwayScore != nil {
                game.HomeScore = *predictedHomeScore
                game.AwayScore = *predictedAwayScore
            } else {
                // No predicted scores, use dummy scores based on picked winner
                if *pickedTeamID == game.HomeTeamID {
                    // Home team picked to win
                    game.HomeScore = 1
                    game.AwayScore = 0
                } else if *pickedTeamID == game.AwayTeamID {
                    // Away team picked to win
                    game.HomeScore = 0
                    game.AwayScore = 1
                } else if *pickedTeamID == 0 {
                    // Tie picked
                    game.HomeScore = 0
                    game.AwayScore = 0
                } else {
                    // Invalid picked_team_id, skip this game
                    continue
                }
            }
            games = append(games, game)
            continue
        }

        // Priority 2: No pick, but game is final - use actual scores
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

func calculateTeamRecords(teams []TeamRecord, games []GameResult) []TeamRecord {
	// Create map for quick lookup
	recordMap := make(map[int]*TeamRecord)
	for i := range teams {
		recordMap[teams[i].TeamID] = &teams[i]
	}

	// Process each game
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

			// Check if division game
			if isDivisionGame {
				homeTeam.DivisionWins++
				awayTeam.DivisionLosses++
			}

			// Check if conference game
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

		// Update points for and against
		homeTeam.PointsFor += game.HomeScore
		homeTeam.PointsAgainst += game.AwayScore
		awayTeam.PointsFor += game.AwayScore
		awayTeam.PointsAgainst += game.HomeScore
	}

	// Calculate win percentages
	for _, team := range recordMap {
		team.WinPct = calculateWinPct(team.Wins, team.Losses, team.Ties)
	}

	return teams
}

func calculateStrengthMetrics(teams []TeamRecord, games []GameResult) {
	teamMap := make(map[int]*TeamRecord)
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
			team.StrengthOfSchedule = calculateWinPct(opponentTotalWins, opponentTotalLosses, opponentTotalTies)
		}
		if defeatedCount > 0 {
			team.StrengthOfVictory = calculateWinPct(defeatedOpponentWins, defeatedOpponentLosses, defeatedOpponentTies)
		}
	}
}

func filterByConference(teams []TeamRecord, conference string) []TeamRecord {
	var filtered []TeamRecord
	for _, team := range teams {
		if team.Conference == conference {
			filtered = append(filtered, team)
		}
	}
	return filtered
}

func calculateConferenceStandings(teams []TeamRecord, games []GameResult) ConferenceStandings {
	// Group teams by division
	divisions := make(map[string][]TeamRecord)
	for _, team := range teams {
		divisions[team.Division] = append(divisions[team.Division], team)
	}

	// Determine division winners
	divisionWinners := []TeamRecord{}
	nonWinners := []TeamRecord{}

	for divName, divTeams := range divisions {
		// Sort division teams
		sortedDiv := applyDivisionTiebreakers(divTeams, games)

		// Calculate division games back
		divLeader := sortedDiv[0]
		for i := range sortedDiv {
			sortedDiv[i].DivisionGamesBack = calculateGamesBack(divLeader, sortedDiv[i])
		}

		// First place is division winner
		divisionWinners = append(divisionWinners, sortedDiv[0])

		// Rest are wild card candidates
		nonWinners = append(nonWinners, sortedDiv[1:]...)

		// Update divisions map with sorted teams
		divisions[divName] = sortedDiv
	}

	// Rank division winners (seeds 1-4)
	divisionWinners = applyConferenceTiebreakers(divisionWinners, games, true)

	// Rank non-division winners (seeds 5-16)
	nonWinners = applyConferenceTiebreakers(nonWinners, games, false)

	// Create playoff seeds
	playoffSeeds := []PlayoffSeed{}
	for i, team := range divisionWinners {
		playoffSeeds = append(playoffSeeds, PlayoffSeed{
			Seed: i + 1,
			Team: team,
			IsDivisionWinner: true,
		})
	}

	for i, team := range nonWinners {
		playoffSeeds = append(playoffSeeds, PlayoffSeed{
			Seed: i + 5,
			Team: team,
			IsDivisionWinner: false,
		})
	}

	// Calculate conference games back
	conferenceLeader := playoffSeeds[0].Team
	for i := range playoffSeeds {
		playoffSeeds[i].Team.ConferenceGamesBack = calculateGamesBack(conferenceLeader, playoffSeeds[i].Team)
	}

	return ConferenceStandings{
		Divisions: divisions,
		PlayoffSeeds: playoffSeeds,
	}
}

func applyDivisionTiebreakers(teams []TeamRecord, games []GameResult) []TeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// Group by win percentage
	pctGroups := make(map[float64][]TeamRecord)
	for _, team := range teams {
		pctGroups[team.WinPct] = append(pctGroups[team.WinPct], team)
	}

	var result []TeamRecord
	for _, group := range pctGroups {
		if len(group) == 1 {
			result = append(result, group[0])
		} else if len(group) == 2 {
			result = append(result, resolveTwoTeamDivisionTie(group, games)...)
		} else {
			result = append(result, resolveMultiTeamDivisionTie(group, games)...)
		}
	}

	// Sort by win percentage
	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct
	})

	return result
}

func resolveTwoTeamDivisionTie(teams []TeamRecord, games []GameResult) []TeamRecord {
	a, b := teams[0], teams[1]

	// Step 1: Win percentage - already tied

	// Step 2: Head-to-head
	h2h := compareHeadToHead(teams, games)
	if len(h2h) == 1 {
		if h2h[0].TeamID == a.TeamID {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 3: Division record
	aDivPct := calculateWinPct(a.DivisionWins, a.DivisionLosses, a.DivisionTies)
	bDivPct := calculateWinPct(b.DivisionWins, b.DivisionLosses, b.DivisionTies)
	if aDivPct != bDivPct {
		if aDivPct > bDivPct {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 4: Common games
	commonResults := compareCommonGames(teams, games, 0)
	if len(commonResults) == 1 {
		if commonResults[0].TeamID == a.TeamID {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 5: Conference record
	aConfPct := calculateWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
	bConfPct := calculateWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
	if aConfPct != bConfPct {
		if aConfPct > bConfPct {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 6: Strength of victory
	if a.StrengthOfVictory != b.StrengthOfVictory {
		if a.StrengthOfVictory > b.StrengthOfVictory {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 7: Strength of schedule
	if a.StrengthOfSchedule != b.StrengthOfSchedule {
		if a.StrengthOfSchedule > b.StrengthOfSchedule {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 8: Point differential
	aDiff := a.PointsFor - a.PointsAgainst
	bDiff := b.PointsFor - b.PointsAgainst
	if aDiff != bDiff {
		if aDiff > bDiff {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 9: Points scored
	if a.PointsFor != b.PointsFor {
		if a.PointsFor > b.PointsFor {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 10: Points allowed (fewer is better)
	if a.PointsAgainst != b.PointsAgainst {
		if a.PointsAgainst < b.PointsAgainst {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 11: Coin toss - not implemented, use TeamID for consistency
	if a.TeamID < b.TeamID {
		return []TeamRecord{a, b}
	}
	return []TeamRecord{b, a}
}

func resolveMultiTeamDivisionTie(teams []TeamRecord, games []GameResult) []TeamRecord {
	// Step 1: Head-to-head (best win pct in games among tied teams)
	h2hWinner := findHeadToHeadWinner(teams, games)
	if h2hWinner != nil {
		remaining := removeTeam(teams, h2hWinner.TeamID)
		result := []TeamRecord{*h2hWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 2: Division record
	divWinner := findBestDivisionRecord(teams)
	if divWinner != nil {
		remaining := removeTeam(teams, divWinner.TeamID)
		result := []TeamRecord{*divWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 3: Common games
	commonWinner := findBestCommonGamesRecord(teams, games, 0)
	if commonWinner != nil {
		remaining := removeTeam(teams, commonWinner.TeamID)
		result := []TeamRecord{*commonWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 4: Conference record
	confWinner := findBestConferenceRecord(teams)
	if confWinner != nil {
		remaining := removeTeam(teams, confWinner.TeamID)
		result := []TeamRecord{*confWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 5: Strength of victory
	sovWinner := findBestStrengthOfVictory(teams)
	if sovWinner != nil {
		remaining := removeTeam(teams, sovWinner.TeamID)
		result := []TeamRecord{*sovWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 6: Strength of schedule
	sosWinner := findBestStrengthOfSchedule(teams)
	if sosWinner != nil {
		remaining := removeTeam(teams, sosWinner.TeamID)
		result := []TeamRecord{*sosWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 7: Point differential
	pdWinner := findBestPointDifferential(teams)
	if pdWinner != nil {
		remaining := removeTeam(teams, pdWinner.TeamID)
		result := []TeamRecord{*pdWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 8: Points scored
	psWinner := findBestPointsScored(teams)
	if psWinner != nil {
		remaining := removeTeam(teams, psWinner.TeamID)
		result := []TeamRecord{*psWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 9: Points allowed (fewer is better)
	paWinner := findBestPointsAllowed(teams)
	if paWinner != nil {
		remaining := removeTeam(teams, paWinner.TeamID)
		result := []TeamRecord{*paWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamDivisionTie(remaining, games)...)
		}
		return result
	}

	// Step 10: Coin toss - not implemented, use TeamID for consistency
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].TeamID < teams[j].TeamID
	})
	return teams
}

func applyConferenceTiebreakers(teams []TeamRecord, games []GameResult, areDivisionWinners bool) []TeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// DEBUG
	log.Printf("=== applyConferenceTiebreakers called with %d teams ===", len(teams))
	for _, t := range teams {
		log.Printf("  - %s %s: %d-%d (%.3f)", t.TeamAbbr, t.Division, t.Wins, t.Losses, t.WinPct)
	}

	// Group by win percentage
	pctGroups := make(map[float64][]TeamRecord)
	for _, team := range teams {
		pctGroups[team.WinPct] = append(pctGroups[team.WinPct], team)
	}

	var result []TeamRecord
	for pct, group := range pctGroups {
		// DEBUG
		log.Printf("Processing win pct group %.3f with %d teams:", pct, len(group))
		for _, t := range group {
			log.Printf("  - %s", t.TeamAbbr)
		}

		if len(group) == 1 {
			result = append(result, group[0])
			// DEBUG
			log.Printf("  -> Single team, added %s", group[0].TeamAbbr)
		} else if len(group) == 2 {
			resolved := resolveTwoTeamConferenceTie(group, games)
			// DEBUG
			log.Printf("  --> Two teams resolved to: %s, %s", resolved[0].TeamAbbr, resolved[1].TeamAbbr)
			result = append(result, resolved...)
		} else {
			resolved := resolveMultiTeamConferenceTie(group, games)
			// DEBUG
			log.Printf("  -> Multi-team resolved to %d teams:", len(resolved))
			for _, t := range resolved {
				log.Printf("   - %s", t.TeamAbbr)
			}
			result = append(result, resolved...)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct
	})

	// DEBUG
	log.Printf("=== Final result has %d teams ===", len(result))
	for i, t := range result {
		log.Printf("  %d. %s %s: %d-%d", i+1, t.TeamAbbr, t.Division, t.Wins, t.Losses)
	}

	return result
}

func resolveTwoTeamConferenceTie(teams []TeamRecord, games []GameResult) []TeamRecord {
	a, b := teams[0], teams[1]

	// If same division, use division tiebreakers
	if a.Division == b.Division {
		return resolveTwoTeamDivisionTie(teams, games)
	}

	// Step 2: Head-to-head
	h2h := compareHeadToHead(teams, games)
	if len(h2h) == 1 {
		if h2h[0].TeamID == a.TeamID {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 3: Conference record
	aConfPct := calculateWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
	bConfPct := calculateWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
	if aConfPct != bConfPct {
		if aConfPct > bConfPct {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 4: Common games (minimum of 4)
	commonResults := compareCommonGames(teams, games, 4)
	if len(commonResults) == 1 {
		if commonResults[0].TeamID == a.TeamID {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 5: Strength of victory
	if a.StrengthOfVictory != b.StrengthOfVictory {
		if a.StrengthOfVictory > b.StrengthOfVictory {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 6: Strength of schedule
	if a.StrengthOfSchedule != b.StrengthOfSchedule {
		if a.StrengthOfSchedule > b.StrengthOfSchedule {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 7: Point differential
	aDiff := a.PointsFor - a.PointsAgainst
	bDiff := b.PointsFor - b.PointsAgainst
	if aDiff != bDiff {
		if aDiff > bDiff {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 8: Points scored
	if a.PointsFor != b.PointsFor {
		if a.PointsFor > b.PointsFor {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 9: Points allowed (fewer is better)
	if a.PointsAgainst != b.PointsAgainst {
		if a.PointsAgainst < b.PointsAgainst {
			return []TeamRecord{a, b}
		}
		return []TeamRecord{b, a}
	}

	// Step 10: Coin toss - not implemented, use TeamID for consistency
	if a.TeamID < b.TeamID {
		return []TeamRecord{a, b}
	}
	return []TeamRecord{b, a}
}

func resolveMultiTeamConferenceTie(teams []TeamRecord, games []GameResult) []TeamRecord {
	// DEBUG
	log.Printf("   >>> resolveMultiTeamConferenceTie called with %d teams", len(teams))
	for _, t := range teams {
		log.Printf("        - %s (%s)", t.TeamAbbr, t.Division)
	}

	if len(teams) == 1 {
		// DEBUG
		log.Printf("   >>> Only one team, returning %s", teams[0].TeamAbbr)
		return teams
	}
	if len(teams) == 2 {
		// DEBUG
		log.Printf("   >>> Two teams, deferring to resolveTwoTeamConferenceTie")
		return resolveTwoTeamConferenceTie(teams, games)
	}

	// Check if all teams are from the same division
	allSameDivision := true
	firstDivision := teams[0].Division
	for _, team := range teams[1:] {
		if team.Division != firstDivision {
			allSameDivision = false
			break
		}
	}

	// If all from same division, use division tiebreakers
	if allSameDivision {
		// DEBUG
		log.Printf("   >>> All teams from same division (%s), deferring to resolveMultiTeamDivisionTie", firstDivision)
		return resolveMultiTeamDivisionTie(teams, games)
	}

	// Step 1: Apply division tiebreaker to get best from each division
	// DEBUG
	log.Printf("   >>> Applying division tiebreakers to get best from each division")
	divGroups := make(map[string][]TeamRecord)

	for _, team := range teams {
		divGroups[team.Division] = append(divGroups[team.Division], team)
	}

	// DEBUG
	log.Printf("   >>> Found %d divisions in tie", len(divGroups))
	for div, divTeams := range divGroups {
		log.Printf("        - Division %s has %d teams", div, len(divTeams))
		for _, t := range divTeams {
			log.Printf("            %s", t.TeamAbbr)
		}
	}

	var filtered []TeamRecord
	for div, divTeams := range divGroups {
		if len(divTeams) == 1 {
			filtered = append(filtered, divTeams[0])
			// DEBUG
			log.Printf("        - Division %s has single team %s, added to filtered", div, divTeams[0].TeamAbbr)
		} else {
			sorted := applyDivisionTiebreakers(divTeams, games)
			filtered = append(filtered, sorted[0])
			// DEBUG
			log.Printf("        - Division %s has multiple teams, added best team %s to filtered", div, sorted[0].TeamAbbr)
		}
	}

	// DEBUG
	log.Printf("   >>> After division tiebreakers, %d teams remain:", len(filtered))
	for _, t := range filtered {
		log.Printf("        - %s (%s)", t.TeamAbbr, t.Division)
	}

	if len(filtered) == 1 {
		// DEBUG
		log.Printf("   >>> Only one team after division tiebreakers, returning %s", filtered[0].TeamAbbr)
		winner := filtered[0]
		remaining := removeTeam(teams, winner.TeamID)
		result := []TeamRecord{winner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// if len(filtered) == 2 {
	// 	// DEBUG
	// 	log.Printf("   >>> Two teams after division tiebreakers, deferring to resolveTwoTeamConferenceTie")
	// 	return resolveTwoTeamConferenceTie(filtered, games)
	// }

	// Step 2: Head-to-head sweep
	// DEBUG
	log.Printf("   >>> Checking for head-to-head sweep among %d teams", len(filtered))
	sweepWinner := checkHeadToHeadSweep(filtered, games)
	if sweepWinner != nil {
		// DEBUG
		log.Printf("   >>> Found head-to-head sweep winner: %s", sweepWinner.TeamAbbr)
		remaining := removeTeam(teams, sweepWinner.TeamID)
		// DEBUG
		log.Printf("   >>> %d teams remain after removing sweep winner", len(remaining))
		result := []TeamRecord{*sweepWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}
	// DEBUG
	log.Printf("   >>> No head-to-head sweep winner found")

	// Step 3: Conference record
	// DEBUG
	log.Printf("   >>> Checking conference records among %d teams", len(filtered))
	confWinner := findBestConferenceRecord(filtered)
	if confWinner != nil {
		// DEBUG
		log.Printf("   >>> Found conference record winner: %s", confWinner.TeamAbbr)
		remaining := removeTeam(teams, confWinner.TeamID)
		// DEBUG
		log.Printf("   >>> %d teams remain after removing conference record winner", len(remaining))
		result := []TeamRecord{*confWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}
	// DEBUG
	log.Printf("   >>> No conference record winner found")

	// Step 4: Common games (minimum of 4)
	// DEBUG
	log.Printf("   >>> Checking common games among %d teams", len(filtered))
	commonWinner := findBestCommonGamesRecord(filtered, games, 4)
	if commonWinner != nil {
		// DEBUG
		log.Printf("   >>> Found common games winner: %s", commonWinner.TeamAbbr)
		remaining := removeTeam(teams, commonWinner.TeamID)
		// DEBUG
		log.Printf("   >>> %d teams remain after removing common games winner", len(remaining))
		result := []TeamRecord{*commonWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}
	// DEBUG
	log.Printf("   >>> No common games winner found")

	// Step 5: Strength of victory
	// DEBUG
	log.Printf("   >>> Checking strength of victory among %d teams", len(filtered))
	if sovWinner := findBestStrengthOfVictory(filtered); sovWinner != nil {
		// DEBUG
		log.Printf("   >>> Found strength of victory winner: %s", sovWinner.TeamAbbr)
		remaining := removeTeam(teams, sovWinner.TeamID)
		// DEBUG
		log.Printf("   >>> %d teams remain after removing strength of victory winner", len(remaining))
		result := []TeamRecord{*sovWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}
	// DEBUG
	log.Printf("   >>> No strength of victory winner found")

	// Step 6: Strength of schedule
	// DEBUG
	log.Printf("   >>> Checking strength of schedule among %d teams", len(filtered))
	if sosWinner := findBestStrengthOfSchedule(filtered); sosWinner != nil {
		// DEBUG
		log.Printf("   >>> Found strength of schedule winner: %s", sosWinner.TeamAbbr)
		remaining := removeTeam(teams, sosWinner.TeamID)
		// DEBUG
		log.Printf("   >>> %d teams remain after removing strength of schedule winner", len(remaining))
		result := []TeamRecord{*sosWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}
	// DEBUG
	log.Printf("   >>> No strength of schedule winner found")

	// Point differential
	// DEBUG
	log.Printf("   >>> Checking point differential among %d teams", len(filtered))
	if pdWinner := findBestPointDifferential(filtered); pdWinner != nil {
		// DEBUG
		log.Printf("   >>> Found point differential winner: %s", pdWinner.TeamAbbr)
		remaining := removeTeam(teams, pdWinner.TeamID)
		// DEBUG
		log.Printf("   >>> %d teams remain after removing point differential winner", len(remaining))
		result := []TeamRecord{*pdWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}
	// DEBUG
	log.Printf("   >>> No point differential winner found")

	// Points scored
	if psWinner := findBestPointsScored(filtered); psWinner != nil {
		remaining := removeTeam(teams, psWinner.TeamID)
		result := []TeamRecord{*psWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Points allowed
	if paWinner := findBestPointsAllowed(filtered); paWinner != nil {
		remaining := removeTeam(teams, paWinner.TeamID)
		result := []TeamRecord{*paWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
		}
		return result
	}

	// Coin toss - not implemented, use TeamID for consistency
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TeamID < filtered[j].TeamID
	})
	
	winner := filtered[0]
	remaining := removeTeam(teams, winner.TeamID)
	result := []TeamRecord{winner}
	if len(remaining) > 0 {
		result = append(result, resolveMultiTeamConferenceTie(remaining, games)...)
	}
	return result
}

func findHeadToHeadWinner(teams []TeamRecord, games []GameResult) *TeamRecord {
	type h2hRecord struct {
		teamID int
		wins int
		losses int
		ties int
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
		pct := calculateWinPct(rec.wins, rec.losses, rec.ties)
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

func findBestDivisionRecord(teams []TeamRecord) *TeamRecord {
	var bestPct float64 = -2.0
	var bestTeam *TeamRecord
	var tie bool

	for i := range teams {
		pct := calculateWinPct(teams[i].DivisionWins, teams[i].DivisionLosses, teams[i].DivisionTies)
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

func findBestConferenceRecord(teams []TeamRecord) *TeamRecord {
	var bestPct float64 = -2.0
	var bestTeam *TeamRecord
	var tie bool

	for i := range teams {
		pct := calculateWinPct(teams[i].ConferenceWins, teams[i].ConferenceLosses, teams[i].ConferenceTies)
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

func findBestCommonGamesRecord(teams []TeamRecord, games []GameResult, minGames int) *TeamRecord {
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
		wins int
		losses int
		ties int
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
		pct := calculateWinPct(rec.wins, rec.losses, rec.ties)
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

func findBestStrengthOfVictory(teams []TeamRecord) *TeamRecord {
	var best float64 = -2.0
	var bestTeam *TeamRecord
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

func findBestStrengthOfSchedule(teams []TeamRecord) *TeamRecord {
	var best float64 = -2.0
	var bestTeam *TeamRecord
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

func findBestPointDifferential(teams []TeamRecord) *TeamRecord {
	var bestDiff int = -1 << 31
	var bestTeam *TeamRecord
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

func findBestPointsScored(teams []TeamRecord) *TeamRecord {
	var best int = -1
	var bestTeam *TeamRecord
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

func findBestPointsAllowed(teams []TeamRecord) *TeamRecord {
	var best int = 1 << 31
	var bestTeam *TeamRecord
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



// func applyDivisionTiebreakers(teams []TeamRecord, games []GameResult) []TeamRecord {
// 	return sortTeamsWithTiebreakers(teams, games, func(a, b TeamRecord) int {
// 		// Step 1: Win percentage
// 		if a.WinPct != b.WinPct {
// 			if a.WinPct > b.WinPct {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 2: Head-to-head record (if applicable)
// 		h2hResult := compareHeadToHead([]TeamRecord{a, b}, games)
// 		if len(h2hResult) == 1 {
// 			if h2hResult[0].TeamID == a.TeamID {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 3: Division record
// 		aDivPct := calculateWinPct(a.DivisionWins, a.DivisionLosses, a.DivisionTies)
// 		bDivPct := calculateWinPct(b.DivisionWins, b.DivisionLosses, b.DivisionTies)
// 		if aDivPct != bDivPct {
// 			if aDivPct > bDivPct {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 4: Common games
// 		commonResult := compareCommonGames([]TeamRecord{a, b}, games, 0) // No minimum for division
// 		if len(commonResult) == 1 {
// 			if commonResult[0].TeamID == a.TeamID {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 5: Conference record
// 		aConfPct := calculateWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
// 		bConfPct := calculateWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
// 		if aConfPct != bConfPct {
// 			if aConfPct > bConfPct {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 6: Strength of victory
// 		if a.StrengthOfVictory != b.StrengthOfVictory {
// 			if a.StrengthOfVictory > b.StrengthOfVictory {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 7: Strength of schedule
// 		if a.StrengthOfSchedule != b.StrengthOfSchedule {
// 			if a.StrengthOfSchedule > b.StrengthOfSchedule {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 8: Point differential
// 		aDiff := a.PointsFor - a.PointsAgainst
// 		bDiff := b.PointsFor - b.PointsAgainst
// 		if aDiff != bDiff {
// 			if aDiff > bDiff {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 9: Points scored
// 		if a.PointsFor != b.PointsFor {
// 			if a.PointsFor > b.PointsFor {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 10: Points allowed (fewer is better)
// 		if a.PointsAgainst != b.PointsAgainst {
// 			if a.PointsAgainst < b.PointsAgainst {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 11: Coin toss - not implemented, use TeamID for consistency
// 		if a.TeamID < b.TeamID {
// 			return -1
// 		}
// 		return 1
// 	})
// }

// func applyConferenceTiebreakers(teams []TeamRecord, games []GameResult, areDivisionWinners bool) []TeamRecord {
// 	return sortTeamsWithTiebreakers(teams, games, func(a, b TeamRecord) int {
// 		// Step 1: Win percentage
// 		if a.WinPct != b.WinPct {
// 			if a.WinPct > b.WinPct {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 2: If same, division, skip to division tiebreakers
// 		if a.Division == b.Division {
// 			divResult := applyDivisionTiebreakers([]TeamRecord{a, b}, games)
// 			if divResult[0].TeamID == a.TeamID {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 3: Head-to-head record (if applicable)
// 		h2hResult := compareHeadToHead([]TeamRecord{a, b}, games)
// 		if len(h2hResult) == 1 {
// 			if h2hResult[0].TeamID == a.TeamID {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 4: Conference record
// 		aConfPct := calculateWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
// 		bConfPct := calculateWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
// 		if aConfPct != bConfPct {
// 			if aConfPct > bConfPct {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 5: Common games (minimum 4 required)
// 		commonResult := compareCommonGames([]TeamRecord{a, b}, games, 4)
// 		if len(commonResult) == 1 {
// 			if commonResult[0].TeamID == a.TeamID {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 6: Strength of victory
// 		if a.StrengthOfVictory != b.StrengthOfVictory {
// 			if a.StrengthOfVictory > b.StrengthOfVictory {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 7: Strength of schedule
// 		if a.StrengthOfSchedule != b.StrengthOfSchedule {
// 			if a.StrengthOfSchedule > b.StrengthOfSchedule {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 8: Point differential
// 		aDiff := a.PointsFor - a.PointsAgainst
// 		bDiff := b.PointsFor - b.PointsAgainst
// 		if aDiff != bDiff {
// 			if aDiff > bDiff {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 9: Points scored
// 		if a.PointsFor != b.PointsFor {
// 			if a.PointsFor > b.PointsFor {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 10: Points allowed (fewer is better)
// 		if a.PointsAgainst != b.PointsAgainst {
// 			if a.PointsAgainst < b.PointsAgainst {
// 				return -1
// 			}
// 			return 1
// 		}

// 		// Step 11: Coin toss - not implemented, use TeamID for consistency
// 		if a.TeamID < b.TeamID {
// 			return -1
// 		}
// 		return 1
// 	})
// }

// Helper function to sort teams with multi-team tiebreaker support
func sortTeamsWithTiebreakers(teams []TeamRecord, games []GameResult, twoTeamCompare func(TeamRecord, TeamRecord) int) []TeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// Group teams by win percentage
	pctGroups := make(map[float64][]TeamRecord)
	for _, team := range teams {
		pctGroups[team.WinPct] = append(pctGroups[team.WinPct], team)
	}

	var result []TeamRecord

	// Sort each group
	for _, group := range pctGroups {
		if len(group) == 1{
			result = append(result, group[0])
		} else if len(group) == 2 {
			// Use two-team tiebreakaer
			if twoTeamCompare(group[0], group[1]) < 0 {
				result = append(result, group[0], group[1])
			} else {
				result = append(result, group[1], group[0])
			}
		} else {
			// Multi-team tiebreaker
			sorted := resolveMultiTeamTie(group, games, twoTeamCompare)
			result = append(result, sorted...)
		}
	}

	// Sort groups by win percentage
	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct
	})

	return result
}

func resolveMultiTeamTie(teams []TeamRecord, games []GameResult, twoTeamCompare func(TeamRecord, TeamRecord) int) []TeamRecord {
	// Try head-to-head sweep first
	sweepWinner := checkHeadToHeadSweep(teams, games)
	if sweepWinner != nil {
		remaining := []TeamRecord{}
		for _, team := range teams {
			if team.TeamID != sweepWinner.TeamID {
				remaining = append(remaining, team)
			}
		}

		result := []TeamRecord{*sweepWinner}
		if len(remaining) > 0 {
			result = append(result, resolveMultiTeamTie(remaining, games, twoTeamCompare)...)
		}
		return result
	}

	// Fall back to sorting with two-team comparisons
	sorted := make([]TeamRecord, len(teams))
	copy(sorted, teams)

	sort.Slice(sorted, func(i, j int) bool {
		return twoTeamCompare(sorted[i], sorted[j]) < 0
	})

	return sorted
}

func compareHeadToHead(teams []TeamRecord, games []GameResult) []TeamRecord {
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
		return []TeamRecord{teamA}
	} else if bWins > aWins {
		return []TeamRecord{teamB}
	}

	return teams
}

func checkHeadToHeadSweep(teams []TeamRecord, games []GameResult) *TeamRecord {
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

func compareCommonGames(teams []TeamRecord, games []GameResult, minCommonGames int) []TeamRecord {
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

	aPct := calculateWinPct(aWins, aLosses, aTies)
	bPct := calculateWinPct(bWins, bLosses, bTies)

	if aPct > bPct {
		return []TeamRecord{teamA}
	} else if bPct > aPct {
		return []TeamRecord{teamB}
	}

	return teams
}

func removeTeam(teams []TeamRecord, teamID int) []TeamRecord {
	var result []TeamRecord
	for _, team := range teams {
		if team.TeamID != teamID {
			result = append(result, team)
		}
	}
	return result
}

func calculateDraftOrder(allTeams []TeamRecord, afc ConferenceStandings, nfc ConferenceStandings) []DraftPick {
	// Get playoff teams (first 7 from each conference)
	playoffTeamIDs := make(map[int]bool)
	for i := 0; i < 7 && i < len(afc.PlayoffSeeds); i++ {
		playoffTeamIDs[afc.PlayoffSeeds[i].Team.TeamID] = true
	}
	for i := 0; i < 7 && i < len(nfc.PlayoffSeeds); i++ {
		playoffTeamIDs[nfc.PlayoffSeeds[i].Team.TeamID] = true
	}

	// Get non-playoff teams
	var nonPlayoffTeams []TeamRecord
	for _, team := range allTeams {
		if !playoffTeamIDs[team.TeamID] {
			nonPlayoffTeams = append(nonPlayoffTeams, team)
		}
	}

	// Sort non-playoff teams by record (worst to best)
	sort.Slice(nonPlayoffTeams, func(i, j int) bool {
		a, b := nonPlayoffTeams[i], nonPlayoffTeams[j]
		if a.WinPct != b.WinPct {
			return a.WinPct < b.WinPct // Worst teams first
		}
		// Tiebreaker: worse point differential picks first
		aDiff := a.PointsFor - a.PointsAgainst
		bDiff := b.PointsFor - b.PointsAgainst

		return aDiff < bDiff
	})

	// Build draft order
	draftOrder := []DraftPick{}
	pickNum := 1

	// Picks 1-18: Non-playoff teams
	for _, team := range nonPlayoffTeams {
		draftOrder = append(draftOrder, DraftPick{
			Pick: pickNum,
			Team: team,
			Reason: "Non-playoff",
		})
		pickNum++
	}

	// Picks 19-32: Playoff teams (ordered by seed, worst to best)
	// Note: In real NFL, this is determined by playoff results
	// For now, we'll just reverse the playoff seeding

	// Combine both conferences' playoff teams and sort by seed (descending)
	var playoffTeams []PlayoffSeed
	playoffTeams = append(playoffTeams, afc.PlayoffSeeds[:7]...)
	playoffTeams = append(playoffTeams, nfc.PlayoffSeeds[:7]...)

	// Sort by seed (7 seed picks before 6 seed, etc.)
	sort.Slice(playoffTeams, func(i, j int) bool {
		// If different seeds, higher seed picks first
		if playoffTeams[i].Seed != playoffTeams[j].Seed {
			return playoffTeams[i].Seed > playoffTeams[j].Seed
		}
		// If same seed, worse record picks first
		return playoffTeams[i].Team.WinPct < playoffTeams[j].Team.WinPct
	})

	for _, seed := range playoffTeams {
		draftOrder = append(draftOrder, DraftPick{
			Pick: pickNum,
			Team: seed.Team,
			Reason: fmt.Sprintf("Playoff seed %d", seed.Seed),
		})
		pickNum++
	}

	return draftOrder
}

func calculateWinPct(wins, losses, ties int) float64 {
	total := wins + losses + ties
	if total == 0 {
		return -1.0 // Treat 0-0 as worse than any actual record
	}
	return (float64(wins) + 0.5*float64(ties)) / float64(total)
}

func calculateGamesBack(leader TeamRecord, team TeamRecord) float64 {
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