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
		} else if sportID == 2 {
			nbaStandings, err := standings.CalculateNBAStandings(db, sID, seasonID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			response = formatNBAStandings(nbaStandings)
		}

		return c.JSON(response)
	}
}

func formatNFLStandings(standings *standings.NFLStandings) map[string]interface{} {
	return map[string]interface{}{
		"afc": map[string]interface{}{
			"divisions":     formatNFLDivisionsAsSeeds(standings.AFC.Divisions, standings.AFC.PlayoffSeeds),
			"playoff_seeds": formatNFLPlayoffSeeds(standings.AFC.PlayoffSeeds),
		},
		"nfc": map[string]interface{}{
			"divisions":     formatNFLDivisionsAsSeeds(standings.NFC.Divisions, standings.NFC.PlayoffSeeds),
			"playoff_seeds": formatNFLPlayoffSeeds(standings.NFC.PlayoffSeeds),
		},
		"draft_order": formatNFLDraftOrder(standings.DraftOrder),
	}
}

func formatNFLDivisionsAsSeeds(divisions map[string][]standings.NFLTeamRecord, allSeeds []standings.NFLPlayoffSeed) map[string]interface{} {
	result := make(map[string]interface{})

	// Create a map of team_id to seed for quick lookup
	seedMap := make(map[int]standings.NFLPlayoffSeed)
	for _, seed := range allSeeds {
		seedMap[seed.Team.TeamID] = seed
	}

	for divName, teams := range divisions {
		formattedTeams := []map[string]interface{}{}
		for _, team := range teams {
			// Find the corresponding seed
			seed, exists := seedMap[team.TeamID]
			if !exists {
				continue
			}

			formattedTeams = append(formattedTeams, map[string]interface{}{
				"seed":                  seed.Seed,
				"team_id":               team.TeamID,
				"team_name":             team.TeamName,
				"team_city":             team.TeamCity,
				"team_abbr":             team.TeamAbbr,
				"wins":                  team.Wins,
				"losses":                team.Losses,
				"ties":                  team.Ties,
				"win_pct":               team.WinPct,
                "home_wins":             team.HomeWins,
				"home_losses":           team.HomeLosses,
				"home_ties":             team.HomeTies,
				"away_wins":             team.AwayWins,
				"away_losses":           team.AwayLosses,
				"away_ties":             team.AwayTies,
                "division_wins":         team.DivisionWins,
				"division_losses":       team.DivisionLosses,
				"division_ties":         team.DivisionTies,
				"conference_wins":       team.ConferenceWins,
				"conference_losses":     team.ConferenceLosses,
				"conference_ties":       team.ConferenceTies,
				"division_games_back":   team.DivisionGamesBack,
				"conference_games_back": team.ConferenceGamesBack,
				"points_for":            team.PointsFor,
				"points_against":        team.PointsAgainst,
				"point_diff":            team.PointsFor - team.PointsAgainst,
				"strength_of_schedule":  team.StrengthOfSchedule,
				"strength_of_victory":   team.StrengthOfVictory,
				"is_division_winner":    seed.IsDivisionWinner,
				"logo_url":              team.LogoURL,
				"team_primary_color":    team.TeamPrimaryColor,
				"team_secondary_color":  team.TeamSecondaryColor,
			})
		}
		result[divName] = formattedTeams
	}

	return result
}

func formatNFLPlayoffSeeds(seeds []standings.NFLPlayoffSeed) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, seed := range seeds {
		result = append(result, map[string]interface{}{
			"seed":                  seed.Seed,
			"team_id":               seed.Team.TeamID,
			"team_name":             seed.Team.TeamName,
			"team_city":             seed.Team.TeamCity,
			"team_abbr":             seed.Team.TeamAbbr,
			"wins":                  seed.Team.Wins,
			"losses":                seed.Team.Losses,
			"ties":                  seed.Team.Ties,
			"win_pct":               seed.Team.WinPct,
			"home_wins":             seed.Team.HomeWins,
			"home_losses":           seed.Team.HomeLosses,
			"home_ties":             seed.Team.HomeTies,
			"away_wins":             seed.Team.AwayWins,
			"away_losses":           seed.Team.AwayLosses,
			"away_ties":             seed.Team.AwayTies,
			"division_wins":         seed.Team.DivisionWins,
			"division_losses":       seed.Team.DivisionLosses,
			"division_ties":         seed.Team.DivisionTies,
			"conference_wins":       seed.Team.ConferenceWins,
			"conference_losses":     seed.Team.ConferenceLosses,
			"conference_ties":       seed.Team.ConferenceTies,
			"division_games_back":   seed.Team.DivisionGamesBack,
			"conference_games_back": seed.Team.ConferenceGamesBack,
			"points_for":            seed.Team.PointsFor,
			"points_against":        seed.Team.PointsAgainst,
			"point_diff":            seed.Team.PointsFor - seed.Team.PointsAgainst,
			"strength_of_schedule":  seed.Team.StrengthOfSchedule,
			"strength_of_victory":   seed.Team.StrengthOfVictory,
            "is_division_winner":    seed.IsDivisionWinner,
			"logo_url":              seed.Team.LogoURL,
			"team_primary_color":    seed.Team.TeamPrimaryColor,
			"team_secondary_color":  seed.Team.TeamSecondaryColor,
		})
	}

	return result
}

func formatNFLDraftOrder(picks []standings.NFLDraftPick) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, pick := range picks {
		result = append(result, map[string]interface{}{
			"pick":                 pick.Pick,
			"team_id":              pick.Team.TeamID,
			"team_name":            pick.Team.TeamName,
			"team_abbr":            pick.Team.TeamAbbr,
			"record":               fmt.Sprintf("%d-%d-%d", pick.Team.Wins, pick.Team.Losses, pick.Team.Ties),
			"logo_url":             pick.Team.LogoURL,
			"team_primary_color":   pick.Team.TeamPrimaryColor,
			"team_secondary_color": pick.Team.TeamSecondaryColor,
		})
	}

	return result
}

