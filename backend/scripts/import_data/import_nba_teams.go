// Imports NBA teams data from JSON file into the database

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"gamescript/internal/database"
	"gamescript/internal/models"
)

func main() {
	// Open JSON file
	file, err := os.Open("database/nba/teams/nba_teams_2025.json")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode JSON data
	var teams []models.Team
	if err := json.NewDecoder(file).Decode(&teams); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := database.NewConnection()
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Insert teams into database
	if err := insertNBATeams(db, teams); err != nil {
		fmt.Printf("Error inserting teams: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully imported %d teams from %s.\n", len(teams), "database/nba/teams/nba_teams_2025.json")
}

func insertNBATeams(db *database.DB, teams []models.Team) error {
	stmt := `
		INSERT INTO teams (
			sport_id, season_id, espn_id, abbreviation, city, name,
			conference, division, primary_color, secondary_color, logo_url, alternate_logo_url,
			created_at
			) VALUES (
			 $1, $2, $3, $4, $5, $6,
			 $7, $8, $9, $10, $11, $12, $13
		)
		ON CONFLICT (season_id, espn_id)
		DO UPDATE SET
			alternate_logo_url = EXCLUDED.alternate_logo_url,
			logo_url = EXCLUDED.logo_url
	`

	for _, team := range teams {
		_, err := db.Conn.Exec(
			stmt,
			team.SportID,
			team.SeasonID,
			team.ESPNID,
			team.Abbreviation,
			team.City,
			team.Name,
			team.Conference,
			team.Division,
			team.PrimaryColor,
			team.SecondaryColor,
			team.LogoURL,
			team.AlternateLogoURL,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("Error inserting team %s: %v\n", team.Name, err)
		}
	}

	return nil
}