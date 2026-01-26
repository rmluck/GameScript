// NBA standings logic

package standings

import (
	"fmt"
	"sort"

	"gamescript/internal/database"
)


type NBATeamRecord struct {
	TeamID int
	TeamCity string
	TeamName string
	TeamAbbr string
	Conference string
	Division string
	Wins int
	Losses int
	HomeWins int
	HomeLosses int
	AwayWins int
	AwayLosses int
	DivisionWins int
	DivisionLosses int
	ConferenceWins int
	ConferenceLosses int
	PointsFor int
	PointsAgainst int
	GamesWithScores int
	WinPct float64
	ConferenceGamesBack float64
	DivisionGamesBack float64
	StrengthOfSchedule float64
	StrengthOfVictory float64
	IsDivisionWinner bool
	LogoURL string
	TeamPrimaryColor string
	TeamSecondaryColor string
}

type NBAStandings struct {
	Eastern NBAConferenceStandings
	Western NBAConferenceStandings
	DraftOrder []NBADraftPick
}

type NBAConferenceStandings struct {
	Divisions map[string][]NBATeamRecord // Keyed by division name
	PlayoffSeeds []NBAPlayoffSeed
}

type NBAPlayoffSeed struct {
	Seed int
	Team NBATeamRecord
	IsDivisionWinner bool
}

type NBADraftPick struct {
	Pick int
	Team NBATeamRecord
}

type NBAGameResult struct {
	GameID int
	HomeTeamID int
	AwayTeamID int
	HomeScore int
	AwayScore int
	Week int
	HasRealScores bool
}

func CalculateNBAStandings(db *database.DB, scenarioID int, seasonID int) (*NBAStandings, error) {
	// Get all teams for the season
	teams, err := getNBATeams(db, seasonID)
	if err != nil {
		return nil, fmt.Errorf("error getting teams: %w", err)
	}

	// Get all game results for the scenario
	games, err := getNBAGameResults(db, scenarioID, seasonID)
	if err != nil {
		return nil, fmt.Errorf("error getting game results: %w", err)
	}

	// Calculate team records
	records := calculateNBATeamRecords(teams, games)

	// Calculate strength metrics
	calculateNBAStrengthMetrics(records, games)

	// Separate by conference
	easternTeams := filterByNBAConference(records, "Eastern")
	westernTeams := filterByNBAConference(records, "Western")

	// Calculate standings for each conference
	easternStandings := calculateNBAConferenceStandings(easternTeams, games)
	westernStandings := calculateNBAConferenceStandings(westernTeams, games)

	// Calculate draft order
	draftOrder := calculateNBADraftOrder(records, easternStandings, westernStandings)

	return &NBAStandings{
		Eastern: easternStandings,
		Western: westernStandings,
		DraftOrder: draftOrder,
	}, nil
}

func getNBATeams(db *database.DB, seasonID int) ([]NBATeamRecord, error) {
	query := `
		SELECT id, city, name, abbreviation, conference, division,
			logo_url, primary_color, secondary_color
		FROM teams
		WHERE season_id = $1
		ORDER BY conference, division, name
	`

	rows, err := db.Query(query, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []NBATeamRecord
	for rows.Next() {
		var team NBATeamRecord
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
			return nil, err
		}
		teams = append(teams, team)
	}
	
	return teams, nil
}

