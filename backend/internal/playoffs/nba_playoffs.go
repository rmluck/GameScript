// NBA playoffs generation and management

package playoffs

import (
	"database/sql"
	"fmt"

	"gamescript/internal/database"
	"gamescript/internal/standings"
)


const (
	RoundPlayInA = 1 // 7v8 and 9v10 games
	RoundPlayInB = 2 // Winner 9v10 vs Loser 7v8
	RoundConferenceQuarterfinals = 3
	RoundConferenceSemifinals = 4
	RoundConferenceFinals = 5
	RoundNBAFinals = 6
)

type NBAPlayoffGenerator struct {
	db *database.DB
}

func NewNBAPlayoffGenerator(db *database.DB) *NBAPlayoffGenerator {
	return &NBAPlayoffGenerator{db: db}
}

func (pg *NBAPlayoffGenerator) CheckAndEnableNBAPlayoffs(scenarioID int, seasonID int) (bool, error) {
	// Count total regular season games
	var totalGames int
	err := pg.db.Conn.QueryRow(`
		SELECT COUNT(*)
		FROM games
		WHERE season_id = $1
	`, seasonID).Scan(&totalGames)
	if err != nil {
		return false, err
	}

	// Count games that are either completed or picked
	var completedOrPickedGames int
	err = pg.db.Conn.QueryRow(`
		SELECT COUNT(DISTINCT g.id)
		FROM games g
		LEFT JOIN picks p ON g.id = p.game_id AND p.scenario_id = $1
		WHERE g.season_id = $2
		AND (
			(p.picked_team_id IS NOT NULL) OR
			(g.status = 'final' AND g.home_score IS NOT NULL AND g.away_score IS NOT NULL)
		)
	`, scenarioID, seasonID).Scan(&completedOrPickedGames)
	if err != nil {
		return false, err
	}

	// Check if all games are completed or picked
	return totalGames == completedOrPickedGames, nil
}

func (pg *NBAPlayoffGenerator) GenerateNBAPlayInRoundA(scenarioID int, seasonID int) error {
	// Get current standings
	nbaStandings, err := standings.CalculateNBAStandings(pg.db, scenarioID, seasonID)
	if err != nil {
		return fmt.Errorf("failed to calculate standings: %v", err)
	}

	// Get or create playoff state
	playoffStateID, err := pg.getOrCreateNBAPlayoffState(scenarioID)
	if err != nil {
		return err
	}

	// Clear existing play-in matchups
	_, err = pg.db.Conn.Exec(`
		DELETE FROM playoff_matchups
		WHERE playoff_state_id = $1 and round = $2
	`, playoffStateID, RoundPlayInA)
	if err != nil {
		return err
	}

	// Generate Eastern Conference play-in Round A (7v8 and 9v10)
	err = pg.generateNBAConferencePlayInA(playoffStateID, "Eastern", nbaStandings.Eastern.PlayoffSeeds)
	if err != nil {
		return err
	}

	// Generate Western Conference play-in Round A (7v8 and 9v10)
	err = pg.generateNBAConferencePlayInA(playoffStateID, "Western", nbaStandings.Western.PlayoffSeeds)
	if err != nil {
		return err
	}

	// Update playoff state
	_, err = pg.db.Conn.Exec(`
		UPDATE playoff_states
		SET is_enabled = true, current_round = $1,
		updated_at = NOW()
		WHERE id = $2
	`, RoundPlayInA, playoffStateID)

	return err
}

