package handlers

import (
	"gamescript/internal/database"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *database.DB) {
	api := app.Group("/api")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "GameScript API is running",
		})
	})

	// Sports routes
	api.Get("/sports", getSports(db))

	// TODO: Add more routes as needed
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
				"id": id,
				"name": name,
				"short_name": shortName,
				"created_at": createdAt,
			})
		}

		return c.JSON(sports)
	}
}