func getNBAGameResults(db *database.DB, scenarioID int, seasonID int) ([]NBAGameResult, error) {
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

	var games []NBAGameResult
	for rows.Next() {
		var game NBAGameResult
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
			return nil, err
		}

		// Priority 1: User has made a pick
		if pickedTeamID != nil {
			// If user provided predicted scores, use those
			if predictedHomeScore != nil && predictedAwayScore != nil {
				game.HomeScore = *predictedHomeScore
				game.AwayScore = *predictedAwayScore
				game.HasRealScores = true
			} else {
				// No predicted scores, use dummy scores based on picked winner
				game.HasRealScores = false
				if *pickedTeamID == game.HomeTeamID {
					game.HomeScore = 1
					game.AwayScore = 0
				} else if *pickedTeamID == game.AwayTeamID {
					game.AwayScore = 1
					game.HomeScore = 0
				} else {
					// Invalid picked team ID
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
			game.HasRealScores = true
			games = append(games, game)
			continue
		}

		// Priority 3: No pick and game not final - skip this game
		// (don't add to games slice)
	}

	return games, nil
}

func calculateNBATeamRecords(teams []NBATeamRecord, games []NBAGameResult) []NBATeamRecord {
	// Initialize record map for easy lookup
	recordMap := make(map[int]*NBATeamRecord)
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
		} else {
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
		}

		if game.HasRealScores {
			// Update points for/against
			homeTeam.PointsFor += game.HomeScore
			homeTeam.PointsAgainst += game.AwayScore
			homeTeam.GamesWithScores++
			awayTeam.PointsFor += game.AwayScore
			awayTeam.PointsAgainst += game.HomeScore
			awayTeam.GamesWithScores++
		}
	}

	// Calculate win percentages
	for _, team := range recordMap {
		team.WinPct = calculateNBAWinPct(team.Wins, team.Losses)
	}

	return teams
}

func calculateNBAStrengthMetrics(teams []NBATeamRecord, games []NBAGameResult) {
	// Initialize team map for easy lookup
	teamMap := make(map[int]*NBATeamRecord)
	for i := range teams {
		teamMap[teams[i].TeamID] = &teams[i]
	}

	// Calculate strength of schedule and strength of victory
	for i := range teams {
		team := &teams[i]

		var opponentTotalWins, opponentTotalLosses int
		var opponentCount int
		var defeatedOpponentWins, defeatedOpponentLosses int
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
			opponentCount++

			// Strength of victory: only defeated opponents
			if teamWon {
				defeatedOpponentWins += opponent.Wins
				defeatedOpponentLosses += opponent.Losses
				defeatedCount++
			}
		}

		// Calculate averages
		if opponentCount > 0 {
			team.StrengthOfSchedule = calculateNBAWinPct(opponentTotalWins, opponentTotalLosses)
		}

		if defeatedCount > 0 {
			team.StrengthOfVictory = calculateNBAWinPct(defeatedOpponentWins, defeatedOpponentLosses)
		}
	}
}

func filterByNBAConference(teams []NBATeamRecord, conference string) []NBATeamRecord {
	var filtered []NBATeamRecord
	for _, team := range teams {
		if team.Conference == conference {
			filtered = append(filtered, team)
		}
	}
	return filtered
}

func calculateNBAConferenceStandings(teams []NBATeamRecord, games []NBAGameResult) NBAConferenceStandings {
	// Group teams by division
	divisions := make(map[string][]NBATeamRecord)
	for _, team := range teams {
		divisions[team.Division] = append(divisions[team.Division], team)
	}

	// Initialize maps for tracking games back and division ranks
	divisionGamesBackMap := make(map[int]float64)
	divisionRankMap := make(map[int]int)

	// Determine division winners first (must be broken before other ties)
	divisionWinners := make(map[string]NBATeamRecord)
	for divName, divTeams := range divisions {
		// Sort division teams with tiebreakers
		sortedDiv := applyNBADivisionTiebreakers(divTeams, games)
		divisionWinners[divName] = sortedDiv[0]

		// Calculate division games back
		divLeader := sortedDiv[0]
		for i := range sortedDiv {
			sortedDiv[i].DivisionGamesBack = calculateNBAGamesBack(divLeader, sortedDiv[i])
			divisionGamesBackMap[sortedDiv[i].TeamID] = sortedDiv[i].DivisionGamesBack
			divisionRankMap[sortedDiv[i].TeamID] = i + 1
		}
		divisions[divName] = sortedDiv
	}

	// Mark division winners
	for i := range teams {
		if dgb, exists := divisionGamesBackMap[teams[i].TeamID]; exists {
			teams[i].DivisionGamesBack = dgb
		}

		for _, winner := range divisionWinners {
			if teams[i].TeamID == winner.TeamID {
				teams[i].IsDivisionWinner = true
				break
			}
		}
	}

	// Apply conference-wide tiebreakers to rank all teams (seeds 1-15)
	rankedTeams := applyNBAConferenceTiebreakers(teams, games, divisionWinners)

	// Calculate conference games back
	conferenceLeader := rankedTeams[0]
	for i := range rankedTeams {
		rankedTeams[i].ConferenceGamesBack = calculateNBAGamesBack(conferenceLeader, rankedTeams[i])
	}

	// Create playoff seeds
	playoffSeeds := []NBAPlayoffSeed{}
	for i, team := range rankedTeams {
		isDivWinner := false
		for _, winner := range divisionWinners {
			if team.TeamID == winner.TeamID {
				isDivWinner = true
				break
			}
		}

		playoffSeeds = append(playoffSeeds, NBAPlayoffSeed{
			Seed: i + 1,
			Team: team,
			IsDivisionWinner: isDivWinner,
		})
	}

	// Update divisions map
	teamRecordMap := make(map[int]NBATeamRecord)
	for _, team := range rankedTeams {
		teamRecordMap[team.TeamID] = team
	}
	for divName, divTeams := range divisions {
		for i := range divTeams {
			if updatedTeam, exists := teamRecordMap[divTeams[i].TeamID]; exists {
				divTeams[i].ConferenceGamesBack = updatedTeam.ConferenceGamesBack
			}
		}
		divisions[divName] = divTeams
	}

	return NBAConferenceStandings{
		Divisions: divisions,
		PlayoffSeeds: playoffSeeds,
	}
}

