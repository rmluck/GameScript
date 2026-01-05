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

		generator := playoffs.NewPlayoffGenerator(db)

		// Check if all regular season games are complete
		allComplete, err := generator.CheckAndEnablePlayoffs(sID, seasonID)
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

        generator := playoffs.NewPlayoffGenerator(db)

        // Verify all games are complete
        allComplete, err := generator.CheckAndEnablePlayoffs(sID, seasonID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        if !allComplete {
            return c.Status(400).JSON(fiber.Map{"error": "Not all regular season games are complete"})
        }

        // Generate wild card round (no need to pass seasonID/sportID to playoff state creation)
        err = generator.GenerateWildCardRound(sID, seasonID, sportID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        return c.JSON(fiber.Map{"message": "Playoffs enabled successfully"})
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

		// Get playoff state ID
		var playoffStateID int
		err = db.Conn.QueryRow("SELECT id FROM playoff_states WHERE scenario_id = $1", sID).Scan(&playoffStateID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Playoffs not enabled for this scenario"})
		}

		// Get matchups for this round
		query := `
            SELECT 
                m.id, m.round, m.matchup_order, m.conference,
                m.higher_seed_team_id, m.lower_seed_team_id,
                m.higher_seed, m.lower_seed,
                m.picked_team_id, m.predicted_higher_seed_score, m.predicted_lower_seed_score,
                m.status, m.created_at, m.updated_at,
                ht.abbreviation as higher_abbr, ht.city as higher_city, ht.name as higher_name,
                ht.logo_url as higher_logo, ht.primary_color as higher_color, ht.secondary_color as higher_secondary,
                lt.abbreviation as lower_abbr, lt.city as lower_city, lt.name as lower_name,
                lt.logo_url as lower_logo, lt.primary_color as lower_color, lt.secondary_color as lower_secondary
            FROM playoff_matchups m
            JOIN teams ht ON m.higher_seed_team_id = ht.id
            JOIN teams lt ON m.lower_seed_team_id = lt.id
            WHERE m.playoff_state_id = $1 AND m.round = $2
            ORDER BY m.conference, m.matchup_order
        `

		rows, err := db.Query(query, playoffStateID, round)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var matchups []map[string]interface{}
		for rows.Next() {
			var id, round, matchupOrder, higherSeed, lowerSeed, higherTeamID, lowerTeamID int
			var pickedTeamID, predictedHigherScore, predictedLowerScore *int
			var conference, status *string
			var createdAt, updatedAt string
			var higherAbbr, higherCity, higherName, higherColor, higherSecondary string
			var lowerAbbr, lowerCity, lowerName, lowerColor, lowerSecondary string
			var higherLogo, lowerLogo *string

			err := rows.Scan(
				&id, &round, &matchupOrder, &conference,
				&higherTeamID, &lowerTeamID,
				&higherSeed, &lowerSeed,
				&pickedTeamID, &predictedHigherScore, &predictedLowerScore,
				&status, &createdAt, &updatedAt,
				&higherAbbr, &higherCity, &higherName, &higherLogo, &higherColor, &higherSecondary,
				&lowerAbbr, &lowerCity, &lowerName, &lowerLogo, &lowerColor, &lowerSecondary,
			)
			if err != nil {
				continue
			}

			matchups = append(matchups, map[string]interface{}{
				"id":                          id,
				"round":                       round,
				"matchup_order":               matchupOrder,
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
					"primary_color":   higherColor,
					"secondary_color": higherSecondary,
				},
				"lower_seed_team": map[string]interface{}{
					"id":              lowerTeamID,
					"abbreviation":    lowerAbbr,
					"city":            lowerCity,
					"name":            lowerName,
					"logo_url":        lowerLogo,
					"primary_color":   lowerColor,
					"secondary_color": lowerSecondary,
				},
			})
		}

		return c.JSON(matchups)
	}
}

func updatePlayoffPick(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		matchupID := c.Params("matchup_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		type UpdatePlayoffPickRequest struct {
			PickedTeamID             *int `json:"picked_team_id"`
			PredictedHigherSeedScore *int `json:"predicted_higher_seed_score"`
			PredictedLowerSeedScore  *int `json:"predicted_lower_seed_score"`
		}

		var req UpdatePlayoffPickRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// If both scores are provided, validate that picked team id matches the winning team
		if req.PredictedHigherSeedScore != nil && req.PredictedLowerSeedScore != nil && req.PickedTeamID != nil {
			// Get matchup details to determine higher and lower seed team IDs
			var higherSeedTeamID, lowerSeedTeamID int
			err := db.Conn.QueryRow(`
				SELECT higher_seed_team_id, lower_seed_team_id
				FROM playoff_matchups
				WHERE id = $1
			`, matchupID).Scan(&higherSeedTeamID, &lowerSeedTeamID)

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

		mID, _ := strconv.Atoi(matchupID)
		sID, _ := strconv.Atoi(scenarioID)

		// Get current matchup round
		var currentRound int
		err := db.Conn.QueryRow("SELECT round FROM playoff_matchups WHERE id = $1", mID).Scan(&currentRound)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Matchup not found"})
		}

		// Delete subsequent rounds since we're modifying an earlier round
		generator := playoffs.NewPlayoffGenerator(db)
		err = generator.DeleteSubsequentRounds(sID, currentRound)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to reset subsequent rounds"})
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
		isComplete, err := generator.CheckRoundComplete(sID, currentRound)
		if err == nil && isComplete && currentRound < 4 {
			// Get season ID from scenario
			var seasonID int
			db.Conn.QueryRow("SELECT season_id FROM scenarios WHERE id = $1", sID).Scan(&seasonID)

			// Generate next round
			generator.GenerateNextRound(sID, seasonID, currentRound)
		}

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