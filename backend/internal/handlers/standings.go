package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"gamescript/internal/database"
	"gamescript/internal/standings"
)

func getStandings(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scenarioID := c.Params("scenario_id")

		// Convert scenarioID to int
		sID, err := strconv.Atoi(scenarioID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid scenario ID"})
		}

		// Get scenario to find season
		var seasonID int
		var sportID int
		query := `SELECT season_id, sport_id FROM scenarios WHERE id = $1`
		err = db.Conn.QueryRow(query, sID).Scan(&seasonID, &sportID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
		}

		// Calculate standings based on sport
		var response map[string]interface{}
		if sportID == 1 {
			nflStandings, err := standings.CalculateNFLStandings(db, sID, seasonID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			response = formatNFLStandings(nflStandings)
		}

		return c.JSON(response)
	}
}

// func calculateStandings(db *database.DB) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		scenarioID := c.Params("scenario_id")

// 		// Convert scenarioID to int
// 		sID, err := strconv.Atoi(scenarioID)
// 		if err != nil {
// 			return c.Status(400).JSON(fiber.Map{"error": "Invalid scenarioID"})
// 		}

// 		// Get scenario to find season
// 		var seasonID int
// 		var sportID int
// 		query := `SELECT season_id, sport_id FROM scenarios WHERE id = $1`
// 		err = db.Conn.QueryRow(query, sID).Scan(&seasonID, &sportID)
// 		if err != nil {
// 			return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
// 		}

// 		// Calculate standings based on sport
// 		var response map[string]interface{}
// 		if sportID == 1 {
// 			nflStandings, err := standings.CalculateNFLStandings(db, sID, seasonID)
// 			if err != nil {
// 				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 			}

// 			response = formatNFLStandings(nflStandings)
// 		}

// 		return c.JSON(response)
// 	}
// }

func formatNFLStandings(standings *standings.NFLStandings) map[string]interface{} {
	return map[string]interface{}{
		"afc": map[string]interface{}{
			"divisions": formatDivisions(standings.AFC.Divisions),
			"playoff_seeds": formatPlayoffSeeds(standings.AFC.PlayoffSeeds),
		},
		"nfc": map[string]interface{}{
			"divisions": formatDivisions(standings.NFC.Divisions),
			"playoff_seeds": formatPlayoffSeeds(standings.NFC.PlayoffSeeds),
		},
		"draft_order": formatDraftOrder(standings.DraftOrder),
	}
}

// NOTE: Need to add games back for both conference and division standings
func formatDivisions(divisions map[string][]standings.TeamRecord) map[string]interface {} {
	result := make(map[string]interface{})

	for divName, teams := range divisions {
		formattedTeams := []map[string]interface{}{}
		for i, team := range teams {
			formattedTeams = append(formattedTeams, map[string]interface{}{
				"rank": i + 1,
				"team_id": team.TeamID,
				"team_name": team.TeamName,
				"team_abbr": team.TeamAbbr,
				"wins": team.Wins,
				"losses": team.Losses,
				"ties": team.Ties,
				"win_pct": team.WinPct,
				"division_record": fmt.Sprintf("%d-%d-%d", team.DivisionWins, team.DivisionLosses, team.DivisionTies),
				"conference_record": fmt.Sprintf("%d-%d-%d", team.ConferenceWins, team.ConferenceLosses, team.ConferenceTies),
				"points_for": team.PointsFor,
				"points_against": team.PointsAgainst,
				"point_diff": team.PointsFor - team.PointsAgainst,
				"division_games_back": team.DivisionGamesBack,
				"conference_games_back": team.ConferenceGamesBack,
				"logo_url": team.LogoURL,
				"team_primary_color": team.TeamPrimaryColor,
				"team_secondary_color": team.TeamSecondaryColor,
			})
		}
		result[divName] = formattedTeams
	}

	return result
}

func formatPlayoffSeeds(seeds []standings.PlayoffSeed) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, seed := range seeds {
		result = append(result, map[string]interface{}{
			"seed": seed.Seed,
			"team_id": seed.Team.TeamID,
			"team_name": seed.Team.TeamName,
			"team_abbr": seed.Team.TeamAbbr,
			"wins": seed.Team.Wins,
			"losses": seed.Team.Losses,
			"ties": seed.Team.Ties,
			"is_division_winner": seed.IsDivisionWinner,
			"logo_url": seed.Team.LogoURL,
			"team_primary_color": seed.Team.TeamPrimaryColor,
			"team_secondary_color": seed.Team.TeamSecondaryColor,
		})
	}

	return result
}

func formatDraftOrder(picks []standings.DraftPick) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, pick := range picks {
		result = append(result, map[string]interface{}{
			"pick": pick.Pick,
			"team_id": pick.Team.TeamID,
			"team_name": pick.Team.TeamName,
			"team_abbr": pick.Team.TeamAbbr,
			"record": fmt.Sprintf("%d-%d-%d", pick.Team.Wins, pick.Team.Losses, pick.Team.Ties),
			"reason": pick.Reason,
			"logo_url": pick.Team.LogoURL,
			"team_primary_color": pick.Team.TeamPrimaryColor,
			"team_secondary_color": pick.Team.TeamSecondaryColor,
		})
	}

	return result
}