package handlers

import (
    "gamescript/internal/database"

    "github.com/gofiber/fiber/v2"
)

func getGamesBySeason(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        seasonID := c.Params("season_id")

        query := `
            SELECT
                game.id, game.season_id, game.espn_id,
                game.home_team_id, game.away_team_id,
                game.start_time, game.day_of_week, game.week,
                game.location, game.primetime, game.network,
                game.home_score, game.away_score, game.status, 
                ht.id as home_id, ht.abbreviation as home_abbr, ht.city as home_city, ht.name as home_name, ht.conference as home_conference, ht.division as home_division, ht.primary_color as home_primary_color, ht.secondary_color as home_secondary_color, ht.logo_url as home_logo_url, ht.alternate_logo_url as home_alternate_logo_url,
                at.id as away_id, at.abbreviation as away_abbr, at.city as away_city, at.name as away_name, at.conference as away_conference, at.division as away_division, at.primary_color as away_primary_color, at.secondary_color as away_secondary_color, at.logo_url as away_logo_url, at.alternate_logo_url as away_alternate_logo_url
            FROM games game
            JOIN teams ht ON game.home_team_id = ht.id
            JOIN teams at ON game.away_team_id = at.id
            WHERE game.season_id = $1
            ORDER BY game.start_time
        `

        rows, err := db.Query(query, seasonID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        defer rows.Close()

        var games[]map[string]interface {}
        for rows.Next() {
            var id, seasonID, homeTeamID, awayTeamID int
            var week, homeScore, awayScore *int
            var espnID, startTime, homeAbbr, homeCity, homeName, homePrimaryColor, homeSecondaryColor, awayAbbr, awayCity, awayName, awayPrimaryColor, awaySecondaryColor string
            var dayOfWeek, location, primetime, network, status, homeConference, homeDivision, homeLogoURL, homeAlternateLogoURL, awayConference, awayDivision, awayLogoURL, awayAlternateLogoURL *string

            err := rows.Scan(
                &id, &seasonID, &espnID,
                &homeTeamID, &awayTeamID,  // ← ADD THESE to Scan
                &startTime, &dayOfWeek, &week,
                &location, &primetime, &network,
                &homeScore, &awayScore, &status, 
                &homeTeamID, &homeAbbr, &homeCity, &homeName, &homeConference, &homeDivision, &homePrimaryColor, &homeSecondaryColor, &homeLogoURL, &homeAlternateLogoURL,
                &awayTeamID, &awayAbbr, &awayCity, &awayName, &awayConference, &awayDivision, &awayPrimaryColor, &awaySecondaryColor, &awayLogoURL, &awayAlternateLogoURL,
            )
            if err != nil {
                continue
            }

            games = append(games, map[string]interface{}{
                "id": id,
                "season_id": seasonID,
                "espn_id": espnID,
                "home_team_id": homeTeamID,
                "away_team_id": awayTeamID,
                "start_time": startTime,
                "day_of_week": dayOfWeek,
                "week": week,
                "location": location,
                "primetime": primetime,
                "network": network,
                "home_score": homeScore,
                "away_score": awayScore,
                "status": status,
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
                    "alternate_logo_url": homeAlternateLogoURL,
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
                    "alternate_logo_url": awayAlternateLogoURL,
                },
            })
        }

        return c.JSON(games)
    }
}