func (pg *NBAPlayoffGenerator) generateNBAConferencePlayInA(playoffStateID int, conference string, seeds []standings.NBAPlayoffSeed) error {
	if len(seeds) < 10 {
		return fmt.Errorf("not enough playoff seeds for conference %s", conference)
	}

	// Create 7v8 matchup
	seed7 := seeds[6]
	seed8 := seeds[7]
	_, err := pg.db.Conn.Exec(`
		INSERT INTO playoff_matchups (
			playoff_state_id, round, matchup_order, conference, higher_seed_team_id, lower_seed_team_id, higher_seed, lower_seed, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
	`, playoffStateID, RoundPlayInA, 1, conference, seed7.Team.TeamID, seed8.Team.TeamID, seed7.Seed, seed8.Seed)
	if err != nil {
		return err
	}

	// Create 9v10 matchup
	seed9 := seeds[8]
	seed10 := seeds[9]
	_, err = pg.db.Conn.Exec(`
		INSERT INTO playoff_matchups (
			playoff_state_id, round, matchup_order, conference, higher_seed_team_id, lower_seed_team_id, higher_seed, lower_seed, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
	`, playoffStateID, RoundPlayInA, 2, conference, seed9.Team.TeamID, seed10.Team.TeamID, seed9.Seed, seed10.Seed)

	return err
}

// Generates the final play-in game (winner 9v10 vs loser 7v8)
func (pg *NBAPlayoffGenerator) GenerateNBAPlayInRoundB(scenarioID int) error {
	playoffStateID, err := pg.getNBAPlayoffStateID(scenarioID)
	if err != nil {
		return err
	}

	// Get winners and losers from Round A
	playInResults, err := pg.getNBAPlayInRoundAResults(playoffStateID)
	if err != nil {
		return err
	}

	// Clear existing Round B matchups
	_, err = pg.db.Conn.Exec(`
		DELETE FROM playoff_matchups
		WHERE playoff_state_id = $1 AND round = $2
	`, playoffStateID, RoundPlayInB)
	if err != nil {
		return err
	}

	// Generate matchup for each conference
	for conference, results := range playInResults {
		winner9v10 := results.Winner9v10
		loser7v8 := results.Loser7v8

		// Higher seed is the loser of 7v8 (because they were seeded higher)
		_, err = pg.db.Conn.Exec(`
			INSERT INTO playoff_matchups (
				playoff_state_id, round, matchup_order, conference,
				higher_seed_team_id, lower_seed_team_id, higher_seed, lower_seed, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
		`, playoffStateID, RoundPlayInB, 1, conference, loser7v8.TeamID, winner9v10.TeamID, loser7v8.Seed, winner9v10.Seed)
		if err != nil {
			return err
		}
	}

	// Update playoff state
	_, err = pg.db.Conn.Exec(`
		UPDATE playoff_states
		SET current_round = $1, updated_at = NOW()
		WHERE id = $2
	`, RoundPlayInB, playoffStateID)

	return err
}

type PlayInRoundAResult struct {
	Winner7v8 TeamSeed
	Loser7v8 TeamSeed
	Winner9v10 TeamSeed
	Loser9v10 TeamSeed
}

func (pg *NBAPlayoffGenerator) getNBAPlayInRoundAResults(playoffStateID int) (map[string]PlayInRoundAResult, error) {
	// Fetch all Round A matchups
	query := `
		SELECT
			conference, matchup_order,
			higher_seed_team_id, lower_seed_team_id,
			higher_seed, lower_seed,
			picked_team_id
		FROM playoff_matchups
		WHERE playoff_state_id = $1 AND round = $2
		ORDER BY conference, matchup_order
	`
	rows, err := pg.db.Query(query, playoffStateID, RoundPlayInA)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]PlayInRoundAResult)
	for rows.Next() {
		var conference string
		var matchupOrder int
		var higherSeedTeamID, lowerSeedTeamID, higherSeed, lowerSeed int
		var pickedTeamID *int

		err := rows.Scan(&conference, &matchupOrder, &higherSeedTeamID, &lowerSeedTeamID, &higherSeed, &lowerSeed, &pickedTeamID)
		if err != nil {
			return nil, err
		}

		if pickedTeamID == nil {
			return nil, fmt.Errorf("play-in round A not complete for %s conference", conference)
		}

		result := results[conference]

		// Determine winner and loser
		var winner, loser TeamSeed
		if *pickedTeamID == higherSeedTeamID {
			winner = TeamSeed{TeamID: higherSeedTeamID, Seed: higherSeed}
			loser = TeamSeed{TeamID: lowerSeedTeamID, Seed: lowerSeed}
		} else {
			winner = TeamSeed{TeamID: lowerSeedTeamID, Seed: lowerSeed}
			loser = TeamSeed{TeamID: higherSeedTeamID, Seed: higherSeed}
		}

		if matchupOrder == 1 { // 7v8 game
			result.Winner7v8 = winner
			result.Loser7v8 = loser	
		} else { // 9v10 game
			result.Winner9v10 = winner
			result.Loser9v10 = loser
		}

		results[conference] = result
	}

	return results, nil
}

