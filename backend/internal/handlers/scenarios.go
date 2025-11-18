package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"gamescript/internal/database"
)

func getScenarios(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAuthenticated := c.Locals("is_authenticated").(bool)

		var query string
		var args []interface{}

		if isAuthenticated {
			userID := c.Locals("user_id").(int)
			query = `
				SELECT
					scenario.id, scenario.user_id, scenario.name, scenario.sport_id, scenario.season_id, scenario.is_public, scenario.created_at, scenario.updated_at,
					sport.short_name as sport_short_name, season.start_year, season.end_year
				FROM
					scenarios scenario
					JOIN sports sport ON scenario.sport_id = sport.id
					JOIN seasons season ON scenario.season_id = season.id
				WHERE
					scenario.user_id = $1
				ORDER BY
					scenario.created_at DESC
				`
				args = []interface{}{userID}
		} else {
			sessionToken := c.Locals("session_token").(string)
			query = `
				SELECT
					scenario.id, scenario.session_token, scenario.name, scenario.sport_id, scenario.season_id, scenario.is_public, scenario.created_at, scenario.updated_at,
					sport.short_name as sport_short_name, season.start_year, season.end_year
				FROM
					scenarios scenario
					JOIN sports sport ON scenario.sport_id = sport.id
					JOIN seasons season ON scenario.season_id = season.id
				WHERE
					scenario.session_token = $1
				ORDER BY
					scenario.updated_at DESC
			`
			args = []interface{}{sessionToken}
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var scenarios []map[string]interface{}
		for rows.Next() {
			var id, sportID, seasonID, startYear int
			var userID, endYear *int
			var sessionToken *string
			var name, sportShortName string
			var isPublic bool
			var createdAt, updatedAt time.Time

			if isAuthenticated {
				err = rows.Scan(&id, &userID, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt, &sportShortName, &startYear, &endYear)
			} else {
				err = rows.Scan(&id, &sessionToken, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt, &sportShortName, &startYear, &endYear)
			}

			if err != nil {
				continue
			}

			scenarios = append(scenarios, map[string]interface{}{
				"id": id,
				"name": name,
				"sport_id": sportID,
				"season_id": seasonID,
				"season_start_year": startYear,
				"season_end_year": endYear,
				"is_public": isPublic,
				"created_at": createdAt,
				"updated_at": updatedAt,
				"sport_short_name": sportShortName,
			})
		}

		return c.JSON(scenarios)
	}
}

func getScenario(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		query := `
			SELECT
				scenario.id, scenario.user_id, scenario.session_token, scenario.name, scenario.sport_id, scenario.season_id, scenario.is_public, scenario.created_at, scenario.updated_at,
				sport.short_name as sport_short_name, season.start_year, season.end_year
			FROM
				scenarios scenario
				JOIN sports sport ON scenario.sport_id = sport.id
				JOIN seasons season ON scenario.season_id = season.id
			WHERE
				scenario.id = $1
		`

		var id, sportID, seasonID, startYear int
		var userID, endYear *int
		var sessionToken *string
		var name, sportShortName string
		var isPublic bool
		var createdAt, updatedAt time.Time

		err := db.Conn.QueryRow(query, scenarioID).Scan(
			&id, &userID, &sessionToken, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt,
			&sportShortName, &startYear, &endYear,
		)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
		}

		// Verify ownership
		if isAuthenticated {
			currentUserID := c.Locals("user_id").(int)
			if userID == nil || *userID != currentUserID {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		} else {
			currentSessionToken := c.Locals("session_token").(string)
			if sessionToken == nil || *sessionToken != currentSessionToken {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		}

		return c.JSON(map[string]interface{}{
			"id": id,
			"name": name,
			"sport_id": sportID,
			"season_id": seasonID,
			"season_start_year": startYear,
			"season_end_year": endYear,
			"is_public": isPublic,
			"created_at": createdAt,
			"updated_at": updatedAt,
			"sport_short_name": sportShortName,
		})
	}
}