func formatNBAStandings(standings *standings.NBAStandings) map[string]interface{} {
	return map[string]interface{}{
		"eastern": map[string]interface{}{
			"divisions":     formatNBADivisionsAsSeeds(standings.Eastern.Divisions, standings.Eastern.PlayoffSeeds),
			"playoff_seeds": formatNBAPlayoffSeeds(standings.Eastern.PlayoffSeeds),
		},
		"western": map[string]interface{}{
			"divisions":     formatNBADivisionsAsSeeds(standings.Western.Divisions, standings.Western.PlayoffSeeds),
			"playoff_seeds": formatNBAPlayoffSeeds(standings.Western.PlayoffSeeds),
		},
		"draft_order": formatNBADraftOrder(standings.DraftOrder),
	}
}

func formatNBADivisionsAsSeeds(divisions map[string][]standings.NBATeamRecord, allSeeds []standings.NBAPlayoffSeed) map[string]interface{} {
    result := make(map[string]interface{})

    // Create a map of team_id to seed for quick lookup
    seedMap := make(map[int]standings.NBAPlayoffSeed)
    for _, seed := range allSeeds {
        seedMap[seed.Team.TeamID] = seed
    }

    for divName, teams := range divisions {
        formattedTeams := []map[string]interface{}{}
        for _, team := range teams {
            // Find the corresponding seed
            seed, exists := seedMap[team.TeamID]
            if !exists {
                continue
            }

            formattedTeams = append(formattedTeams, map[string]interface{}{
                "seed":                 seed.Seed,
                "team_id":              team.TeamID,
                "team_name":            team.TeamName,
                "team_city":            team.TeamCity,
                "team_abbr":            team.TeamAbbr,
                "wins":                 team.Wins,
                "losses":               team.Losses,
                "win_pct":              team.WinPct,
                "home_wins":            team.HomeWins,
                "home_losses":          team.HomeLosses,
                "away_wins":            team.AwayWins,
                "away_losses":          team.AwayLosses,
                "division_wins":        team.DivisionWins,
                "division_losses":      team.DivisionLosses,
                "conference_wins":      team.ConferenceWins,
                "conference_losses":    team.ConferenceLosses,
                "division_games_back":  team.DivisionGamesBack,
                "conference_games_back": team.ConferenceGamesBack,
                "points_for":           team.PointsFor,
                "points_against":       team.PointsAgainst,
                "games_with_scores":    team.GamesWithScores,
                "strength_of_schedule": team.StrengthOfSchedule,
                "strength_of_victory":  team.StrengthOfVictory,
                "is_division_winner":   seed.IsDivisionWinner,
                "logo_url":             team.LogoURL,
                "team_primary_color":   team.TeamPrimaryColor,
                "team_secondary_color": team.TeamSecondaryColor,
            })
        }
        result[divName] = formattedTeams
    }

    return result
}

func formatNBAPlayoffSeeds(seeds []standings.NBAPlayoffSeed) []map[string]interface{} {
    result := []map[string]interface{}{}

    for _, seed := range seeds {
        result = append(result, map[string]interface{}{
            "seed":                 seed.Seed,
            "team_id":              seed.Team.TeamID,
            "team_name":            seed.Team.TeamName,
            "team_city":            seed.Team.TeamCity,
            "team_abbr":            seed.Team.TeamAbbr,
            "wins":                 seed.Team.Wins,
            "losses":               seed.Team.Losses,
            "win_pct":              seed.Team.WinPct,
            "home_wins":            seed.Team.HomeWins,
            "home_losses":          seed.Team.HomeLosses,
            "away_wins":            seed.Team.AwayWins,
            "away_losses":          seed.Team.AwayLosses,
            "division_wins":        seed.Team.DivisionWins,
            "division_losses":      seed.Team.DivisionLosses,
            "conference_wins":      seed.Team.ConferenceWins,
            "conference_losses":    seed.Team.ConferenceLosses,
            "division_games_back":  seed.Team.DivisionGamesBack,
            "conference_games_back": seed.Team.ConferenceGamesBack,
            "points_for":           seed.Team.PointsFor,
            "points_against":       seed.Team.PointsAgainst,
            "games_with_scores":    seed.Team.GamesWithScores,
            "strength_of_schedule": seed.Team.StrengthOfSchedule,
            "strength_of_victory":  seed.Team.StrengthOfVictory,
            "is_division_winner":   seed.IsDivisionWinner,
            "logo_url":             seed.Team.LogoURL,
            "team_primary_color":   seed.Team.TeamPrimaryColor,
            "team_secondary_color": seed.Team.TeamSecondaryColor,
        })
    }

    return result
}

func formatNBADraftOrder(picks []standings.NBADraftPick) []map[string]interface{} {
    result := []map[string]interface{}{}

    for _, pick := range picks {
        result = append(result, map[string]interface{}{
            "pick":                 pick.Pick,
            "team_id":              pick.Team.TeamID,
            "team_name":            pick.Team.TeamName,
            "team_abbr":            pick.Team.TeamAbbr,
            "record":               fmt.Sprintf("%d-%d", pick.Team.Wins, pick.Team.Losses),
            "logo_url":             pick.Team.LogoURL,
            "team_primary_color":   pick.Team.TeamPrimaryColor,
            "team_secondary_color": pick.Team.TeamSecondaryColor,
        })
    }

    return result
}