func applyNBADivisionTiebreakers(teams []NBATeamRecord, games []NBAGameResult) []NBATeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// Group by win percentage (rounded to avoid floating point issues)
	pctGroups := make(map[string][]NBATeamRecord)
	for _, team := range teams {
		key := fmt.Sprintf("%.6f", team.WinPct)
		pctGroups[key] = append(pctGroups[key], team)
	}

	var result []NBATeamRecord
	for _, group := range pctGroups {
		if len(group) == 1 {
			result = append(result, group[0])
		} else if len(group) == 2 {
			result = append(result, resolveNBATwoTeamTie(group, games, nil, true)...)
		} else {
			result = append(result, resolveNBAMultiTeamTie(group, games, nil, true)...)
		}
	}

	// Sort by win percentage
	sort.Slice(result, func(i, j int) bool {
		return result[i].WinPct > result[j].WinPct
	})

	return result
}

func applyNBAConferenceTiebreakers(teams []NBATeamRecord, games []NBAGameResult, divisionWinners map[string]NBATeamRecord) []NBATeamRecord {
	if len(teams) <= 1 {
		return teams
	}

	// First, sort by division rank to preserve division tiebreaker results
    // Teams are already marked with IsDivisionWinner from calculateNBAConferenceStandings
    sort.Slice(teams, func(i, j int) bool {
        // If same division, preserve their division order
        if teams[i].Division == teams[j].Division {
            // Both division winners or both not - maintain current order
            if teams[i].IsDivisionWinner == teams[j].IsDivisionWinner {
                // Use division games back to preserve division ranking
                if teams[i].DivisionGamesBack != teams[j].DivisionGamesBack {
                    return teams[i].DivisionGamesBack < teams[j].DivisionGamesBack
                }
            }
            // Division winner always comes first
            return teams[i].IsDivisionWinner
        }
        // Different divisions - compare win percentage
        return teams[i].WinPct > teams[j].WinPct
    })

	// Group by win percentage (rounded to avoid floating point issues)
	pctGroups := make(map[string][]NBATeamRecord)
    groupOrder := []string{} // Track order of groups
    for _, team := range teams {
        key := fmt.Sprintf("%.6f", team.WinPct)
        if _, exists := pctGroups[key]; !exists {
            groupOrder = append(groupOrder, key)
        }
        pctGroups[key] = append(pctGroups[key], team)
    }

	// Resolve ties within each win percentage group
	var result []NBATeamRecord
	sort.Slice(groupOrder, func(i, j int) bool {
        return groupOrder[i] > groupOrder[j]
    })
	for _, key := range groupOrder {
        group := pctGroups[key]
        if len(group) == 1 {
            result = append(result, group[0])
        } else if len(group) == 2 {
            result = append(result, resolveNBATwoTeamTie(group, games, divisionWinners, false)...)
        } else {
            result = append(result, resolveNBAMultiTeamTie(group, games, divisionWinners, false)...)
        }
    }

	return result
}