func (pg *NBAPlayoffGenerator) GenerateNBAConferenceQuarterfinals(scenarioID int, seasonID int) error {
	playoffStateID, err := pg.getNBAPlayoffStateID(scenarioID)
	if err != nil {
		return err
	}

	// Get current standings for seeds 1-6
	nbaStandings, err := standings.CalculateNBAStandings(pg.db, scenarioID, seasonID)
	if err != nil {
		return fmt.Errorf("failed to calculate standings: %w", err)
	}

	// Get play-in tournament results for seeds 7-8
	playInResults, err := pg.getNBAPlayInFinalSeeds(playoffStateID)
	if err != nil {
		return err
	}

	// Clear existing quarterfinal matchups
	_, err = pg.db.Conn.Exec(`
		DELETE FROM playoff_series
		WHERE playoff_state_id = $1 AND round = $2
	`, playoffStateID, RoundConferenceQuarterfinals)
	if err != nil {
		return err
	}

	// Generate Eastern Conference quarterfinals
	err = pg.generateNBAConferenceQuarterfinalsForConference(
		playoffStateID, "Eastern", nbaStandings.Eastern.PlayoffSeeds, playInResults["Eastern"])
	if err != nil {
		return err
	}

	// Generate Western Conference quarterfinals
	err = pg.generateNBAConferenceQuarterfinalsForConference(
		playoffStateID, "Western", nbaStandings.Western.PlayoffSeeds, playInResults["Western"])
	if err != nil {
		return err
	}

	// Update playoff state
	_, err = pg.db.Conn.Exec(`
		UPDATE playoff_states
		SET current_round = $1, updated_at = NOW()
		WHERE id = $2
	`, RoundConferenceQuarterfinals, playoffStateID)

	return err
}

func (pg *NBAPlayoffGenerator) generateNBAConferenceQuarterfinalsForConference(playoffStateID int, conference string, seeds []standings.NBAPlayoffSeed, playInSeeds PlayInFinalSeeds) error {
	// Build full playoff bracket
	bracket := make([]TeamSeed, 8)

	// Seeds 1-6 from regular season standings
	for i := 0; i < 6; i++ {
		bracket[i] = TeamSeed{
			TeamID: seeds[i].Team.TeamID,
			Seed: seeds[i].Seed,
		}
	}

	// Seeds 7-8 from play-in winners
	bracket[6] = playInSeeds.Seed7
	bracket[7] = playInSeeds.Seed8

	// Create matchups
	matchups := []struct {
		higher int
		lower int
		order int
	}{
		{0, 7, 1}, // 1v8
		{1, 6, 2}, // 2v7
		{2, 5, 3}, // 3v6
		{3, 4, 4}, // 4v5
	}

	for _, m := range matchups {
		higherSeed := bracket[m.higher]
		lowerSeed := bracket[m.lower]

		// Create the series
		var seriesID int
		err := pg.db.Conn.QueryRow(`
			INSERT INTO playoff_series (
				playoff_state_id, round, series_order, conference,
				higher_seed_team_id, lower_seed_team_id, higher_seed, lower_seed, best_of, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 7, 'pending')
			RETURNING id
		`, playoffStateID, RoundConferenceQuarterfinals, m.order, conference, higherSeed.TeamID, lowerSeed.TeamID, higherSeed.Seed, lowerSeed.Seed).Scan(&seriesID)
		if err != nil {
			return err
		}
	}

	return nil
}

type PlayInFinalSeeds struct {
	Seed7 TeamSeed
	Seed8 TeamSeed
}

