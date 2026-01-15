package playoffs

import (
	"database/sql"
	"fmt"
	"gamescript/internal/database"
	"gamescript/internal/standings"
)

const (
	RoundWildCard               = 1
	RoundDivisional             = 2
	RoundConferenceChampionship = 3
	RoundSuperBowl              = 4
)

type NFLPlayoffGenerator struct {
	db *database.DB
}

func NewNFLPlayoffGenerator(db *database.DB) *NFLPlayoffGenerator {
	return &NFLPlayoffGenerator{db: db}
}

// Checks if all regular season games are complete/picked
func (pg *NFLPlayoffGenerator) CheckAndEnableNFLPlayoffs(scenarioID int, seasonID int) (bool, error) {
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

	// Count games that are either:
	// 1. Picked by the user (has a pick with picked_team_id)
	// 2. Already completed (status = 'final')
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

	return totalGames == completedOrPickedGames, nil
}

// Generates the wild card matchups based on standings
func (pg *NFLPlayoffGenerator) GenerateNFLWildCardRound(scenarioID int, seasonID int, sportID int) error {
	// Get current standings
	nflStandings, err := standings.CalculateNFLStandings(pg.db, scenarioID, seasonID)
	if err != nil {
		return fmt.Errorf("failed to calculate standings: %w", err)
	}

	// Create or get playoff state
	playoffStateID, err := pg.getOrCreateNFLPlayoffState(scenarioID)
	if err != nil {
		return err
	}

	// Clear existing wild card matchups (in case of regeneration)
	_, err = pg.db.Conn.Exec(`
        DELETE FROM playoff_matchups 
        WHERE playoff_state_id = $1 AND round = $2
    `, playoffStateID, RoundWildCard)
	if err != nil {
		return err
	}

	// Generate AFC Wild Card matchups
	err = pg.generateNFLConferenceWildCard(playoffStateID, "AFC", nflStandings.AFC.PlayoffSeeds)
	if err != nil {
		return err
	}

	// Generate NFC Wild Card matchups
	err = pg.generateNFLConferenceWildCard(playoffStateID, "NFC", nflStandings.NFC.PlayoffSeeds)
	if err != nil {
		return err
	}

	// Update playoff state
	_, err = pg.db.Conn.Exec(`
        UPDATE playoff_states 
        SET is_enabled = true, current_round = $1, updated_at = NOW()
        WHERE id = $2
    `, RoundWildCard, playoffStateID)

	return err
}

func (pg *NFLPlayoffGenerator) generateNFLConferenceWildCard(playoffStateID int, conference string, seeds []standings.NFLPlayoffSeed) error {
	// Wild Card matchups: (2 vs 7), (3 vs 6), (4 vs 5)
	// Seed 1 gets a bye

	matchups := []struct {
		higher int
		lower  int
		order  int
	}{
		{2, 7, 1},
		{3, 6, 2},
		{4, 5, 3},
	}

	for _, m := range matchups {
		if m.higher > len(seeds) || m.lower > len(seeds) {
			continue
		}

		higherSeed := seeds[m.higher-1]
		lowerSeed := seeds[m.lower-1]

		_, err := pg.db.Conn.Exec(`
            INSERT INTO playoff_matchups (
                playoff_state_id, round, matchup_order, conference,
                higher_seed_team_id, lower_seed_team_id, 
                higher_seed, lower_seed, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
        `, playoffStateID, RoundWildCard, m.order, conference,
			higherSeed.Team.TeamID, lowerSeed.Team.TeamID,
			higherSeed.Seed, lowerSeed.Seed)

		if err != nil {
			return err
		}
	}

	return nil
}

// Generates matchups for the next playoff round
func (pg *NFLPlayoffGenerator) GenerateNFLNextRound(scenarioID int, seasonID int, currentRound int) error {
	playoffStateID, err := pg.getNFLPlayoffStateID(scenarioID)
	if err != nil {
		return err
	}

	// Get winners from current round
	winners, err := pg.getNFLRoundWinners(playoffStateID, currentRound)
	if err != nil {
		return err
	}

	nextRound := currentRound + 1

	// Clear existing matchups for next round
	_, err = pg.db.Conn.Exec(`
        DELETE FROM playoff_matchups 
        WHERE playoff_state_id = $1 AND round >= $2
    `, playoffStateID, nextRound)
	if err != nil {
		return err
	}

	var genErr error
	switch nextRound {
	case RoundDivisional:
		genErr = pg.generateNFLDivisionalRound(playoffStateID, scenarioID, seasonID, winners)
	case RoundConferenceChampionship:
		genErr = pg.generateNFLConferenceChampionships(playoffStateID, winners)
	case RoundSuperBowl:
		genErr = pg.generateNFLSuperBowl(playoffStateID, winners)
	}

	if genErr != nil {
		return genErr
	}

	// Update playoff state to reflect that new round is available
	_, err = pg.db.Conn.Exec(`
        UPDATE playoff_states
        SET current_round = $1, updated_at = NOW()
        WHERE id = $2
    `, nextRound, playoffStateID)

	return err
}