func getGamesByWeek(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        seasonID := c.Params("season_id")
        week := c.Params("week")

        query := `
            SELECT
                game.id, game.season_id, game.espn_id,
                game.home_team_id, game.away_team_id,
                game.start_time, game.day_of_week, game.week,
                game.location, game.primetime, game.network,
                game.home_score, game.away_score, game.status, 
                ht.id as home_id, ht.abbreviation as home_abbr, ht.city as home_city, ht.name as home_name, ht.conference as home_conference, ht.division as home_division, ht.primary_color as home_primary_color, ht.secondary_color as home_secondary_color, ht.logo_url as home_logo_url, ht.alternate_logo_url as home_alternate_logo_url,
                at.id as away_id, at.abbreviation as away_abbr, at.city as away_city, at.name as away_name, at.conference as away_conference, at.division as away_division, at.primary_color as away_primary_color, at.secondary_color as away_secondary_color, at.logo_url as away_logo_url, at.alternate_logo_url as away_alternate_logo_url
            FROM games game
            JOIN teams ht ON game.home_team_id = ht.id
            JOIN teams at ON game.away_team_id = at.id
            WHERE game.season_id = $1 AND game.week = $2
            ORDER BY game.start_time
        `

        rows, err := db.Query(query, seasonID, week)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        defer rows.Close()

        var games []map[string]interface {}
        for rows.Next() {
            var id, seasonID, homeTeamID, awayTeamID int
            var week, homeScore, awayScore *int
            var espnID, startTime, homeAbbr, homeCity, homeName, homePrimaryColor, homeSecondaryColor, awayAbbr, awayCity, awayName, awayPrimaryColor, awaySecondaryColor string
            var dayOfWeek, location, primetime, network, status, homeConference, homeDivision, homeLogoURL, homeAlternateLogoURL, awayConference, awayDivision, awayLogoURL, awayAlternateLogoURL *string

            err := rows.Scan(
                &id, &seasonID, &espnID,
                &homeTeamID, &awayTeamID,  // ← ADD THESE
                &startTime, &dayOfWeek, &week,
                &location, &primetime, &network,
                &homeScore, &awayScore, &status, 
                &homeTeamID, &homeAbbr, &homeCity, &homeName, &homeConference, &homeDivision, &homePrimaryColor, &homeSecondaryColor, &homeLogoURL, &homeAlternateLogoURL,
                &awayTeamID, &awayAbbr, &awayCity, &awayName, &awayConference, &awayDivision, &awayPrimaryColor, &awaySecondaryColor, &awayLogoURL, &awayAlternateLogoURL,
            )
            if err != nil {
                continue
            }

            games = append(games, map[string]interface{}{
                "id": id,
                "season_id": seasonID,
                "espn_id": espnID,
                "home_team_id": homeTeamID,  // ← ADD THIS
                "away_team_id": awayTeamID,  // ← ADD THIS
                "start_time": startTime,
                "day_of_week": dayOfWeek,
                "week": week,
                "location": location,
                "primetime": primetime,
                "network": network,
                "home_score": homeScore,
                "away_score": awayScore,
                "status": status,
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
                    "alternate_logo_url": homeAlternateLogoURL,
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
                    "alternate_logo_url": awayAlternateLogoURL,
                },
            })
        }

        return c.JSON(games)
    }
}

func getGamesByTeam(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        teamID := c.Params("team_id")

        query := `
            SELECT
                game.id, game.season_id, game.espn_id,
                game.home_team_id, game.away_team_id,
                game.start_time, game.day_of_week, game.week,
                game.location, game.primetime, game.network,
                game.home_score, game.away_score, game.status, 
                ht.id as home_id, ht.abbreviation as home_abbr, ht.city as home_city, ht.name as home_name, ht.conference as home_conference, ht.division as home_division, ht.primary_color as home_primary_color, ht.secondary_color as home_secondary_color, ht.logo_url as home_logo_url, ht.alternate_logo_url as home_alternate_logo_url,
                at.id as away_id, at.abbreviation as away_abbr, at.city as away_city, at.name as away_name, at.conference as away_conference, at.division as away_division, at.primary_color as away_primary_color, at.secondary_color as away_secondary_color, at.logo_url as away_logo_url, at.alternate_logo_url as away_alternate_logo_url
            FROM games game
            JOIN teams ht ON game.home_team_id = ht.id
            JOIN teams at ON game.away_team_id = at.id
            WHERE game.home_team_id = $1 OR game.away_team_id = $1
            ORDER BY game.start_time
        `

        rows, err := db.Query(query, teamID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        defer rows.Close()

        var games []map[string]interface{}
        for rows.Next() {
            var id, seasonID, homeTeamID, awayTeamID int
            var week, homeScore, awayScore *int
            var espnID, startTime, homeAbbr, homeCity, homeName, homePrimaryColor, homeSecondaryColor, awayAbbr, awayCity, awayName, awayPrimaryColor, awaySecondaryColor string
            var dayOfWeek, location, primetime, network, status, homeConference, homeDivision, homeLogoURL, homeAlternateLogoURL, awayConference, awayDivision, awayLogoURL, awayAlternateLogoURL *string

            err := rows.Scan(
                &id, &seasonID, &espnID,
                &homeTeamID, &awayTeamID,  // ← ADD THESE
                &startTime, &dayOfWeek, &week,
                &location, &primetime, &network,
                &homeScore, &awayScore, &status, 
                &homeTeamID, &homeAbbr, &homeCity, &homeName, &homeConference, &homeDivision, &homePrimaryColor, &homeSecondaryColor, &homeLogoURL, &homeAlternateLogoURL,
                &awayTeamID, &awayAbbr, &awayCity, &awayName, &awayConference, &awayDivision, &awayPrimaryColor, &awaySecondaryColor, &awayLogoURL, &awayAlternateLogoURL,
            )
            if err != nil {
                continue
            }

            games = append(games, map[string]interface{}{
                "id": id,
                "season_id": seasonID,
                "espn_id": espnID,
                "home_team_id": homeTeamID,  // ← ADD THIS
                "away_team_id": awayTeamID,  // ← ADD THIS
                "start_time": startTime,
                "day_of_week": dayOfWeek,
                "week": week,
                "location": location,
                "primetime": primetime,
                "network": network,
                "home_score": homeScore,
                "away_score": awayScore,
                "status": status,
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
                    "alternate_logo_url": homeAlternateLogoURL,
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
                    "alternate_logo_url": awayAlternateLogoURL,
                },
            })
        }

        return c.JSON(games)
    }
}

