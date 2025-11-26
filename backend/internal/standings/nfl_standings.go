package standings

import (
	"fmt"
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
	IsPostseason	bool
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
			id, city, name, abbreviation, conference, division
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
            game.id, game.home_team_id, game.away_team_id, game.week, game.is_postseason,
            COALESCE(game.home_score, pick.predicted_home_score) AS home_score,
            COALESCE(game.away_score, pick.predicted_away_score) AS away_score,
            CASE
                WHEN game.status = 'final' THEN 'final'
                WHEN pick.picked_team_id IS NOT NULL THEN 'predicted'
                ELSE 'unpicked'
            END AS result_type
        FROM games game
        LEFT JOIN picks pick ON game.id = pick.game_id AND pick.scenario_id = $1
        WHERE game.season_id = $2
        AND game.is_postseason = false
        AND (
            game.status = 'final'
            OR pick.picked_team_id IS NOT NULL
        )
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
		var resultType string
		var homeScore, awayScore *int

		err := rows.Scan(
			&game.GameID,
			&game.HomeTeamID,
			&game.AwayTeamID,
			&game.Week,
			&game.IsPostseason,
			&homeScore,
			&awayScore,
			&resultType,
		)
		if err != nil {
			continue
		}

		// Skip games without results
		if homeScore == nil || awayScore == nil {
			continue
		}

		game.HomeScore = *homeScore
		game.AwayScore = *awayScore
		games = append(games, game)
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

		// Determine winner
		if game.HomeScore > game.AwayScore {
			// Home team wins
			homeTeam.Wins++
			awayTeam.Losses++

			// Check if division game
			if homeTeam.Division == awayTeam.Division {
				homeTeam.DivisionWins++
				awayTeam.DivisionLosses++
			}

			// Check if conference game
			if homeTeam.Conference == awayTeam.Conference {
				homeTeam.ConferenceWins++
				awayTeam.ConferenceLosses++
			}
		} else if game.AwayScore > game.HomeScore {
			// Away team wins
			awayTeam.Wins++
			homeTeam.Losses++

			if homeTeam.Division == awayTeam.Division {
				awayTeam.DivisionWins++
				homeTeam.DivisionLosses++
			}

			if homeTeam.Conference == awayTeam.Conference {
				awayTeam.ConferenceWins++
				homeTeam.ConferenceLosses++
			}
		} else {
			// Tie
			homeTeam.Ties++
			awayTeam.Ties++

			if homeTeam.Division == awayTeam.Division {
				homeTeam.DivisionTies++
				awayTeam.DivisionTies++
			}

			if homeTeam.Conference == awayTeam.Conference {
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
	divisionWinners = applyConferenceTiebreakers(divisionWinners, games)

	// Rank non-division winners (seeds 5-16)
	nonWinners = applyConferenceTiebreakers(nonWinners, games)

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

// Only two-team tiebreakers implemented so far
// TODO: Expand to multi-team tiebreakers
func applyDivisionTiebreakers(teams []TeamRecord, games []GameResult) []TeamRecord {
	sort.Slice(teams, func(i, j int) bool {
		a, b := teams[i], teams[j]

		// 1. Win percentage
		if a.WinPct != b.WinPct {
			return a.WinPct > b.WinPct
		}

		// 2. Head-to-head record (if they played)
		h2hWinner := getHeadToHeadWinner(a, b, games)
		if h2hWinner != 0 {
			return h2hWinner == a.TeamID
		}

		// 3. Division record
		aDivPct := calculateWinPct(a.DivisionWins, a.DivisionLosses, a.DivisionTies)
		bDivPct := calculateWinPct(b.DivisionWins, b.DivisionLosses, b.DivisionTies)
		if aDivPct != bDivPct {
			return aDivPct > bDivPct
		}

		// 4. Common games record - Skipped right now

		// 5. Conference record
		aConfPct := calculateWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
		bConfPct := calculateWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
		if aConfPct != bConfPct {
			return aConfPct > bConfPct
		}

		// 6. Strength of victory - Skipped right now

		// 7. Strength of schedule - Skipped right now

		// 8. Point differential
		aDiff := a.PointsFor - a.PointsAgainst
		bDiff := b.PointsFor - b.PointsAgainst
		if aDiff != bDiff {
			return aDiff > bDiff
		}

		// 9. Points scored
		if a.PointsFor != b.PointsFor {
			return a.PointsFor > b.PointsFor
		}

		// 10. Points allowed (fewer is better)
		if a.PointsAgainst != b.PointsAgainst {
			return a.PointsAgainst < b.PointsAgainst
		}

		// 11. Coin toss - Skipped for now
		return a.TeamID < b.TeamID // Arbitrary but consistent
	})

	return teams
}

// Only two-team tiebreakers implemented so far
// TODO: Expand to multi-team tiebreakers
func applyConferenceTiebreakers(teams []TeamRecord, games []GameResult) []TeamRecord {
	sort.Slice(teams, func(i, j int) bool {
		a, b := teams[i], teams[j]

		// 1. Win percentage
		if a.WinPct != b.WinPct {
			return a.WinPct > b.WinPct
		}

		// 2. Head-to-head record (if they played)
		h2hWinner := getHeadToHeadWinner(a, b, games)
		if h2hWinner != 0 {
			return h2hWinner == a.TeamID
		}

		// 3. Conference record
		aConfPct := calculateWinPct(a.ConferenceWins, a.ConferenceLosses, a.ConferenceTies)
		bConfPct := calculateWinPct(b.ConferenceWins, b.ConferenceLosses, b.ConferenceTies)
		if aConfPct != bConfPct {
			return aConfPct > bConfPct
		}

		// 4. Common games record - Skipped right now

		// 5. Strength of victory - Skipped right now

		// 6. Strength of schedule - Skipped right now

		// 7. Point differential
		aDiff := a.PointsFor - a.PointsAgainst
		bDiff := b.PointsFor - b.PointsAgainst
		if aDiff != bDiff {
			return aDiff > bDiff
		}

		// 8. Points scored
		if a.PointsFor != b.PointsFor {
			return a.PointsFor > b.PointsFor
		}

		// 9. Points allowed (fewer is better)
		if a.PointsAgainst != b.PointsAgainst {
			return a.PointsAgainst < b.PointsAgainst
		}

		// 10. Coin toss - Skipped for now
		return a.TeamID < b.TeamID // Arbitrary but consistent
	})

	return teams
}

func getHeadToHeadWinner(teamA TeamRecord, teamB TeamRecord, games []GameResult) int {
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
		return teamA.TeamID
	} else if bWins > aWins {
		return teamB.TeamID
	}
	return 0
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func calculateWinPct(wins, losses, ties int) float64 {
	total := wins + losses + ties
	if total == 0 {
		return 0
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