func (pg *NFLPlayoffGenerator) generateNFLDivisionalRound(playoffStateID int, scenarioID int, seasonID int, wildCardWinners map[string][]TeamSeed) error {
	// Get the 1 seeds for both conferences from standings
	nflStandings, err := standings.CalculateNFLStandings(pg.db, scenarioID, seasonID)
	if err != nil {
		return fmt.Errorf("failed to calculate standings: %w", err)
	}

	// Map conference to their 1 seed team ID
	oneSeeds := map[string]int{
		"AFC": nflStandings.AFC.PlayoffSeeds[0].Team.TeamID,
		"NFC": nflStandings.NFC.PlayoffSeeds[0].Team.TeamID,
	}

	for conference, winners := range wildCardWinners {
		if len(winners) != 3 {
			return fmt.Errorf("expected 3 wild card winners for %s, got %d", conference, len(winners))
		}

		// Sort winners by seed (already sorted from query)
		// Matchups: 1 seed vs lowest remaining, higher remaining vs lower remaining
		lowestSeed := winners[len(winners)-1] // Highest seed number
		middleSeed := winners[len(winners)-2] // Middle seed number
		highestSeed := winners[0]             // Lowest seed number (best team)

		// Get the 1 seed team ID for this conference
		oneSeedTeamID, exists := oneSeeds[conference]
		if !exists {
			return fmt.Errorf("no 1 seed found for conference %s", conference)
		}

		// 1 seed vs lowest remaining seed
		_, err := pg.db.Conn.Exec(`
            INSERT INTO playoff_matchups (
                playoff_state_id, round, matchup_order, conference,
                higher_seed_team_id, lower_seed_team_id,
                higher_seed, lower_seed, status
            ) VALUES ($1, $2, $3, $4, $5, $6, 1, $7, 'pending')
        `, playoffStateID, RoundDivisional, 1, conference,
			oneSeedTeamID, lowestSeed.TeamID, lowestSeed.Seed)
		if err != nil {
			return err
		}

		// Higher remaining seed vs lower remaining seed
		var higherTeam, lowerTeam TeamSeed
		if highestSeed.Seed < middleSeed.Seed {
			higherTeam = highestSeed
			lowerTeam = middleSeed
		} else {
			higherTeam = middleSeed
			lowerTeam = highestSeed
		}

		_, err = pg.db.Conn.Exec(`
            INSERT INTO playoff_matchups (
                playoff_state_id, round, matchup_order, conference,
                higher_seed_team_id, lower_seed_team_id,
                higher_seed, lower_seed, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
        `, playoffStateID, RoundDivisional, 2, conference,
			higherTeam.TeamID, lowerTeam.TeamID,
			higherTeam.Seed, lowerTeam.Seed)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg *NFLPlayoffGenerator) generateNFLConferenceChampionships(playoffStateID int, divisionalWinners map[string][]TeamSeed) error {
	for conference, winners := range divisionalWinners {
		if len(winners) != 2 {
			return fmt.Errorf("expected 2 divisional winners for %s, got %d", conference, len(winners))
		}

		higherSeed := winners[0]
		lowerSeed := winners[1]

		_, err := pg.db.Conn.Exec(`
            INSERT INTO playoff_matchups (
                playoff_state_id, round, matchup_order, conference,
                higher_seed_team_id, lower_seed_team_id,
                higher_seed, lower_seed, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
        `, playoffStateID, RoundConferenceChampionship, 1, conference,
			higherSeed.TeamID, lowerSeed.TeamID,
			higherSeed.Seed, lowerSeed.Seed)

		if err != nil {
			return err
		}
	}

	return nil
}

func (pg *NFLPlayoffGenerator) generateNFLSuperBowl(playoffStateID int, conferenceWinners map[string][]TeamSeed) error {
	afcWinner := conferenceWinners["AFC"][0]
	nfcWinner := conferenceWinners["NFC"][0]

	// Higher seed is determined by original playoff seed
	var higherSeed, lowerSeed TeamSeed
	if afcWinner.Seed < nfcWinner.Seed {
		higherSeed = afcWinner
		lowerSeed = nfcWinner
	} else {
		higherSeed = nfcWinner
		lowerSeed = afcWinner
	}

	_, err := pg.db.Conn.Exec(`
        INSERT INTO playoff_matchups (
            playoff_state_id, round, matchup_order, conference,
            higher_seed_team_id, lower_seed_team_id,
            higher_seed, lower_seed, status
        ) VALUES ($1, $2, $3, NULL, $4, $5, $6, $7, 'pending')
    `, playoffStateID, RoundSuperBowl, 1,
		higherSeed.TeamID, lowerSeed.TeamID,
		higherSeed.Seed, lowerSeed.Seed)

	return err
}

type TeamSeed struct {
	TeamID int
	Seed   int
}

func (pg *NFLPlayoffGenerator) getNFLRoundWinners(playoffStateID int, round int) (map[string][]TeamSeed, error) {
	query := `
        SELECT conference, picked_team_id, higher_seed, lower_seed,
               higher_seed_team_id, lower_seed_team_id
        FROM playoff_matchups
        WHERE playoff_state_id = $1 AND round = $2 AND picked_team_id IS NOT NULL
        ORDER BY conference, matchup_order
    `

	rows, err := pg.db.Query(query, playoffStateID, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	winners := make(map[string][]TeamSeed)

	for rows.Next() {
		var conference *string
		var pickedTeamID, higherSeed, lowerSeed, higherTeamID, lowerTeamID int

		err := rows.Scan(&conference, &pickedTeamID, &higherSeed, &lowerSeed, &higherTeamID, &lowerTeamID)
		if err != nil {
			return nil, err
		}

		conf := "SUPER BOWL"
		if conference != nil {
			conf = *conference
		}

		// Determine winner's seed
		var winnerSeed int
		if pickedTeamID == higherTeamID {
			winnerSeed = higherSeed
		} else {
			winnerSeed = lowerSeed
		}

		winners[conf] = append(winners[conf], TeamSeed{
			TeamID: pickedTeamID,
			Seed:   winnerSeed,
		})
	}

	// Sort each conference's winners by seed
	for conf := range winners {
		winners[conf] = sortByNFLSeed(winners[conf])
	}

	return winners, nil
}

func sortByNFLSeed(teams []TeamSeed) []TeamSeed {
	// Simple bubble sort
	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			if teams[j].Seed < teams[i].Seed {
				teams[i], teams[j] = teams[j], teams[i]
			}
		}
	}
	return teams
}