func (pg *NBAPlayoffGenerator) getNBAPlayInFinalSeeds(playoffStateID int) (map[string]PlayInFinalSeeds, error) {
	// Get Round A results
	roundAResults, err := pg.getNBAPlayInRoundAResults(playoffStateID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT conference, picked_team_id,
			higher_seed_team_id, lower_seed_team_id,
			higher_seed, lower_seed
		FROM playoff_matchups
		WHERE playoff_state_id = $1 AND round = $2
	`
	rows, err := pg.db.Query(query, playoffStateID, RoundPlayInB)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]PlayInFinalSeeds)
	for rows.Next() {
		var conference string
		var pickedTeamID *int
		var higherSeedTeamID, lowerSeedTeamID, higherSeed, lowerSeed int

		err := rows.Scan(&conference, &pickedTeamID, &higherSeedTeamID, &lowerSeedTeamID, &higherSeed, &lowerSeed)
		if err != nil {
			return nil, err
		}

		if pickedTeamID == nil {
			return nil, fmt.Errorf("play-in round B not complete for %s conference", conference)
		}

		var seed8Winner TeamSeed
		if *pickedTeamID == higherSeedTeamID {
			seed8Winner = TeamSeed{TeamID: higherSeedTeamID, Seed: 8}
		} else {
			seed8Winner = TeamSeed{TeamID: lowerSeedTeamID, Seed: 8}
		}

		results[conference] = PlayInFinalSeeds{
			Seed7: roundAResults[conference].Winner7v8,
			Seed8: seed8Winner,
		}
	}

	return results, nil
}

func (pg *NBAPlayoffGenerator) GenerateNBANextRound(scenarioID int, currentRound int) error {
	playoffStateID, err := pg.getNBAPlayoffStateID(scenarioID)
	if err != nil {
		return err
	}

	// Get series winners from current round
	winners, err := pg.getNBASeriesWinners(playoffStateID, currentRound)
	if err != nil {
		return err
	}

	nextRound := currentRound + 1

	// Clear existing series for next round
	_, err = pg.db.Conn.Exec(`
		DELETE FROM playoff_series WHERE playoff_state_id = $1 AND round = $2
	`, playoffStateID, nextRound)
	if err != nil {
		return err
	}

	// Generate series for next round
	var genErr error
	switch nextRound {
	case RoundConferenceSemifinals:
		genErr = pg.generateNBAConferenceSemifinals(playoffStateID, winners)
	case RoundConferenceFinals:
		genErr = pg.generateNBAConferenceFinals(playoffStateID, winners)
	case RoundNBAFinals:
		genErr = pg.generateNBAFinals(playoffStateID, winners)
	default:
		return fmt.Errorf("invalid next round: %d", nextRound)
	}
	if genErr != nil {
		return genErr
	}

	// Update playoff state
	_, err = pg.db.Conn.Exec(`
		UPDATE playoff_states
		SET current_round = $1, updated_at = NOW()
		WHERE id = $2
	`, nextRound, playoffStateID)

	return err
}

func (pg *NBAPlayoffGenerator) generateNBAConferenceSemifinals(playoffStateID int, conferenceQuarterfinalsWinners map[string][]TeamSeed) error {
	// For each conference, match up the quarterfinal winners
	for conference, winners := range conferenceQuarterfinalsWinners {
		if len(winners) != 4 {
			return fmt.Errorf("expected 4 quarterfinal winners for %s conference, got %d", conference, len(winners))
		}

		matchups := []struct {
			team1 TeamSeed
			team2 TeamSeed
			order int
		}{
			{winners[0], winners[3], 1}, // 1v8 winner vs 4v5 winner
			{winners[1], winners[2], 2}, // 2v7 winner vs 3v6 winner
		}

		for _, m := range matchups {
			var higherSeed, lowerSeed TeamSeed
			if m.team1.Seed < m.team2.Seed {
				higherSeed = m.team1
				lowerSeed = m.team2
			} else {
				higherSeed = m.team2
				lowerSeed = m.team1
			}

			// Create series
			var seriesID int
			err := pg.db.Conn.QueryRow(`
				INSERT INTO playoff_series (
					playoff_state_id, round, series_order, conference,
					higher_seed_team_id, lower_seed_team_id, higher_seed, lower_seed, best_of, status
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 7, 'pending')
				RETURNING id
			`, playoffStateID, RoundConferenceSemifinals, m.order, conference, higherSeed.TeamID, lowerSeed.TeamID, higherSeed.Seed, lowerSeed.Seed).Scan(&seriesID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pg *NBAPlayoffGenerator) generateNBAConferenceFinals(playoffStateID int, conferenceSemifinalsWinners map[string][]TeamSeed) error {
	// For each conference, match up the semifinal winners
	for conference, winners := range conferenceSemifinalsWinners {
		if len(winners) != 2 {
			return fmt.Errorf("expected 2 semifinals winners for %s conference, got %d", conference, len(winners))
		}

		var higherSeed, lowerSeed TeamSeed
		if winners[0].Seed < winners[1].Seed {
			higherSeed = winners[0]
			lowerSeed = winners[1]
		} else {
			higherSeed = winners[1]
			lowerSeed = winners[0]
		}

		// Create series
		var seriesID int
		err := pg.db.Conn.QueryRow(`
			INSERT INTO playoff_series (
				playoff_state_id, round, series_order, conference,
				higher_seed_team_id, lower_seed_team_id, higher_seed, lower_seed, best_of, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 7, 'pending')
			RETURNING id
		`, playoffStateID, RoundConferenceFinals, 1, conference, higherSeed.TeamID, lowerSeed.TeamID, higherSeed.Seed, lowerSeed.Seed).Scan(&seriesID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg *NBAPlayoffGenerator) generateNBAFinals(playoffStateID int, conferenceFinalsWinners map[string][]TeamSeed) error {
	// Get Eastern and Western Conference champions
	easternWinner := conferenceFinalsWinners["Eastern"][0]
	westernWinner := conferenceFinalsWinners["Western"][0]

	// TODO: Need to compare regular season wins to determine home field advantage
	var higherSeed, lowerSeed TeamSeed
	if easternWinner.Seed < westernWinner.Seed {
		higherSeed = easternWinner
		lowerSeed = westernWinner
	} else {
		higherSeed = westernWinner
		lowerSeed = easternWinner
	}

	// Create series
	var seriesID int
	err := pg.db.Conn.QueryRow(`
		INSERT INTO playoff_series (
			playoff_state_id, round, series_order, conference,
			higher_seed_team_id, lower_seed_team_id, higher_seed, lower_seed, best_of, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 7, 'pending')
		RETURNING id
	`, playoffStateID, RoundNBAFinals, 1, "NBA Finals", higherSeed.TeamID, lowerSeed.TeamID, higherSeed.Seed, lowerSeed.Seed).Scan(&seriesID)
	if err != nil {
		return err
	}

	return nil
}

func (pg *NBAPlayoffGenerator) getNBASeriesWinners(playoffStateID int, round int) (map[string][]TeamSeed, error) {
	query := `
		SELECT conference, picked_team_id, higher_seed, lower_seed, higher_seed_team_id, lower_seed_team_id
		FROM playoff_series
		WHERE playoff_state_id = $1 AND round = $2
		AND picked_team_id IS NOT NULL
		ORDER BY conference, series_order
	`

	rows, err := pg.db.Query(query, playoffStateID, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	winners := make(map[string][]TeamSeed)
	for rows.Next() {
		var conference *string
		var pickedTeamID, higherSeed, lowerSeed, higherSeedTeamID, lowerSeedTeamID int

		err := rows.Scan(&conference, &pickedTeamID, &higherSeed, &lowerSeed, &higherSeedTeamID, &lowerSeedTeamID)
		if err != nil {
			return nil, err
		}

		conf := "NBA Finals"
		if conference != nil {
			conf = *conference
		}

		var winner TeamSeed
		if pickedTeamID == higherSeedTeamID {
			winner = TeamSeed{TeamID: higherSeedTeamID, Seed: higherSeed}
		} else {
			winner = TeamSeed{TeamID: lowerSeedTeamID, Seed: lowerSeed}
		}

		winners[conf] = append(winners[conf], winner)
	}

	return winners, nil
}

func (pg *NBAPlayoffGenerator) CheckNBASeriesComplete(scenarioID int, round int) (bool, error) {
	playoffStateID, err := pg.getNBAPlayoffStateID(scenarioID)
	if err != nil {
		return false, err
	}

	// Count total series and picked series for the round
	var totalSeries, pickedSeries int
	err = pg.db.Conn.QueryRow(`
		SELECT
			COUNT(*) as total,
			COUNT(picked_team_id) AS picked
		FROM playoff_series
		WHERE playoff_state_id = $1 AND round = $2
	`, playoffStateID, round).Scan(&totalSeries, &pickedSeries)
	if err != nil {
		return false, err
	}

	// Check if all series for a round have been picked
	return totalSeries > 0 && totalSeries == pickedSeries, nil
}

// Checks if all picks for a round are complete
func (pg *NBAPlayoffGenerator) CheckNBARoundComplete(scenarioID int, round int) (bool, error) {
	playoffStateID, err := pg.getNBAPlayoffStateID(scenarioID)
	if err != nil {
		return false, err
	}

	// Check if all play-in matchups are complete
	if round == RoundPlayInA || round == RoundPlayInB {
		var totalMatchups, pickedMatchups int
		err = pg.db.Conn.QueryRow(`
			SELECT
				COUNT(*) as total,
				COUNT(picked_team_id) AS picked
			FROM playoff_matchups
			WHERE playoff_state_id = $1 AND round = $2
		`, playoffStateID, round).Scan(&totalMatchups, &pickedMatchups)

		if err != nil {
			return false, err
		}

		return totalMatchups > 0 && totalMatchups == pickedMatchups, nil
	}

	// For series rounds, delegate to CheckNBASeriesComplete
	return pg.CheckNBASeriesComplete(scenarioID, round)
}

func (pg *NBAPlayoffGenerator) DeleteNBASubsequentRounds(scenarioID int, round int) error {
	var playoffStateID int
	err := pg.db.Conn.QueryRow(`
		SELECT id FROM playoff_states
		WHERE scenario_id = $1
	`, scenarioID).Scan(&playoffStateID)
	if err != nil {
		return err
	}

	// Delete all series from rounds after specified round
	_, err = pg.db.Conn.Exec(`
		DELETE FROM playoff_series
		WHERE playoff_state_id = $1 AND round > $2
	`, playoffStateID, round)
	if err != nil {
		return err
	}

	// Delete all matchups from rounds after specified round
	_, err = pg.db.Conn.Exec(`
		DELETE FROM playoff_matchups
		WHERE playoff_state_id = $1 AND round > $2
	`, playoffStateID, round)
	if err != nil {
		return err
	}

	// Update current round in playoff state
	_, err = pg.db.Conn.Exec(`
		UPDATE playoff_states
		SET current_round = $1, updated_at = NOW()
		WHERE id = $2
	`, round, playoffStateID)

	return err
}

func (pg *NBAPlayoffGenerator) getOrCreateNBAPlayoffState(scenarioID int) (int, error) {
	var id int
	err := pg.db.Conn.QueryRow(`
		SELECT id FROM playoff_states WHERE scenario_id = $1
	`, scenarioID).Scan(&id)

	if err == sql.ErrNoRows {
		err = pg.db.Conn.QueryRow(`
			INSERT INTO playoff_states (scenario_id)
			VALUES ($1)
			RETURNING id
		`, scenarioID).Scan(&id)
	}

	return id, err
}

func (pg *NBAPlayoffGenerator) getNBAPlayoffStateID(scenarioID int) (int, error) {
	var id int
	err := pg.db.Conn.QueryRow(`
		SELECT id FROM playoff_states WHERE scenario_id = $1
	`, scenarioID).Scan(&id)

	return id, err
}