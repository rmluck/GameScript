package handlers

import (
	"time"
	"crypto/rand"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"

	"gamescript/internal/database"
)

func getScenarios(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAuthenticated := false
		if val, ok := c.Locals("is_authenticated").(bool); ok {
			isAuthenticated = val
		}

		userID := 0
		if val, ok := c.Locals("user_id").(int); ok {
			userID = val
		}

		sessionToken := ""
		if val, ok := c.Locals("session_token").(string); ok {
			sessionToken = val
		}

		var scenarios []map[string]interface{}
		var query string
		var args []interface{}

		if isAuthenticated && userID > 0 {
			query = `
				SELECT
					scenario.id, scenario.name, scenario.sport_id, scenario.season_id, scenario.is_public, scenario.created_at, scenario.updated_at,
					sport.short_name as sport_short_name, season.start_year AS season_start_year, season.end_year AS season_end_year
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
		} else if sessionToken != "" {
			query = `
				SELECT
					scenario.id, scenario.name, scenario.sport_id, scenario.season_id, scenario.is_public, scenario.created_at, scenario.updated_at,
					sport.short_name as sport_short_name, season.start_year AS season_start_year, season.end_year AS season_end_year
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
		} else {
			return c.JSON([]map[string]interface{}{})
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		for rows.Next() {
			var scenario map[string]interface{}
			var id, sportID, seasonID, startYear int
			var endYear *int
			var name, sportShortName string
			var isPublic bool
			var createdAt, updatedAt time.Time

			err := rows.Scan(&id, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt, &sportShortName, &startYear, &endYear)
			if err != nil {
				continue
			}

			scenario = map[string]interface{}{
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
			}
			scenarios = append(scenarios, scenario)
		}

		if scenarios == nil {
			scenarios = []map[string]interface{}{}
		}

		return c.JSON(scenarios)
	}
}

func getScenario(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")

		isAuthenticated := false
		if val, ok := c.Locals("is_authenticated").(bool); ok {
			isAuthenticated = val
		}

		userID := 0
		if val, ok := c.Locals("user_id").(int); ok {
			userID = val
		}

		sessionToken := ""
		if val, ok := c.Locals("session_token").(string); ok {
			sessionToken = val
		}

		query := `
			SELECT
				scenario.id, scenario.user_id, scenario.session_token, scenario.name, scenario.sport_id, scenario.season_id, scenario.is_public, scenario.created_at, scenario.updated_at,
				sport.short_name as sport_short_name, season.start_year AS season_start_year, season.end_year AS season_end_year
			FROM
				scenarios scenario
				JOIN sports sport ON scenario.sport_id = sport.id
				JOIN seasons season ON scenario.season_id = season.id
			WHERE
				scenario.id = $1
		`

		var id, sportID, seasonID, startYear int
		var ownerUserID, endYear *int
		var ownerSessionToken *string
		var name, sportShortName string
		var isPublic bool
		var createdAt, updatedAt time.Time

		err := db.Conn.QueryRow(query, scenarioID).Scan(
			&id, &ownerUserID, &ownerSessionToken, &name, &sportID, &seasonID, &isPublic, &createdAt, &updatedAt,
			&sportShortName, &startYear, &endYear,
		)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
		}

		// Verify ownership
		if isAuthenticated && userID > 0 {
			if ownerUserID == nil || *ownerUserID != userID {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		} else if sessionToken != "" {
			currentSessionToken := c.Locals("session_token").(string)
			if ownerSessionToken == nil || *ownerSessionToken != currentSessionToken {
				return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
			}
		} else {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
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

		isAuthenticated := false
		if val, ok := c.Locals("is_authenticated").(bool); ok {
			isAuthenticated = val
		}

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
			// Generate session token for guest
			sessionToken := c.Cookies("session_token")
			if sessionToken == "" {
				sessionToken = generateSessionToken()
				c.Cookie(&fiber.Cookie{
					Name:     "session_token",
					Value:    sessionToken,
					MaxAge:  7 * 24 * 60 * 60, // 7 days
					HTTPOnly: true,
					SameSite: "Lax",
				})
			}
			c.Locals("session_token", sessionToken)

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

		var sportShortName string
		var startYear int
		var endYear *int

		seasonQuery := `
			SELECT sport.short_name, season.start_year, season.end_year
			FROM sports sport
			JOIN seasons season ON sport.id = season.sport_id
			WHERE sport.id = $1 AND season.id = $2
		`

		err := db.Conn.QueryRow(seasonQuery, sportID, seasonID).Scan(&sportShortName, &startYear, &endYear)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve sport information"})
		}

		return c.Status(201).JSON(map[string]interface{}{
			"id": id,
			"name": name,
			"sport_id": sportID,
			"season_id": seasonID,
			"sport_short_name": sportShortName,
			"season_start_year": startYear,
			"season_end_year": endYear,
			"is_public": isPublic,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}
}

func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
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