func (pg *NFLPlayoffGenerator) getOrCreateNFLPlayoffState(scenarioID int) (int, error) {
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

func (pg *NFLPlayoffGenerator) getNFLPlayoffStateID(scenarioID int) (int, error) {
	var id int
	err := pg.db.Conn.QueryRow(`
        SELECT id FROM playoff_states WHERE scenario_id = $1
    `, scenarioID).Scan(&id)
	return id, err
}

// CheckNFLRoundComplete checks if all picks for a round are complete
func (pg *NFLPlayoffGenerator) CheckNFLRoundComplete(scenarioID int, round int) (bool, error) {
	playoffStateID, err := pg.getNFLPlayoffStateID(scenarioID)
	if err != nil {
		return false, err
	}

	var totalMatchups, pickedMatchups int
	err = pg.db.Conn.QueryRow(`
        SELECT 
            COUNT(*) as total,
            COUNT(picked_team_id) as picked
        FROM playoff_matchups
        WHERE playoff_state_id = $1 AND round = $2
    `, playoffStateID, round).Scan(&totalMatchups, &pickedMatchups)

	if err != nil {
		return false, err
	}

	return totalMatchups > 0 && totalMatchups == pickedMatchups, nil
}

// DeleteSubsequentNFLRounds deletes all playoff matchups after the specified round
func (pg *NFLPlayoffGenerator) DeleteSubsequentNFLRounds(scenarioID int, round int) error {
	// Get playoff state ID
	var playoffStateID int
	err := pg.db.Conn.QueryRow(`
        SELECT id FROM playoff_states WHERE scenario_id = $1
    `, scenarioID).Scan(&playoffStateID)
	if err != nil {
		return err
	}

	// Delete all matchups from rounds greater than the specified round
	_, err = pg.db.Conn.Exec(`
        DELETE FROM playoff_matchups
        WHERE playoff_state_id = $1 AND round > $2
    `, playoffStateID, round)
	if err != nil {
		return err
	}

	// Update playoff state current_round
	_, err = pg.db.Conn.Exec(`
        UPDATE playoff_states
        SET current_round = $1, updated_at = NOW()
        WHERE id = $2
    `, round, playoffStateID)

	return err
}
