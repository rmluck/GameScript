package handlers

import (
	"gamescript/internal/database"
	"gamescript/internal/playoffs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getPlayoffState(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Check if playoffs should be enabled
		sID, _ := strconv.Atoi(scenarioID)

		// Get scenario to find season and sport
		var seasonID, sportID int
		err := db.Conn.QueryRow("SELECT season_id, sport_id FROM scenarios WHERE id = $1", sID).Scan(&seasonID, &sportID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
		}

		if sportID == 1 {
			generator := playoffs.NewNFLPlayoffGenerator(db)

			// Check if all regular season games are complete
			allComplete, err := generator.CheckAndEnableNFLPlayoffs(sID, seasonID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			// Try to get existing playoff state
			query := `
				SELECT id, scenario_id, current_round, is_enabled, created_at, updated_at
				FROM playoff_states
				WHERE scenario_id = $1
			`

			var playoffState map[string]interface{}
			var id, currentRound int
			var isEnabled bool
			var createdAt, updatedAt string

			err = db.Conn.QueryRow(query, sID).Scan(&id, &sID, &currentRound, &isEnabled, &createdAt, &updatedAt)

			if err == nil {
				playoffState = map[string]interface{}{
					"id":            id,
					"scenario_id":   sID,
					"current_round": currentRound,
					"is_enabled":    isEnabled,
					"created_at":    createdAt,
					"updated_at":    updatedAt,
				}
			} else {
				// No playoff state exists yet
				playoffState = nil
			}

			return c.JSON(fiber.Map{
				"playoff_state": playoffState,
				"can_enable":    allComplete,
			})
		} else if sportID == 2 {
			generator := playoffs.NewNBAPlayoffGenerator(db)

			// Check if all regular season games are complete
			allComplete, err := generator.CheckAndEnableNBAPlayoffs(sID, seasonID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			// Try to get existing playoff state
			query := `
				SELECT id, scenario_id, current_round, is_enabled, created_at, updated_at
				FROM playoff_states
				WHERE scenario_id = $1
			`

			var playoffState map[string]interface{}
			var id, currentRound int
			var isEnabled bool
			var createdAt, updatedAt string

			err = db.Conn.QueryRow(query, sID).Scan(&id, &sID, &currentRound, &isEnabled, &createdAt, &updatedAt)

			if err == nil {
				playoffState = map[string]interface{}{
					"id":            id,
					"scenario_id":   sID,
					"current_round": currentRound,
					"is_enabled":    isEnabled,
					"created_at":    createdAt,
					"updated_at":    updatedAt,
				}
			} else {
				// No playoff state exists yet
				playoffState = nil
			}

			return c.JSON(fiber.Map{
				"playoff_state": playoffState,
				"can_enable":    allComplete,
			})
		}
		
		return c.Status(400).JSON(fiber.Map{"error": "Playoffs not supported for this sport"})
	}
}

func enablePlayoffs(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		sID, _ := strconv.Atoi(scenarioID)

		// Get scenario details (only need seasonID for generating matchups)
		var seasonID, sportID int
		err := db.Conn.QueryRow("SELECT season_id, sport_id FROM scenarios WHERE id = $1", sID).Scan(&seasonID, &sportID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
		}

		if sportID == 1 {
			generator := playoffs.NewNFLPlayoffGenerator(db)

			// Verify all games are complete
			allComplete, err := generator.CheckAndEnableNFLPlayoffs(sID, seasonID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			if !allComplete {
				return c.Status(400).JSON(fiber.Map{"error": "Not all regular season games are complete"})
			}

			// Generate wild card round (no need to pass seasonID/sportID to playoff state creation)
			err = generator.GenerateNFLWildCardRound(sID, seasonID, sportID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			return c.JSON(fiber.Map{"message": "NFL playoffs enabled successfully"})
		} else if sportID == 2 {
			generator := playoffs.NewNBAPlayoffGenerator(db)

			// Verify all games are complete
			allComplete, err := generator.CheckAndEnableNBAPlayoffs(sID, seasonID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			if !allComplete {
				return c.Status(400).JSON(fiber.Map{"error": "Not all regular season games are complete"})
			}

			// Generate play-in round A
			err = generator.GenerateNBAPlayInRoundA(sID, seasonID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			return c.JSON(fiber.Map{"message": "NBA playoffs enabled successfully"})
		}
		
		return c.Status(400).JSON(fiber.Map{"error": "Playoffs not supported for this sport"})
	}
}

func getPlayoffMatchups(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		roundStr := c.Params("round")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		round, err := strconv.Atoi(roundStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid round"})
		}

		sID, _ := strconv.Atoi(scenarioID)

		// Get playoff state ID and sport ID
		var playoffStateID, sportID int
		err = db.Conn.QueryRow(`
			SELECT ps.id, s.sport_id
			FROM playoff_states ps
			JOIN scenarios s ON ps.scenario_id = s.id
			WHERE ps.scenario_id = $1
		`, sID).Scan(&playoffStateID, &sportID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Playoffs not enabled for this scenario"})
		}

		// For NBA, check if this is a series round
		if sportID == 2 && round >= playoffs.RoundConferenceQuarterfinals {
			return getPlayoffSeries(db, c, playoffStateID, round)
		}

		// Get matchups for this round
		query := `
            SELECT 
                m.id, m.round, m.matchup_order, m.game_number, m.conference,
                m.higher_seed_team_id, m.lower_seed_team_id,
                m.higher_seed, m.lower_seed,
                m.picked_team_id, m.predicted_higher_seed_score, m.predicted_lower_seed_score,
                m.status, m.created_at, m.updated_at,
                ht.abbreviation as higher_abbr, ht.city as higher_city, ht.name as higher_name,
                ht.logo_url as higher_logo, ht.alternate_logo_url as higher_alt_logo, ht.primary_color as higher_color, ht.secondary_color as higher_secondary,
                lt.abbreviation as lower_abbr, lt.city as lower_city, lt.name as lower_name,
                lt.logo_url as lower_logo, lt.alternate_logo_url as lower_alt_logo, lt.primary_color as lower_color, lt.secondary_color as lower_secondary
            FROM playoff_matchups m
            JOIN teams ht ON m.higher_seed_team_id = ht.id
            JOIN teams lt ON m.lower_seed_team_id = lt.id
            WHERE m.playoff_state_id = $1 AND m.round = $2
            ORDER BY m.conference, m.matchup_order, m.game_number
        `

		rows, err := db.Query(query, playoffStateID, round)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var matchups []map[string]interface{}
		for rows.Next() {
			var id, round, matchupOrder, higherSeed, lowerSeed, higherTeamID, lowerTeamID int
			var pickedTeamID, predictedHigherScore, predictedLowerScore, gameNumber *int
			var conference, status *string
			var createdAt, updatedAt string
			var higherAbbr, higherCity, higherName, higherColor, higherSecondary string
			var lowerAbbr, lowerCity, lowerName, lowerColor, lowerSecondary string
			var higherLogo, higherAltLogo, lowerLogo, lowerAltLogo *string

			err := rows.Scan(
				&id, &round, &matchupOrder, &gameNumber, &conference,
				&higherTeamID, &lowerTeamID,
				&higherSeed, &lowerSeed,
				&pickedTeamID, &predictedHigherScore, &predictedLowerScore,
				&status, &createdAt, &updatedAt,
				&higherAbbr, &higherCity, &higherName, &higherLogo, &higherAltLogo, &higherColor, &higherSecondary,
				&lowerAbbr, &lowerCity, &lowerName, &lowerLogo, &lowerAltLogo, &lowerColor, &lowerSecondary,
			)
			if err != nil {
				continue
			}

			matchups = append(matchups, map[string]interface{}{
				"id":                          id,
				"round":                       round,
				"matchup_order":               matchupOrder,
				"game_number":                 gameNumber,
				"conference":                  conference,
				"higher_seed":                 higherSeed,
				"lower_seed":                  lowerSeed,
				"higher_seed_team_id":         higherTeamID,
				"lower_seed_team_id":          lowerTeamID,
				"picked_team_id":              pickedTeamID,
				"predicted_higher_seed_score": predictedHigherScore,
				"predicted_lower_seed_score":  predictedLowerScore,
				"status":                      status,
				"created_at":                  createdAt,
				"updated_at":                  updatedAt,
				"higher_seed_team": map[string]interface{}{
					"id":              higherTeamID,
					"abbreviation":    higherAbbr,
					"city":            higherCity,
					"name":            higherName,
					"logo_url":        higherLogo,
					"alternate_logo_url": higherAltLogo,
					"primary_color":   higherColor,
					"secondary_color": higherSecondary,
				},
				"lower_seed_team": map[string]interface{}{
					"id":              lowerTeamID,
					"abbreviation":    lowerAbbr,
					"city":            lowerCity,
					"name":            lowerName,
					"logo_url":        lowerLogo,
					"alternate_logo_url": lowerAltLogo,
					"primary_color":   lowerColor,
					"secondary_color": lowerSecondary,
				},
			})
		}

		return c.JSON(matchups)
	}
}

func getPlayoffSeries(db *database.DB, c *fiber.Ctx, playoffStateID int, round int) error {
	query := `
		SELECT
			ps.id, ps.round, ps.series_order, ps.conference,
			ps.higher_seed_team_id, ps.lower_seed_team_id,
			ps.higher_seed, ps.lower_seed,
			ps.picked_team_id, ps.predicted_higher_seed_wins, ps.predicted_lower_seed_wins,
			ps.best_of, ps.status, ps.created_at, ps.updated_at,
			ht.abbreviation as higher_abbr, ht.city as higher_city, ht.name as higher_name,
			ht.logo_url as higher_logo, ht.alternate_logo_url as higher_alt_logo, ht.primary_color as higher_color, ht.secondary_color as higher_secondary,
			lt.abbreviation as lower_abbr, lt.city as lower_city, lt.name as lower_name,
			lt.logo_url as lower_logo, lt.alternate_logo_url as lower_alt_logo, lt.primary_color as lower_color, lt.secondary_color as lower_secondary
		FROM playoff_series ps
		JOIN teams ht ON ps.higher_seed_team_id = ht.id
		JOIN teams lt ON ps.lower_seed_team_id = lt.id
		WHERE ps.playoff_state_id = $1 AND ps.round = $2
		ORDER BY ps.conference, ps.series_order
	`

	rows, err := db.Query(query, playoffStateID, round)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var series []map[string]interface{}
	for rows.Next() {
		var id, round, seriesOrder, higherSeed, lowerSeed, higherTeamID, lowerTeamID, bestOf int
		var pickedTeamID, predictedHigherWins, predictedLowerWins *int
		var conference, status *string
		var createdAt, updatedAt string
		var higherAbbr, higherCity, higherName, higherColor, higherSecondary string
		var lowerAbbr, lowerCity, lowerName, lowerColor, lowerSecondary string
		var higherLogo, higherAltLogo, lowerLogo, lowerAltLogo *string

		err := rows.Scan(
			&id, &round, &seriesOrder, &conference,
			&higherTeamID, &lowerTeamID,
			&higherSeed, &lowerSeed,
			&pickedTeamID, &predictedHigherWins, &predictedLowerWins,
			&bestOf, &status, &createdAt, &updatedAt,
			&higherAbbr, &higherCity, &higherName, &higherLogo, &higherAltLogo, &higherColor, &higherSecondary,
			&lowerAbbr, &lowerCity, &lowerName, &lowerLogo, &lowerAltLogo, &lowerColor, &lowerSecondary,
		)
		if err != nil {
			continue
		}

		series = append(series, map[string]interface{}{
			"id":                         id,
			"round":                      round,
			"series_order":               seriesOrder,
			"conference":                 conference,
			"higher_seed":                higherSeed,
			"lower_seed":                 lowerSeed,
			"higher_seed_team_id":        higherTeamID,
			"lower_seed_team_id":         lowerTeamID,
			"picked_team_id":             pickedTeamID,
			"predicted_higher_seed_wins": predictedHigherWins,
			"predicted_lower_seed_wins":  predictedLowerWins,
			"best_of":                    bestOf,
			"status":                     status,
			"created_at":                 createdAt,
			"updated_at":                 updatedAt,
			"higher_seed_team": map[string]interface{}{
				"id":              higherTeamID,
				"abbreviation":    higherAbbr,
				"city":            higherCity,
				"name":            higherName,
				"logo_url":        higherLogo,
				"alternate_logo_url": higherAltLogo,
				"primary_color":   higherColor,
				"secondary_color": higherSecondary,
			},
			"lower_seed_team": map[string]interface{}{
				"id":              lowerTeamID,
				"abbreviation":    lowerAbbr,
				"city":            lowerCity,
				"name":            lowerName,
				"logo_url":        lowerLogo,
				"alternate_logo_url": lowerAltLogo,
				"primary_color":   lowerColor,
				"secondary_color": lowerSecondary,
			},
		})
	}

	return c.JSON(series)
}

type UpdatePlayoffPickRequest struct {
	PickedTeamID             *int `json:"picked_team_id"`
	PredictedHigherSeedScore *int `json:"predicted_higher_seed_score"`
	PredictedLowerSeedScore  *int `json:"predicted_lower_seed_score"`
	PredictedHigherSeedWins  *int `json:"predicted_higher_seed_wins"`
	PredictedLowerSeedWins   *int `json:"predicted_lower_seed_wins"`
}

func updatePlayoffPick(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		itemID := c.Params("matchup_id") // This can be matchup or series ID
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var req UpdatePlayoffPickRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		mID, _ := strconv.Atoi(itemID)
		sID, _ := strconv.Atoi(scenarioID)

		// Get sport ID and current round
		var sportID, currentRound int
		err := db.Conn.QueryRow(`
			SELECT s.sport_id, ps.round
			FROM playoff_series ps
			JOIN playoff_states pst ON ps.playoff_state_id = pst.id
			JOIN scenarios s ON pst.scenario_id = s.id
			WHERE ps.id = $1
		`, mID).Scan(&sportID, &currentRound)
		if err == nil {
			// This is a playoff series
			return updatePlayoffSeriesPick(db, c, sID, mID, currentRound, &req)
		}

		// If not a series, look for matchup
		err = db.Conn.QueryRow(`
			SELECT s.sport_id, pm.round
			FROM playoff_matchups pm
			JOIN playoff_states ps ON pm.playoff_state_id = ps.id
			JOIN scenarios s ON ps.scenario_id = s.id
			WHERE pm.id = $1
		`, mID).Scan(&sportID, &currentRound)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Matchup or series not found"})
		}

		// Single elimination game logic
		// If both scores are provided, validate that picked team id matches the winning team
		if req.PredictedHigherSeedScore != nil && req.PredictedLowerSeedScore != nil && req.PickedTeamID != nil {
			// Get matchup details to determine higher and lower seed team IDs
			var higherSeedTeamID, lowerSeedTeamID int
			err := db.Conn.QueryRow(`
				SELECT higher_seed_team_id, lower_seed_team_id
				FROM playoff_matchups
				WHERE id = $1
			`, itemID).Scan(&higherSeedTeamID, &lowerSeedTeamID)

			if err != nil {
				return c.Status(404).JSON(fiber.Map{"error": "Matchup not found"})
			}

			// Determine winner based on scores
			var expectedTeamID int
			if *req.PredictedHigherSeedScore > *req.PredictedLowerSeedScore {
				expectedTeamID = higherSeedTeamID
			} else {
				expectedTeamID = lowerSeedTeamID
			}

			// Override the picked team ID with calculated winner
			req.PickedTeamID = &expectedTeamID
		}

		// Delete subsequent rounds since we're modifying an earlier round
		if sportID == 1 {
			generator := playoffs.NewNFLPlayoffGenerator(db)
			err = generator.DeleteSubsequentNFLRounds(sID, currentRound)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to reset subsequent rounds"})
			}
		} else if sportID == 2 {
			generator := playoffs.NewNBAPlayoffGenerator(db)
			err = generator.DeleteNBASubsequentRounds(sID, currentRound)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to reset subsequent rounds"})
			}
		}

		// Update the pick
		query := `
            UPDATE playoff_matchups
            SET picked_team_id = $1, 
                predicted_higher_seed_score = $2, 
                predicted_lower_seed_score = $3,
				status = 'pending',
                updated_at = NOW()
            WHERE id = $4
            RETURNING id, picked_team_id, predicted_higher_seed_score, predicted_lower_seed_score
        `

		var id int
		var pickedTeamID, predictedHigherScore, predictedLowerScore *int

		err = db.Conn.QueryRow(query, req.PickedTeamID, req.PredictedHigherSeedScore, req.PredictedLowerSeedScore, mID).Scan(
			&id, &pickedTeamID, &predictedHigherScore, &predictedLowerScore,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Check if current round is complete and generate next round
		// if sportID == 1 {
		// 	generator := playoffs.NewNFLPlayoffGenerator(db)
		// 	isComplete, err := generator.CheckNFLRoundComplete(sID, currentRound)
		// 	if err == nil && isComplete && currentRound < playoffs.RoundSuperBowl {
		// 		var seasonID int
		// 		db.Conn.QueryRow("SELECT season_id FROM scenarios WHERE id = $1", sID).Scan(&seasonID)
		// 		generator.GenerateNFLNextRound(sID, seasonID, currentRound)
		// 	}
		// } else if sportID == 2 {
		// 	generator := playoffs.NewNBAPlayoffGenerator(db)
		// 	isComplete, err := generator.CheckNBARoundComplete(sID, currentRound)
		// 	if err == nil && isComplete {
		// 		if currentRound == playoffs.RoundPlayInA {
		// 			generator.GenerateNBAPlayInRoundB(sID)
		// 		} else if currentRound == playoffs.RoundPlayInB {
		// 			var seasonID int
		// 			db.Conn.QueryRow("SELECT season_id FROM scenarios WHERE id = $1", sID).Scan(&seasonID)
		// 			generator.GenerateNBAConferenceQuarterfinals(sID, seasonID)
		// 		}
		// 	}
		// }

		// Update scenario's updated_at timestamp
		_, updateErr := db.Conn.Exec(`
			UPDATE scenarios
			SET updated_at = NOW()
			WHERE id = $1
		`, sID)
		if updateErr != nil {
			c.Locals("scenario_update_error", updateErr.Error())
		}

		return c.JSON(map[string]interface{}{
			"id":                          id,
			"picked_team_id":              pickedTeamID,
			"predicted_higher_seed_score": predictedHigherScore,
			"predicted_lower_seed_score":  predictedLowerScore,
		})
	}
}

