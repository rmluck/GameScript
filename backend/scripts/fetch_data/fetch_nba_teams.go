// Fetches NBA teams data

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gamescript/internal/services/espn"
)

func main() {
	// Initialize ESPN client
	client := espn.NewClient()

	// Fetch NBA teams
	teams, err := client.FetchNBATeams()
	if err != nil {
		fmt.Printf("Error fetching teams: %v\n", err)
		os.Exit(1)
	}

	// Write to JSON file
	file, err := os.Create("database/nba/teams/nba_teams.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Encode teams to JSON with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(teams); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("NBA teams data written to nba_teams.json")
}