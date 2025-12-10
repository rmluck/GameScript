package handlers

import (
	"github.com/gofiber/fiber/v2"

	"gamescript/internal/database"
)

func getTeamsBySeason(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		seasonID := c.Params("season_id")

		query := `
			SELECT
				team.id, team.sport_id, team.season_id, team.espn_id,
				team.abbreviation, team.city, team.name, 
				team.conference, team.division,
				team.primary_color, team.secondary_color, team.logo_url, team.alternate_logo_url
			FROM teams team
			WHERE team.season_id = $1
			ORDER BY team.name
		`

		rows, err := db.Query(query, seasonID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var teams []map[string]interface {}
		for rows.Next() {
			var id, sportID, seasonID int
			var espnID, abbreviation, city, name, primaryColor, secondaryColor string
			var conference, division, logoURL, alternateLogoURL *string

			err := rows.Scan(
				&id, &sportID, &seasonID, &espnID,
				&abbreviation, &city, &name,
				&conference, &division,
				&primaryColor, &secondaryColor, &logoURL, &alternateLogoURL,
			)
			if err != nil {
				continue
			}

			teams = append(teams, map[string]interface{}{
				"id": id,
				"sport_id": sportID,
				"season_id": seasonID,
				"espn_id": espnID,
				"abbreviation": abbreviation,
				"city": city,
				"name": name,
				"conference": conference,
				"division": division,
				"primary_color": primaryColor,
				"secondary_color": secondaryColor,
				"logo_url": logoURL,
				"alternate_logo_url": alternateLogoURL,
			})
		}

		return c.JSON(teams)
	}
}

func getTeam(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		teamID := c.Params("team_id")

		query := `
			SELECT
				team.id, team.sport_id, team.season_id, team.espn_id,
				team.abbreviation, team.city, team.name, 
				team.conference, team.division,
				team.primary_color, team.secondary_color, team.logo_url, team.alternate_logo_url
			FROM teams team
			WHERE team.id = $1
		`

		var id, sportID, seasonID int
		var espnID, abbreviation, city, name, primaryColor, secondaryColor string
		var conference, division, logoURL, alternateLogoURL *string

		err := db.Conn.QueryRow(query, teamID).Scan(
			&id, &sportID, &seasonID, &espnID,
			&abbreviation, &city, &name,
			&conference, &division,
			&primaryColor, &secondaryColor, &logoURL, &alternateLogoURL,
		)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Team not found"})
		}

		return c.JSON(map[string]interface{}{
			"id": id,
			"sport_id": sportID,
			"season_id": seasonID,
			"espn_id": espnID,
			"abbreviation": abbreviation,
			"city": city,
			"name": name,
			"conference": conference,
			"division": division,
			"primary_color": primaryColor,
			"secondary_color": secondaryColor,
			"logo_url": logoURL,
			"alternate_logo_url": alternateLogoURL,
		})
	}
}