func resolveNBATwoTeamTie(teams []NBATeamRecord, games []NBAGameResult, divisionWinners map[string]NBATeamRecord, inDivision bool) []NBATeamRecord {
	a, b := teams[0], teams[1]

	// Step 1: Head-to-head
	h2h := compareNBAHeadToHead([]NBATeamRecord{a, b}, games)
	if len(h2h) == 1 {
		if h2h[0].TeamID == a.TeamID {
			return []NBATeamRecord{a, b}
		}
		return []NBATeamRecord{b, a}
	}

	// Step 2: Division winner (if not already determining division winner)
	if !inDivision && divisionWinners != nil {
		aIsDivWinner := false
		bIsDivWinner := false
		for _, winner := range divisionWinners {
			if winner.TeamID == a.TeamID {
				aIsDivWinner = true
			}
			if winner.TeamID == b.TeamID {
				bIsDivWinner = true
			}
		}
		if aIsDivWinner != bIsDivWinner {
			if aIsDivWinner {
				return []NBATeamRecord{a, b}
			}
			return []NBATeamRecord{b, a}
		}
	}

	// Step 3: Division win percentage (only if same division)
	if a.Division == b.Division {
		aDivPct := calculateNBAWinPct(a.DivisionWins, a.DivisionLosses)
		bDivPct := calculateNBAWinPct(b.DivisionWins, b.DivisionLosses)
		if aDivPct != bDivPct {
			if aDivPct > bDivPct {
				return []NBATeamRecord{a, b}
			}
			return []NBATeamRecord{b, a}
		}
	}

	// Step 4: Conference win percentage
	aConfPct := calculateNBAWinPct(a.ConferenceWins, a.ConferenceLosses)
	bConfPct := calculateNBAWinPct(b.ConferenceWins, b.ConferenceLosses)
	if aConfPct != bConfPct {
		if aConfPct > bConfPct {
			return []NBATeamRecord{a, b}
		}
		return []NBATeamRecord{b, a}
	}

	// Step 5: Point differential
	aPointDiff := a.PointsFor - a.PointsAgainst
	bPointDiff := b.PointsFor - b.PointsAgainst
	if aPointDiff != bPointDiff {
		if aPointDiff > bPointDiff {
			return []NBATeamRecord{a, b}
		}
		return []NBATeamRecord{b, a}
	}

	// Random drawing - use TeamID for consistency
	if a.TeamID < b.TeamID {
		return []NBATeamRecord{a, b}
	}
	return []NBATeamRecord{b, a}
}

