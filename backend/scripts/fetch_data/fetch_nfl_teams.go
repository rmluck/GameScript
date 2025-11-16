package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gamescript/internal/services/espn"
)

func main() {
	client := espn.NewClient()

	teams, err := client.FetchNFLTeams()
	if err != nil {
		fmt.Printf("Error fetching teams: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create("database/nfl/teams/nfl_teams_2025.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(teams); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("NFL teams data written to nfl_teams.json")
}