func updatePlayoffSeriesPick(db *database.DB, c *fiber.Ctx, scenarioID int, seriesID int, currentRound int, req *UpdatePlayoffPickRequest) error {
	// Validate series prediction
	if req.PredictedHigherSeedWins != nil && req.PredictedLowerSeedWins != nil && req.PickedTeamID != nil {
		if *req.PredictedHigherSeedWins < 0 || *req.PredictedLowerSeedWins < 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Predicted wins cannot be negative"})
		}
		if *req.PredictedHigherSeedWins > 4 || *req.PredictedLowerSeedWins > 4 {
			return c.Status(400).JSON(fiber.Map{"error": "Predicted wins cannot exceed 4"})
		}
		if *req.PredictedHigherSeedWins == 4 && *req.PredictedLowerSeedWins == 4 {
			return c.Status(400).JSON(fiber.Map{"error": "Both teams cannot have 4 wins"})
		}
		if *req.PredictedHigherSeedWins < 4 && *req.PredictedLowerSeedWins < 4 {
			return c.Status(400).JSON(fiber.Map{"error": "One team must have 4 wins"})
		}

		// Determine winner based on predicted wins
		var higherSeedTeamID, lowerSeedTeamID int
		err := db.Conn.QueryRow(`
			SELECT higher_seed_team_id, lower_seed_team_id
			FROM playoff_series
			WHERE id = $1
		`, seriesID).Scan(&higherSeedTeamID, &lowerSeedTeamID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Series not found"})
		}

		if *req.PredictedHigherSeedWins == 4 {
			req.PickedTeamID = &higherSeedTeamID
		} else {
			req.PickedTeamID = &lowerSeedTeamID
		}
	}

	// Delete subsequent rounds
	generator := playoffs.NewNBAPlayoffGenerator(db)
	err := generator.DeleteNBASubsequentRounds(scenarioID, currentRound)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reset subsequent rounds"})
	}

	// Update the series pick
	query := `
		UPDATE playoff_series
		SET picked_team_id = $1,
			predicted_higher_seed_wins = $2,
			predicted_lower_seed_wins = $3,
			status = 'pending',
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, picked_team_id, predicted_higher_seed_wins, predicted_lower_seed_wins
	`

	var id int
	var pickedTeamID, predictedHigherWins, predictedLowerWins *int

	err = db.Conn.QueryRow(query, req.PickedTeamID, req.PredictedHigherSeedWins, req.PredictedLowerSeedWins, seriesID).Scan(&id, &pickedTeamID, &predictedHigherWins, &predictedLowerWins)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if current round is complete and generate next round
	// isComplete, err := generator.CheckNBASeriesComplete(scenarioID, currentRound)
	// if err == nil && isComplete && currentRound < playoffs.RoundNBAFinals {
	// 	generator.GenerateNBANextRound(scenarioID, currentRound)
	// }

	// Update scenario's updated_at timestamp
	db.Conn.Exec(`UPDATE scenarios SET updated_at = NOW() WHERE id = $1`, scenarioID)

	return c.JSON(map[string]interface{}{
		"id":                         id,
		"picked_team_id":             pickedTeamID,
		"predicted_higher_seed_wins": predictedHigherWins,
		"predicted_lower_seed_wins":  predictedLowerWins,
	})
}