func createScenario(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type CreateScenarioRequest struct {
			Name     string `json:"name"`
			SportID  int    `json:"sport_id"`
			SeasonID int    `json:"season_id"`
			IsPublic bool   `json:"is_public"`
		}

		var req CreateScenarioRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Validate required fields
		if req.Name == "" || req.SportID == 0 || req.SeasonID == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
		}

		isAuthenticated := c.Locals("is_authenticated").(bool)

		var query string
		var args []interface{}

		if isAuthenticated {
			userID := c.Locals("user_id").(int)
			query = `
				INSERT INTO scenarios (user_id, name, sport_id, season_id, is_public)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, user_id, name, sport_id, season_id, is_public, created_at, updated_at
			`
			args = []interface{}{userID, req.Name, req.SportID, req.SeasonID, req.IsPublic}
		} else {
			sessionToken := c.Locals("session_token").(string)
			query = `
				INSERT INTO scenarios (session_token, name, sport_id, season_id, is_public)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, session_token, name, sport_id, season_id, is_public, created_at, updated_at
			`
			args = []interface{}{sessionToken, req.Name, req.SportID, req.SeasonID, req.IsPublic}
		}

		var id, sportID, seasonID int
		var userID *int
		var sessionToken *string
		var name string
		var isPublic bool
		var createdAt, updatedAt time.Time

		if isAuthenticated{
			err := db.Conn.QueryRow(query, args...).Scan(
				&id, &userID, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt,
			)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			err := db.Conn.QueryRow(query, args...).Scan(
				&id, &sessionToken, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt,
			)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		return c.Status(201).JSON(map[string]interface{}{
			"id": id,
			"name": name,
			"sport_id": sportID,
			"season_id": seasonID,
			"is_public": isPublic,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}
}

func updateScenario(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		type UpdateScenarioRequest struct {
			Name     *string `json:"name"`
			IsPublic *bool   `json:"is_public"`
		}

		var req UpdateScenarioRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Verify ownership
		var ownerUserID *int
		var ownerSessionToken *string
		verifyQuery := `
			SELECT user_id, session_token
			FROM scenarios
			WHERE id = $1
		`
		err := db.Conn.QueryRow(verifyQuery, scenarioID).Scan(&ownerUserID, &ownerSessionToken)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
		}

		if isAuthenticated {
			currentUserID := c.Locals("user_id").(int)
			if ownerUserID == nil || *ownerUserID != currentUserID {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		} else {
			currentSessionToken := c.Locals("session_token").(string)
			if ownerSessionToken == nil || *ownerSessionToken != currentSessionToken {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		}

		// Build update query dynamically
		updateFields := []string{}
		args := []interface{}{}
		argCount := 1

		if req.Name != nil {
			updateFields = append(updateFields, "name = $"+string(rune('0'+argCount)))
			args = append(args, *req.Name)
			argCount++
		}

		if req.IsPublic != nil {
			updateFields = append(updateFields, "is_public = $"+string(rune('0'+argCount)))
			args = append(args, *req.IsPublic)
			argCount++
		}

		if len(updateFields) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "No fields to update"})
		}

		updateFields = append(updateFields, "updated_at = NOW()")
		args = append(args, scenarioID)

		query := `UPDATE scenarios SET ` + updateFields[0]
		for i := 1; i < len(updateFields); i++ {
			query += ", " + updateFields[i]
		}
		query += ` WHERE id = $` + string(rune('0'+argCount)) + ` RETURNING id, name, sport_id, season_id, is_public, created_at, updated_at`

		var id, sportID, seasonID int
		var name string
		var isPublic bool
		var createdAt, updatedAt time.Time

		err = db.Conn.QueryRow(query, args...).Scan(&id, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(map[string]interface{}{
			"id": id,
			"name": name,
			"sport_id": sportID,
			"season_id": seasonID,
			"is_public": isPublic,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}
}

func deleteScenario(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		isAuthenticated := c.Locals("is_authenticated").(bool)

		// Verify ownership
		var ownerUserID *int
		var ownerSessionToken *string
		verifyQuery := `
			SELECT user_id, session_token
			FROM scenarios
			WHERE id = $1
		`
		err := db.Conn.QueryRow(verifyQuery, scenarioID).Scan(&ownerUserID, &ownerSessionToken)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
		}

		if isAuthenticated {
			currentUserID := c.Locals("user_id").(int)
			if ownerUserID == nil || *ownerUserID != currentUserID {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		} else {
			currentSessionToken := c.Locals("session_token").(string)
			if ownerSessionToken == nil || *ownerSessionToken != currentSessionToken {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		}

		deleteQuery := `
			DELETE FROM scenarios
			WHERE id = $1
		`
		_, err = db.Conn.Exec(deleteQuery, scenarioID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Scenario deleted successfully"})
	}
}

func claimScenario(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")
		sessionToken := c.Cookies("session_token")
		userID := c.Locals("user_id").(int)

		if sessionToken == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Session token required to claim scenario"})
		}

		updateQuery := `
			UPDATE scenarios
			SET user_id = $1, session_token = NULL, updated_at = NOW()
			WHERE id = $2 AND session_token = $3
			RETURNING id
		`

		var id int
		err := db.Conn.QueryRow(updateQuery, userID, scenarioID, sessionToken).Scan(&id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found or already claimed"})
		}

		return c.JSON(fiber.Map{
			"message": "Scenario claimed successfully",
			"id":      id,
		})
	}
}