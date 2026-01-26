// Fetches NBA schedule data from ESPN API and processes it into internal game models

package espn

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	// "os"

	"gamescript/internal/models"
)


const nbaScheduleURL = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"

func (c *Client) FetchNBASchedule(startDate string, endDate string) ([]models.Game, error) {
	// Format: YYYYMMDD-YYYYMMDD
	url := fmt.Sprintf("%s?dates=%s-%s", nbaScheduleURL, startDate, endDate)
	body, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	var scheduleResp models.ESPNScheduleAPIResponse
	if err := json.Unmarshal(body, &scheduleResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	// // Throw raw response into file for further inspection
	// err = os.WriteFile("nba_full_schedule_response.json", body, 0644)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to write response to file: %w", err)
	// }

	// Load Pacific timezone
	pst, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	var games []models.Game
	for _, event := range scheduleResp.Events {
		// Ensure competition exists and has two competitors
		if len(event.Competitions) == 0 {
			continue
		}
		competition := event.Competitions[0]
		if len(competition.Competitors) < 2 {
			continue
		}

		// Remove NBA Cup Final
		if competition.Type.ID == "39" {
			continue
		}

		// Remove All-Star festivities
		if competition.Competitors[0].Team.ID == "-1" || competition.Competitors[1].Team.ID == "-2" {
			continue
		}

		// Parse gametime
		gameTimeUTC, err := time.Parse("2006-01-02T15:04Z", competition.Date)
		if err != nil {
			continue
		}
		gameTimePST := gameTimeUTC.In(pst)
		dayOfWeek := gameTimePST.Weekday().String()

		// Parse location
		location := competition.Venue.FullName
		if competition.Venue.Address.State == "" {
			location += ", " + competition.Venue.Address.City + ", USA"
		} else {
			location += ", " + competition.Venue.Address.City + ", " + competition.Venue.Address.State + ", USA"
		}

		// Find home and away teams
		var homeTeamID, awayTeamID string
		var homeScore, awayScore *int
		for _, competitor := range competition.Competitors {
			if competitor.HomeAway == "home" {
				homeTeamID = competitor.Team.ID
				if score, err := strconv.Atoi(competitor.Score); err == nil && competitor.Score != "" {
					homeScore = &score
				}
			} else {
				awayTeamID = competitor.Team.ID
				if score, err := strconv.Atoi(competitor.Score); err == nil && competitor.Score != "" {
					awayScore = &score
				}
			}
		}

		// Parse primetime info
		primetime := determineNBAPrimetime(gameTimePST)

		// Parse broadcasts
		var network string
		if len(competition.Broadcasts) > 0 {
			network = parseNBABroadcast(competition.Broadcasts[0].Names)
		}

		// Determine game status
		status := "upcoming"
		if competition.Status.Type.Name == "STATUS_FINAL" {
			status = "final"
		} else if competition.Status.Type.Name == "STATUS_POSTPONED" {
			continue
		}

		game := models.Game{
			SeasonID: 			2,
			ESPNID: 			competition.ID,
			StartTime: 			gameTimeUTC,
			DayOfWeek: 			&dayOfWeek,
			Location: 			&location,
			HomeScore: 			homeScore,
			AwayScore: 			awayScore,
			Primetime: 			&primetime,
			Status: 			&status,
			Network: 			&network,
			HomeTeamESPNID: 	&homeTeamID,
			AwayTeamESPNID: 	&awayTeamID,
		}

		games = append(games, game)
	}

	return games, nil
}

func (c *Client) FetchEntireNBASeason(startYear int) ([]models.Game, error) {
	// Update for current NBA season dates
	seasonStart := time.Date(startYear, time.October, 21, 0, 0, 0, 0, time.UTC)
	seasonEnd := time.Date(startYear + 1, time.April, 13, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Fetching NBA season from %s to %s...\n", seasonStart.Format("2006-01-02"), seasonEnd.Format("2006-01-02"))

	var allGames []models.Game

	// Fetch in weekly increments
	currentDate := seasonStart
	weekNum := 1
	for currentDate.Before(seasonEnd) {
		// Calculate week end date
		weekEnd := currentDate.AddDate(0, 0, 7)
		if weekEnd.After(seasonEnd) {
			weekEnd = seasonEnd
		}

		startDateStr := currentDate.Format("20060102")
		endDateStr := weekEnd.Format("20060102")

		fmt.Printf("Fetching week %d: %s to %s...\n", weekNum, startDateStr, endDateStr)

		// Fetch schedule for the week
		games, err := c.FetchNBASchedule(startDateStr, endDateStr)
		if err != nil {
			fmt.Printf("Error fetching games for week %d: %v\n", weekNum, err)
		} else {
			fmt.Printf("  Found %d games\n", len(games))
			allGames = append(allGames, games...)
		}

		// Move to next week
		currentDate = weekEnd.AddDate(0, 0, 1)
		weekNum++

		// Be polite to the API
		time.Sleep(500 * time.Millisecond)
	}

	// Sort games by start time
	for i := 0; i < len(allGames); i++ {
		for j := i + 1; j < len(allGames); j++ {
			if allGames[i].StartTime.After(allGames[j].StartTime) {
				allGames[i], allGames[j] = allGames[j], allGames[i]
			}
		}
	}

	// Assign weeks based on Monday-Sunday groupings
	allGames = assignNBAWeeks(allGames)

	fmt.Printf("Fetched %d NBA games total\n", len(allGames))

	return allGames, nil
}

func assignNBAWeeks(games []models.Game) []models.Game {
	if len(games) == 0 {
		return games
	}

	// Load Pacific timezone
	pst, _ := time.LoadLocation("America/Los_Angeles")

	// Find the first Monday of the season
	firstGameTime := games[0].StartTime.In(pst)
	daysUntilMonday := (int(time.Monday) - int(firstGameTime.Weekday()) + 7) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 0
	} else {
		daysUntilMonday = -(7 - daysUntilMonday)
	}
	seasonStartMonday := firstGameTime.AddDate(0, 0, daysUntilMonday)
	seasonStartMonday = time.Date(seasonStartMonday.Year(), seasonStartMonday.Month(), seasonStartMonday.Day(), 0, 0, 0, 0, pst)

	// Assign games to weeks
	for i := range games {
		gameTime := games[i].StartTime.In(pst)

		// Calculate which week this game falls into
		daysSinceSeasonStart := int(gameTime.Sub(seasonStartMonday).Hours() / 24)
		week := (daysSinceSeasonStart / 7) + 1

		if week < 1 {
			week = 1
		}

		games[i].Week = &week
	}

	return games
}

func determineNBAPrimetime(gameTime time.Time) string {
	var labels []string

	// Get Pacific time components
	month := gameTime.Month()
	day := gameTime.Day()

	// Christmas game
	if month == time.December && day == 25 {
		labels = append(labels, "Christmas")
	}

	if len(labels) == 0 {
		return ""
	}

	// Join with commas
	result := ""
	for i, label := range labels {
		if i > 0 {
			result += ", "
		}
		result += label
	}

	return result
}

func parseNBABroadcast(broadcast []string) string {
	contains := func(slice []string, item string) bool {
		for _, s := range slice {
			if s == item {
				return true
			}
		}
		return false
	}

	network := "League Pass"
	priorities := []string{"ESPN", "NBC", "Prime Video", "ABC", "Peacock", "NBA TV"}

	for _, p := range priorities {
		if contains(broadcast, p) {
			network = p
			break
		}
	}

	return network
}