func generateNextPlayoffRound(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        scenarioID := c.Params("scenario_id")
        isAuthenticated := c.Locals("is_authenticated").(bool)

        if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
            return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
        }

        sID, _ := strconv.Atoi(scenarioID)

        // Get current playoff state and sport
        var currentRound, sportID, seasonID int
        err := db.Conn.QueryRow(`
            SELECT ps.current_round, s.sport_id, s.season_id
            FROM playoff_states ps
            JOIN scenarios s ON ps.scenario_id = s.id
            WHERE ps.scenario_id = $1
        `, sID).Scan(&currentRound, &sportID, &seasonID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Playoff state not found"})
        }

        // Generate based on sport and current round
		if sportID == 1 {
			generator := playoffs.NewNFLPlayoffGenerator(db)

			// Check if current round is complete
			isComplete, err := generator.CheckNFLRoundComplete(sID, currentRound)
			if err != nil || !isComplete {
				return c.Status(400).JSON(fiber.Map{"error": "Current round is not complete"})
			}

			// Generate next round based on current round
			if currentRound < playoffs.RoundSuperBowl {
				err = generator.GenerateNFLNextRound(sID, seasonID, currentRound)
			} else {
				return c.Status(400).JSON(fiber.Map{"error": "No more rounds to generate"})
			}

			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			return c.JSON(fiber.Map{"message": "Next round generated successfully"})
		} else if sportID == 2 {
            generator := playoffs.NewNBAPlayoffGenerator(db)
            
            // Check if current round is complete
            isComplete, err := generator.CheckNBARoundComplete(sID, currentRound)
            if err != nil || !isComplete {
                return c.Status(400).JSON(fiber.Map{"error": "Current round is not complete"})
            }

            // Generate next round based on current round
            if currentRound == playoffs.RoundPlayInA {
                err = generator.GenerateNBAPlayInRoundB(sID)
            } else if currentRound == playoffs.RoundPlayInB {
                err = generator.GenerateNBAConferenceQuarterfinals(sID, seasonID)
            } else if currentRound < playoffs.RoundNBAFinals {
                err = generator.GenerateNBANextRound(sID, currentRound)
            } else {
                return c.Status(400).JSON(fiber.Map{"error": "No more rounds to generate"})
            }

            if err != nil {
                return c.Status(500).JSON(fiber.Map{"error": err.Error()})
            }

            return c.JSON(fiber.Map{"message": "Next round generated successfully"})
        }

        return c.Status(400).JSON(fiber.Map{"error": "Unsupported sport"})
    }
}

