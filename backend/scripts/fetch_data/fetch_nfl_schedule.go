// Fetches NFL schedule data

package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"

    "gamescript/internal/services/espn"
	"gamescript/internal/models"
)

func main() {
    // Initialize ESPN client
    client := espn.NewClient()

    // Get year and week from command line args
    year := 2025
    week := 0 // 0 means fetch entire season
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

    var games []models.Game
    var filename string

    if week > 0 {
        // Fetch specific week
        fmt.Printf("Fetching NFL schedule for year %d, week %d...\n", year, week)
        weekGames, err := client.FetchNFLSchedule(year, week)
        if err != nil {
            fmt.Printf("Error fetching schedule: %v\n", err)
            os.Exit(1)
        }
        games = weekGames
        filename = fmt.Sprintf("database/nfl/schedules/nfl_schedule_%d_week_%d.json", year, week)
    } else {
        // Fetch entire season
        fmt.Printf("Fetching entire NFL season for year %d...\n", year)
        allGames, err := client.FetchEntireNFLSeason(year)
        if err != nil {
            fmt.Printf("Error fetching season: %v\n", err)
            os.Exit(1)
        }
        games = allGames
        filename = fmt.Sprintf("database/nfl/schedules/nfl_schedule_%d_full.json", year)
    }

    // Write to JSON file
    file, err := os.Create(filename)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    // Encode games to JSON with indentation
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(games); err != nil {
        fmt.Printf("Error encoding JSON: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("NFL schedule data written to %s (%d games)\n", filename, len(games))
}