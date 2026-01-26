// Route setup and handlers for API endpoints

package handlers

import (
	// "time"

	"github.com/gofiber/fiber/v2"

	"gamescript/internal/database"
	"gamescript/internal/middleware"
	"gamescript/internal/scheduler"
)


func SetupRoutes(app *fiber.App, db *database.DB, scheduler *scheduler.Scheduler) {
	api := app.Group("/api")

	// // Health check
	// api.Get("/health", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"status": "ok",
	// 		"message": "GameScript API is running",
	// 	})
	// })

	// Sports routes
	api.Get("/sports", getSports(db))

	// Seasons routes
	api.Get("/sports/:sport_id/seasons", getSeasons(db))
	api.Get("/seasons/:season_id", getSeason(db))

	// Teams routes
	api.Get("/seasons/:season_id/teams", getTeamsBySeason(db))
	api.Get("/teams/:team_id", getTeam(db))

	// Games routes
	api.Get("/seasons/:season_id/games", getGamesBySeason(db))
	api.Get("/seasons/:season_id/weeks/:week/games", getGamesByWeek(db))
	api.Get("/teams/:team_id/games", getGamesByTeam(db))
	api.Get("/games/:game_id", getGame(db))

	// Scenarios (optional auth - guest or user)
	scenarios := api.Group("/scenarios")
	scenarios.Use(middleware.OptionalAuth)
	scenarios.Get("/", getScenarios(db))
	scenarios.Post("/", createScenario(db))
	scenarios.Get("/:scenario_id", getScenario(db))
	scenarios.Put("/:scenario_id", updateScenario(db))
	scenarios.Delete("/:scenario_id", deleteScenario(db))
	scenarios.Post("/:scenario_id/claim", middleware.AuthMiddleware, claimScenario(db))
	scenarios.Get("/:scenario_id/standings", getStandings(db))

	// Picks (optional auth - guest or user)
	picks := api.Group("/picks")
	picks.Use(middleware.OptionalAuth)
	picks.Get("/scenarios/:scenario_id", getPicksByScenario(db))
	picks.Get("/scenarios/:scenario_id/games/:game_id", getPick(db))
	picks.Post("/scenarios/:scenario_id/games/:game_id", createPick(db))
	picks.Put("/scenarios/:scenario_id/games/:game_id", updatePick(db))
	picks.Delete("/scenarios/:scenario_id/games/:game_id", deletePick(db))

	// Playoffs (optional auth - guest or user)
	playoffs := api.Group("/playoffs")
	playoffs.Use(middleware.OptionalAuth)
	playoffs.Get("/scenarios/:scenario_id/state", getPlayoffState(db))
	playoffs.Post("/scenarios/:scenario_id/enable", enablePlayoffs(db))
	playoffs.Get("/scenarios/:scenario_id/rounds/:round", getPlayoffMatchups(db))
	playoffs.Put("/scenarios/:scenario_id/matchups/:matchup_id", updatePlayoffPick(db))
	playoffs.Delete("/scenarios/:scenario_id/matchups/:matchup_id", deletePlayoffPick(db))
	playoffs.Post("/scenarios/:scenario_id/generate", generateNextPlayoffRound(db))

	// Admin routes
	admin := api.Group("/admin")
	admin.Post("/update-schedule/nfl", triggerNFLUpdate(scheduler))
	admin.Post("/update-schedule/nba", triggerNBAUpdate(scheduler))
	// admin.Post("/update-schedule/cfb", triggerCFBUpdate(scheduler))
}

func getSports(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, name, short_name, created_at FROM sports")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var sports []map[string]interface{}
		for rows.Next() {
			var id int
			var name, shortName string
			var createdAt string

			if err := rows.Scan(&id, &name, &shortName, &createdAt); err != nil {
				continue
			}

			sports = append(sports, map[string]interface{}{
				"id":         id,
				"name":       name,
				"short_name": shortName,
				"created_at": createdAt,
			})
		}

		return c.JSON(sports)
	}
}

func getSeasons(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sportID := c.Params("sport_id")

		query := `
			SELECT id, sport_id, start_year, end_year, is_active, created_at
			FROM seasons
			WHERE sport_id = $1
			ORDER BY start_year DESC
		`

		rows, err := db.Query(query, sportID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var seasons []map[string]interface{}
		for rows.Next() {
			var id, sportID, startYear int
			var endYear *int
			var isActive bool
			var createdAt string

			if err := rows.Scan(&id, &sportID, &startYear, &endYear, &isActive, &createdAt); err != nil {
				continue
			}

			seasons = append(seasons, map[string]interface{}{
				"id":         id,
				"sport_id":   sportID,
				"start_year": startYear,
				"end_year":   endYear,
				"is_active":  isActive,
				"created_at": createdAt,
			})
		}

		return c.JSON(seasons)
	}
}

func getSeason(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		seasonID := c.Params("season_id")

		query := `
			SELECT id, sport_id, start_year, end_year, is_active, created_at
			FROM seasons
			WHERE id = $1
		`

		var id, sportID, startYear int
		var endYear *int
		var isActive bool
		var createdAt string

		err := db.Conn.QueryRow(query, seasonID).Scan(&id, &sportID, &startYear, &endYear, &isActive, &createdAt)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Season not found"})
		}

		return c.JSON(map[string]interface{}{
			"id":         id,
			"sport_id":   sportID,
			"start_year": startYear,
			"end_year":   endYear,
			"is_active":  isActive,
			"created_at": createdAt,
		})
	}
}

func triggerNFLUpdate(scheduler *scheduler.Scheduler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scheduler.UpdateNFLSchedule()
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "NFL schedule update triggered",
		})
	}
}

func triggerNBAUpdate(scheduler *scheduler.Scheduler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scheduler.UpdateNBASchedule()
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "NBA schedule update triggered",
		})
	}
}

// func triggerCFBUpdate(scheduler *scheduler.Scheduler) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		scheduler.UpdateCFBSchedule()
// 		return c.JSON(fiber.Map{
// 			"status": "ok",
// 			"message": "CFB schedule update trigerred"
// 		})
// 	}
// }