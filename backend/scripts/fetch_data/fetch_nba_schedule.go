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
	client := espn.NewClient()

	// Get year from command line args
	year := 2025

	if len(os.Args) > 1 {
		if y, err := strconv.Atoi(os.Args[1]); err == nil {
			year = y
		}
	}

	var games []models.Game
	var filename string

	// Fetch entire season
	fmt.Printf("Fetching entire NBA season for %d-%d...\n", year, year + 1)
	games, err := client.FetchEntireNBASeason(year)
	if err != nil {
		fmt.Printf("Error fetching season: %v\n", err)
		os.Exit(1)
	}

	// Write to JSON file
	filename = fmt.Sprintf("database/nba/schedules/nba_schedule_%d-%d_full.json", year, year + 1)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(games); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	// Print week breakdown
	weekCounts := make(map[int]int)
	for _, game := range games {
		if game.Week != nil {
			weekCounts[*game.Week]++
		}
	}
	fmt.Printf("\nGames per week:\n")
	for week := 1; week <= len(weekCounts); week++ {
		if count, ok := weekCounts[week]; ok {
			fmt.Printf("  Week %d: %d games\n", week, count)
		}
	}

	fmt.Printf("\nNBA schedule data written to %s (%d games)\n", filename, len(games))
}