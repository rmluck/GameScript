package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"gamescript/internal/database"
)

func getPicksByScenario(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		query := `
			SELECT
                pick.id, pick.scenario_id, pick.game_id, pick.picked_team_id, 
                pick.predicted_home_score, pick.predicted_away_score, 
                pick.status, pick.created_at, pick.updated_at,
                game.espn_id, game.start_time, game.week, 
                game.home_team_id, game.away_team_id, 
                game.home_score, game.away_score, game.status as game_status,
                home_team.abbreviation, home_team.city, home_team.name, 
                home_team.conference, home_team.division, 
                home_team.primary_color, home_team.secondary_color, home_team.logo_url,
                away_team.abbreviation, away_team.city, away_team.name, 
                away_team.conference, away_team.division, 
                away_team.primary_color, away_team.secondary_color, away_team.logo_url
            FROM picks pick
            JOIN games game ON pick.game_id = game.id
            JOIN teams home_team ON game.home_team_id = home_team.id
            JOIN teams away_team ON game.away_team_id = away_team.id
            WHERE pick.scenario_id = $1
            ORDER BY game.start_time
		`

		rows, err := db.Query(query, scenarioID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var picks []map[string]interface{}
		for rows.Next() {
			var pickID, scenarioID, gameID, homeTeamID, awayTeamID int
			var pickedTeamID *int
			var week *int
			var predictedHomeScore, predictedAwayScore, homeScore, awayScore *int
			var gameESPNID, homeAbbr, homeCity, homeName, homeConference, homeDivision, homePrimaryColor, homeSecondaryColor, homeLogoURL, awayAbbr, awayCity, awayName, awayConference, awayDivision, awayPrimaryColor, awaySecondaryColor, awayLogoURL string
			var pickedAbbr, pickedCity, pickedName *string
			var pickStatus, gameStatus *string
			var startTime, createdAt, updatedAt time.Time

			err := rows.Scan(
				&pickID, &scenarioID, &gameID, &pickedTeamID, &predictedHomeScore, &predictedAwayScore, &pickStatus, &createdAt, &updatedAt, &gameESPNID, &startTime, &week, &homeTeamID, &awayTeamID, &homeScore, &awayScore, &gameStatus,
				&homeAbbr, &homeCity, &homeName, &homeConference, &homeDivision, &homePrimaryColor, &homeSecondaryColor, &homeLogoURL,
				&awayAbbr, &awayCity, &awayName, &awayConference, &awayDivision, &awayPrimaryColor, &awaySecondaryColor, &awayLogoURL,
				&pickedAbbr, &pickedCity, &pickedName,
			)
			if err != nil {
				continue
			}

			pick := map[string]interface{}{
				"id": pickID,
				"scenario_id": scenarioID,
				"game_id": gameID,
				"picked_team_id": pickedTeamID,
				"predicted_home_score": predictedHomeScore,
				"predicted_away_score": predictedAwayScore,
				"status": pickStatus,
				"created_at": createdAt,
				"updated_at": updatedAt,
				"game": map[string]interface{}{
					"espn_id": gameESPNID,
					"start_time": startTime,
					"week": week,
					"home_score": homeScore,
					"away_score": awayScore,
					"status": gameStatus,
					"home_team": map[string]interface{}{
						"id": homeTeamID,
						"abbreviation": homeAbbr,
						"city": homeCity,
						"name": homeName,
						"conference": homeConference,
						"division": homeDivision,
						"primary_color": homePrimaryColor,
						"secondary_color": homeSecondaryColor,
						"logo_url": homeLogoURL,
					},
					"away_team": map[string]interface{}{
						"id": awayTeamID,
						"abbreviation": awayAbbr,
						"city": awayCity,
						"name": awayName,
						"conference": awayConference,
						"division": awayDivision,
						"primary_color": awayPrimaryColor,
						"secondary_color": awaySecondaryColor,
						"logo_url": awayLogoURL,
					},
				},
			}

			if pickedTeamID != nil {
				pick["picked_team"] = map[string]interface{}{
					"id": pickedTeamID,
					"abbreviation": pickedAbbr,
					"city": pickedCity,
					"name": pickedName,
				}
			}

			picks = append(picks, pick)
		}

		return c.JSON(picks)
	}
}