func getGame(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        gameID := c.Params("game_id")

        query := `
            SELECT
                game.id, game.season_id, game.espn_id,
                game.home_team_id, game.away_team_id,
                game.start_time, game.day_of_week, game.week,
                game.location, game.primetime, game.network,
                game.home_score, game.away_score, game.status, 
                ht.id as home_id, ht.abbreviation as home_abbr, ht.city as home_city, ht.name as home_name, ht.conference as home_conference, ht.division as home_division, ht.primary_color as home_primary_color, ht.secondary_color as home_secondary_color, ht.logo_url as home_logo_url, ht.alternate_logo_url as home_alternate_logo_url,
                at.id as away_id, at.abbreviation as away_abbr, at.city as away_city, at.name as away_name, at.conference as away_conference, at.division as away_division, at.primary_color as away_primary_color, at.secondary_color as away_secondary_color, at.logo_url as away_logo_url, at.alternate_logo_url as away_alternate_logo_url
            FROM games game
            JOIN teams ht ON game.home_team_id = ht.id
            JOIN teams at ON game.away_team_id = at.id
            WHERE game.id = $1
        `

        var id, seasonID, homeTeamID, awayTeamID int
        var week, homeScore, awayScore *int
        var espnID, startTime, homeAbbr, homeCity, homeName, homePrimaryColor, homeSecondaryColor, awayAbbr, awayCity, awayName, awayPrimaryColor, awaySecondaryColor string
        var dayOfWeek, location, primetime, network, status, homeConference, homeDivision, homeLogoURL, homeAlternateLogoURL, awayConference, awayDivision, awayLogoURL, awayAlternateLogoURL *string

        err := db.Conn.QueryRow(query, gameID).Scan(
            &id, &seasonID, &espnID,
            &homeTeamID, &awayTeamID,  // ← ADD THESE
            &startTime, &dayOfWeek, &week,
            &location, &primetime, &network,
            &homeScore, &awayScore, &status, 
            &homeTeamID, &homeAbbr, &homeCity, &homeName, &homeConference, &homeDivision, &homePrimaryColor, &homeSecondaryColor, &homeLogoURL, &homeAlternateLogoURL,
            &awayTeamID, &awayAbbr, &awayCity, &awayName, &awayConference, &awayDivision, &awayPrimaryColor, &awaySecondaryColor, &awayLogoURL, &awayAlternateLogoURL,
        )
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Game not found"})
        }

        return c.JSON(map[string]interface{}{
            "id": id,
            "season_id": seasonID,
            "espn_id": espnID,
            "home_team_id": homeTeamID,  // ← ADD THIS
            "away_team_id": awayTeamID,  // ← ADD THIS
            "start_time": startTime,
            "day_of_week": dayOfWeek,
            "week": week,
            "location": location,
            "primetime": primetime,
            "network": network,
            "home_score": homeScore,
            "away_score": awayScore,
            "status": status,
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
                "alternate_logo_url": homeAlternateLogoURL,
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
                "alternate_logo_url": awayAlternateLogoURL,
            },
        })
    }
}