func deletePlayoffPick(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		matchupID := c.Params("matchup_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		mID, err := strconv.Atoi(matchupID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid matchup ID"})
		}

		// Check if this is a series or single matchup
		var playoffSeriesID *int
		err = db.Conn.QueryRow(`
			SELECT playoff_series_id
			FROM playoff_matchups
			WHERE id = $1
		`, mID).Scan(&playoffSeriesID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Matchup not found"})
		}

		if playoffSeriesID != nil {
			// Delete series pick
			query := `
				UPDATE playoff_series
				SET picked_team_id = NULL,
					predicted_higher_seed_wins = NULL,
					predicted_lower_seed_wins = NULL,
					updated_at = NOW()
				WHERE id = $1
				RETURNING id
			`

			var id int
			err = db.Conn.QueryRow(query, *playoffSeriesID).Scan(&id)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to delete playoff series pick"})
			}
		} else {
			// Set pick fields to NULL instead of deleting the row
			query := `
				UPDATE playoff_matchups
				SET picked_team_id = NULL,
					predicted_higher_seed_score = NULL,
					predicted_lower_seed_score = NULL,
					updated_at = NOW()
				WHERE id = $1
				RETURNING id
			`

			var id int
			err = db.Conn.QueryRow(query, mID).Scan(&id)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to delete playoff pick"})
			}
		}

		// Update scenario's updated_at timestamp
		scenarioIDInt, err := strconv.Atoi(scenarioID)
		if err == nil {
			_, updateErr := db.Conn.Exec(`
                UPDATE scenarios
                SET updated_at = NOW()
                WHERE id = $1
            `, scenarioIDInt)
			if updateErr != nil {
				c.Locals("scenario_update_error", updateErr.Error())
			}
		}

		return c.JSON(fiber.Map{"message": "Playoff pick deleted successfully"})
	}
}