func getPick(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		gameID := c.Params("game_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		query := `
			SELECT
				pick.id, pick.scenario_id, pick.game_id, pick.picked_team_id, pick.predicted_home_score, pick.predicted_away_score, pick.status, pick.created_at, pick.updated_at
			FROM picks pick
			WHERE pick.scenario_id = $1 AND pick.game_id = $2
		`

		var id, sID, gID int
		var pickedTeamID *int
		var predictedHomeScore, predictedAwayScore *int
		var status *string
		var createdAt, updatedAt time.Time

		err := db.Conn.QueryRow(query, scenarioID, gameID).Scan(
			&id, &sID, &gID, &pickedTeamID, &predictedHomeScore, &predictedAwayScore, &status, &createdAt, &updatedAt,
		)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Pick not found"})
		}

		return c.JSON(map[string]interface{}{
			"id": id,
			"scenario_id": sID,
			"game_id": gID,
			"picked_team_id": pickedTeamID,
			"predicted_home_score": predictedHomeScore,
			"predicted_away_score": predictedAwayScore,
			"status": status,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}
}

func createPick(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		gameID := c.Params("game_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		type CreatePickRequest struct {
			PickedTeamID *int `json:"picked_team_id"`
			PredictedHomeScore *int `json:"predicted_home_score"`
			PredictedAwayScore *int `json:"predicted_away_score"`
		}

		var req CreatePickRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		query := `
			INSERT INTO picks (scenario_id, game_id, picked_team_id, predicted_home_score, predicted_away_score, status)
			VALUES ($1, $2, $3, $4, $5, 'pending')
			RETURNING id, scenario_id, game_id, picked_team_id, predicted_home_score, predicted_away_score, status, created_at, updated_at
		`

		var id, sID, gID int
		var pickedTeamID *int
		var predictedHomeScore, predictedAwayScore *int
		var status *string
		var createdAt, updatedAt time.Time

		err := db.Conn.QueryRow(query, scenarioID, gameID, req.PickedTeamID, req.PredictedHomeScore, req.PredictedAwayScore).Scan(
			&id, &sID, &gID, &pickedTeamID, &predictedHomeScore, &predictedAwayScore, &status, &createdAt, &updatedAt,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(map[string]interface{}{
			"id": id,
			"scenario_id": sID,
			"game_id": gID,
			"picked_team_id": pickedTeamID,
			"predicted_home_score": predictedHomeScore,
			"predicted_away_score": predictedAwayScore,
			"status": status,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}
}

func updatePick(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		gameID := c.Params("game_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		type UpdatePickRequest struct {
			PickedTeamID *int `json:"picked_team_id"`
			PredictedHomeScore *int `json:"predicted_home_score"`
			PredictedAwayScore *int `json:"predicted_away_score"`
		}

		var req UpdatePickRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		query := `
			UPDATE picks
			SET picked_team_id = $1, predicted_home_score = $2, predicted_away_score = $3, updated_at = NOW()
			WHERE scenario_id = $4 AND game_id = $5
			RETURNING id, scenario_id, game_id, picked_team_id, predicted_home_score, predicted_away_score, status, updated_at
		`

		var id, sID, gID int
		var pickedTeamID *int
		var predictedHomeScore, predictedAwayScore *int
		var status *string
		var updatedAt time.Time

		err := db.Conn.QueryRow(query, req.PickedTeamID, req.PredictedHomeScore, req.PredictedAwayScore, scenarioID, gameID).Scan(
			&id, &sID, &gID, &pickedTeamID, &predictedHomeScore, &predictedAwayScore, &status, &updatedAt,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(map[string]interface{}{
			"id": id,
			"scenario_id": sID,
			"game_id": gID,
			"picked_team_id": pickedTeamID,
			"predicted_home_score": predictedHomeScore,
			"predicted_away_score": predictedAwayScore,
			"status": status,
			"updated_at": updatedAt,
		})
	}
}

func deletePick(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		gameID := c.Params("game_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		if !verifyScenarioOwnership(db, scenarioID, isAuthenticated, c) {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		query := `
			DELETE FROM picks
			WHERE scenario_id = $1 AND game_id = $2
		`
		_, err := db.Conn.Exec(query, scenarioID, gameID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Pick deleted successfully"})
	}
}

func verifyScenarioOwnership(db *database.DB, scenarioID string, isAuthenticated bool, c *fiber.Ctx) bool {
	var ownerUserID *int
	var ownerSessionToken *string

	query := `
		SELECT user_id, session_token
		FROM scenarios
		WHERE id = $1
	`
	err := db.Conn.QueryRow(query, scenarioID).Scan(&ownerUserID, &ownerSessionToken)
	if err != nil {
		return false
	}

	// Get current user info from locals
	currentUserID, _ := c.Locals("user_id").(int)
	currentSessionToken, _ := c.Locals("session_token").(string)

	// Check authenticated user ownership
	if isAuthenticated && currentUserID > 0 {
		return ownerUserID != nil && *ownerUserID == currentUserID
	}
	
	// Check guest session ownership
	if !isAuthenticated && currentSessionToken != "" {
		return ownerSessionToken != nil && *ownerSessionToken == currentSessionToken
	}

	return false
}