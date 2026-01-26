package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"

    "github.com/joho/godotenv"

    "gamescript/internal/database"
    "gamescript/internal/models"
)

func main() {
    // Get year and week from command line args (optional)
    year := 2025
    week := 0 // 0 means import entire season
    if len(os.Args) > 1 {
        if y, err := strconv.Atoi(os.Args[1]); err == nil {
            year = y
        }
    }
    if len(os.Args) > 2 {
        if w, err := strconv.Atoi(os.Args[2]); err == nil {
            week = w
        }
    }

    // Determine filename
    var filename string
    if week > 0 {
        filename = fmt.Sprintf("database/nfl/schedules/nfl_schedule_%d_week_%d.json", year, week)
    } else {
        filename = fmt.Sprintf("database/nfl/schedules/nfl_schedule_%d_full.json", year)
    }

    // Open JSON file
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file %s: %v", filename, err)
        os.Exit(1)
    }
    defer file.Close()

    // Decode JSON data
    var games []models.Game
    if err := json.NewDecoder(file).Decode(&games); err != nil {
        fmt.Printf("Error decoding JSON: %v", err)
        os.Exit(1)
    }

    // Load environment variables
    if err := godotenv.Load(); err != nil {
        fmt.Printf("Error loading .env file: %v", err)
        os.Exit(1)
    }

    // Connect to database
    db, err := database.NewConnection()
    if err != nil {
        fmt.Printf("Error connecting to database: %v", err)
        os.Exit(1)
    }
    defer db.Close()

    // Insert games into database
    if err := insertNFLGames(db, games); err != nil {
        fmt.Printf("Error inserting games: %v", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully imported %d games from %s\n", len(games), filename)
}

func insertNFLGames(db *database.DB, games []models.Game) error {
    stmt := `
        INSERT INTO games (
            season_id, espn_id, home_team_id, away_team_id, start_time, day_of_week,
            week, location, primetime, network, home_score, away_score, status
        ) VALUES (
            $1, $2,
            (SELECT id FROM teams WHERE season_id = $1 AND espn_id = $3),
            (SELECT id FROM teams WHERE season_id = $1 AND espn_id = $4),
            $5, $6, $7, $8, $9, $10, $11, $12, $13
        )
        ON CONFLICT (season_id, espn_id) DO UPDATE SET
            home_score = EXCLUDED.home_score,
            away_score = EXCLUDED.away_score,
            status = EXCLUDED.status,
            primetime = EXCLUDED.primetime,
            network = EXCLUDED.network
    `

    for _, game := range games {
        // Only set scores if status is "final"
        var homeScore, awayScore *int
        if game.Status != nil && *game.Status == "final" {
            homeScore = game.HomeScore
            awayScore = game.AwayScore
        }

        _, err := db.Conn.Exec(
            stmt,
            game.SeasonID,
            game.ESPNID,
            *game.HomeTeamESPNID,
            *game.AwayTeamESPNID,
            game.StartTime,
            game.DayOfWeek,
            game.Week,
            game.Location,
            game.Primetime,
            game.Network,
            homeScore,  // Will be NULL for upcoming games
            awayScore,  // Will be NULL for upcoming games
            game.Status,
        )
        if err != nil {
            return fmt.Errorf("error inserting game: %v", err)
        }
    }

    return nil
}