func resolveNBAMultiTeamTie(teams []NBATeamRecord, games []NBAGameResult, divisionWinners map[string]NBATeamRecord, inDivision bool) []NBATeamRecord {
	if len(teams) <= 1 {
		return teams
	}
	if len(teams) == 2 {
		return resolveNBATwoTeamTie(teams, games, divisionWinners, inDivision)
	}

	// Step 1: Division winner (if not already determining division winner)
	if !inDivision && divisionWinners != nil {
		divWinners := []NBATeamRecord{}
		nonWinners := []NBATeamRecord{}

		for _, team := range teams {
			isDivWinner := false
			for _, winner := range divisionWinners {
				if team.TeamID == winner.TeamID {
					isDivWinner = true
					break
				}
			}
			if isDivWinner {
				divWinners = append(divWinners, team)
			} else {
				nonWinners = append(nonWinners, team)
			}
		}

		// If some are division winners and some aren't, separate them
		if len(divWinners) > 0 && len(nonWinners) > 0 {
			var result []NBATeamRecord
			if len(divWinners) == 1 {
				result = append(result, divWinners[0])
			} else {
				result = append(result, resolveNBAMultiTeamTie(divWinners, games, divisionWinners, inDivision)...)
			}

			if len(nonWinners) == 1 {
				result = append(result, nonWinners[0])
			} else {
				result = append(result, resolveNBAMultiTeamTie(nonWinners, games, divisionWinners, inDivision)...)
			}

			return result
		}
	}

	// Step 2: Head-to-head win percentage
	h2hWinner := findNBAHeadToHeadWinner(teams, games)
	if h2hWinner != nil {
		remaining := removeNBATeam(teams, h2hWinner.TeamID)
		result := []NBATeamRecord{*h2hWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNBAMultiTeamTie(remaining, games, divisionWinners, inDivision)...)
		}
		return result
	}

	// Step 3: Division win percentage (only if all same division)
	allSameDivision := true
	firstDiv := teams[0].Division
	for _, team := range teams[1:] {
		if team.Division != firstDiv {
			allSameDivision = false
			break
		}
	}
	if allSameDivision {
		divWinner := findBestNBADivisionRecord(teams)
		if divWinner != nil {
			remaining := removeNBATeam(teams, divWinner.TeamID)
			result := []NBATeamRecord{*divWinner}
			if len(remaining) > 0 {
				result = append(result, resolveNBAMultiTeamTie(remaining, games, divisionWinners, inDivision)...)
			}
			return result
		}
	}

	// Step 4: Conference win percentage
	confWinner := findBestNBAConferenceRecord(teams)
	if confWinner != nil {
		remaining := removeNBATeam(teams, confWinner.TeamID)
		result := []NBATeamRecord{*confWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNBAMultiTeamTie(remaining, games, divisionWinners, inDivision)...)
		}
		return result
	}

	// Step 5: Point differential
	pointDiffWinner := findBestNBAPointDifferential(teams)
	if pointDiffWinner != nil {
		remaining := removeNBATeam(teams, pointDiffWinner.TeamID)
		result := []NBATeamRecord{*pointDiffWinner}
		if len(remaining) > 0 {
			result = append(result, resolveNBAMultiTeamTie(remaining, games, divisionWinners, inDivision)...)
		}
		return result
	}

	// Random drawing - use TeamID for consistency
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].TeamID < teams[j].TeamID
	})

	return teams
}

func compareNBAHeadToHead(teams []NBATeamRecord, games []NBAGameResult) []NBATeamRecord {
	if len(teams) != 2 {
		return teams
	}

	teamA, teamB := teams[0], teams[1]
	aWins, bWins := 0, 0

	for _, game := range games {
		isMatchup := (game.HomeTeamID == teamA.TeamID && game.AwayTeamID == teamB.TeamID) || (game.HomeTeamID == teamB.TeamID && game.AwayTeamID == teamA.TeamID)

		if !isMatchup {
			continue
		}

		var winner int
		if game.HomeScore > game.AwayScore {
			winner = game.HomeTeamID
		} else {
			winner = game.AwayTeamID
		}

		if winner == teamA.TeamID {
			aWins++
		} else if winner == teamB.TeamID {
			bWins++
		}
	}

	if aWins > bWins {
		return []NBATeamRecord{teamA}
	} else if bWins > aWins {
		return []NBATeamRecord{teamB}
	}

	return teams
}

func findNBAHeadToHeadWinner(teams []NBATeamRecord, games []NBAGameResult) *NBATeamRecord {
	type h2hRecord struct {
		teamID  int
		wins    int
		losses  int
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
		} else {
			records[game.AwayTeamID].wins++
			records[game.HomeTeamID].losses++
		}
	}

	var bestPct float64 = -1.0
	var bestTeamID int
	var tie bool

	for _, rec := range records {
		pct := calculateNBAWinPct(rec.wins, rec.losses)
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

func findBestNBADivisionRecord(teams []NBATeamRecord) *NBATeamRecord {
	var bestPct float64 = -1.0
	var bestTeam *NBATeamRecord
	var tie bool

	for i := range teams {
		pct := calculateNBAWinPct(teams[i].DivisionWins, teams[i].DivisionLosses)
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

func findBestNBAConferenceRecord(teams []NBATeamRecord) *NBATeamRecord {
	var bestPct float64 = -1.0
	var bestTeam *NBATeamRecord
	var tie bool

	for i := range teams {
		pct := calculateNBAWinPct(teams[i].ConferenceWins, teams[i].ConferenceLosses)
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

func findBestNBAPointDifferential(teams []NBATeamRecord) *NBATeamRecord {
	var bestDiff int = -1 << 31
	var bestTeam *NBATeamRecord
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

func removeNBATeam(teams []NBATeamRecord, teamID int) []NBATeamRecord {
	var result []NBATeamRecord
	for _, team := range teams {
		if team.TeamID != teamID {
			result = append(result, team)
		}
	}
	return result
}

func calculateNBADraftOrder(allTeams []NBATeamRecord, eastern NBAConferenceStandings, western NBAConferenceStandings) []NBADraftPick {
	// Get playoff teams (top 10 from each conference)
	playoffTeamIDs := make(map[int]bool)
	for _, seed := range eastern.PlayoffSeeds {
		if seed.Seed <= 10 {
			playoffTeamIDs[seed.Team.TeamID] = true
		}
	}
	for _, seed := range western.PlayoffSeeds {
		if seed.Seed <= 10 {
			playoffTeamIDs[seed.Team.TeamID] = true
		}
	}

	// Get non-playoff teams
	var nonPlayoffTeams []NBATeamRecord
	for _, team := range allTeams {
		if !playoffTeamIDs[team.TeamID] {
			nonPlayoffTeams = append(nonPlayoffTeams, team)
		}
	}

	// Sort non-playoff teams by record (worst to best)
	sort.Slice(nonPlayoffTeams, func(i, j int) bool {
		if nonPlayoffTeams[i].WinPct != nonPlayoffTeams[j].WinPct {
			return nonPlayoffTeams[i].WinPct < nonPlayoffTeams[j].WinPct
		}
		// Tiebreaker: Point differential
		diffI := nonPlayoffTeams[i].PointsFor - nonPlayoffTeams[i].PointsAgainst
		diffJ := nonPlayoffTeams[j].PointsFor - nonPlayoffTeams[j].PointsAgainst
		return diffI < diffJ
	})

	// Build draft order
	draftOrder := []NBADraftPick{}
	pickNum := 1

	// Picks 1-10: Non-playoff teams
	for _, team := range nonPlayoffTeams {
		draftOrder = append(draftOrder, NBADraftPick{
			Pick: pickNum,
			Team: team,
		})
		pickNum++
	}

	// Picks 11-30: Playoff teams (by seed, worst to best)
	// TODO: Update as playoffs progress
	playoffTeams := make([]NBAPlayoffSeed, 0)

	for _, seed := range eastern.PlayoffSeeds {
		if seed.Seed <= 10 {
			playoffTeams = append(playoffTeams, seed)
		}
	}
	for _, seed := range western.PlayoffSeeds {
		if seed.Seed <= 10 {
			playoffTeams = append(playoffTeams, seed)
		}
	}

	sort.Slice(playoffTeams, func(i, j int) bool {
		return playoffTeams[i].Seed > playoffTeams[j].Seed
	})

	for _, seed := range playoffTeams {
		draftOrder = append(draftOrder, NBADraftPick{
			Pick: pickNum,
			Team: seed.Team,
		})
		pickNum++
	}

	return draftOrder
}

func calculateNBAWinPct(wins int, losses int) float64 {
	total := wins + losses
	if total == 0 {
		return 0.0
	}
	return float64(wins) / float64(total)
}

func calculateNBAGamesBack(leader NBATeamRecord, team NBATeamRecord) float64 {
	if leader.TeamID == team.TeamID {
		return 0.0
	}

	// Formula: ((Leader Wins - Team Wins) + (Team Losses - Leader Losses)) / 2
	gamesBack := float64((leader.Wins - team.Wins) + (team.Losses - leader.Losses)) / 2.0
	return gamesBack
}