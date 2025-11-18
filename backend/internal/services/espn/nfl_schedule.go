package espn

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"strings"

	// "os"

	"gamescript/internal/models"
)

const nflScheduleURL = "https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard"

func (c *Client) FetchNFLSchedule(year int, week int) ([]models.Game, error) {
	url := fmt.Sprintf("%s?dates=%d&seasontype=2&week=%d", nflScheduleURL, year, week)

	body, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	var scheduleResp models.ESPNScheduleAPIResponse
	if err := json.Unmarshal(body, &scheduleResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	// Throw raw response into file for further inspection
	// err = os.WriteFile("nfl_full_schedule_response.json", body, 0644)
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
		if len(event.Competitions) == 0 {
			continue
		}

		competition := event.Competitions[0]
		if len(competition.Competitors) < 2 {
			continue
		}

		// Parse the game time
		gameTimeUTC, err := time.Parse("2006-01-02T15:04Z", competition.Date)
		if err != nil {
			continue
		}
		gameTime := gameTimeUTC.In(pst)
		dayOfWeek := gameTime.Weekday().String()

		// Parse the location
		location := competition.Venue.FullName
		if competition.Venue.Address.State == "" {
			location +=  ", " + competition.Venue.Address.City + ", " + competition.Venue.Address.Country
		} else {
			location +=  ", " + competition.Venue.Address.City + ", " + competition.Venue.Address.State + ", " + competition.Venue.Address.Country
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
		primetime := determinePrimetime(gameTime, location)

		// Parse broadcasts
		var network string
		if len(competition.Broadcasts) > 0 {
			network = parseBroadcast(competition.Broadcasts[0].Names)
		}

		// Determine game status
		status := "upcoming"
		if competition.Status.Type.Name == "STATUS_FINAL" {
			status = "final"
		}		

		game := models.Game{
			SeasonID: 			1,
			ESPNID: 			competition.ID,
			StartTime: 			gameTime,
			DayOfWeek: 			&dayOfWeek,
			Week: 				&week,
			Location: 			&location,
			HomeScore: 			homeScore,
			AwayScore: 			awayScore,
			Primetime: 			&primetime,
			Status:    			&status,
			Network: 			&network,
			IsPostseason: 		week > 18,
			HomeTeamESPNID: 	&homeTeamID,
			AwayTeamESPNID: 	&awayTeamID,
		}

		games = append(games, game)
	}

	return games, nil
}

func (c *Client) FetchEntireNFLSeason(year int) ([]models.Game, error) {
	var allGames []models.Game

	// Regular season: weeks 1-18
	for week := 1; week <= 18; week++ {
		fmt.Printf("Fetching NFL week %d...\n", week)
		weekGames, err := c.FetchNFLSchedule(year, week)
		if err != nil {
			fmt.Printf("Error fetching week %d: %v\n", week, err)
			continue
		}
		allGames = append(allGames, weekGames...)

		// Be polite to the API
		time.Sleep(500 * time.Millisecond)
	}

	// TODO: Add postseason weeks if needed

	return allGames, nil
}

func determinePrimetime(gameTime time.Time, location string) string {
	var labels []string

	// Get Pacific time components
	weekday := gameTime.Weekday()
	hour := gameTime.Hour()
	month := gameTime.Month()
	day := gameTime.Day()

	// Thursday Night Football
	if weekday == time.Thursday {
		labels = append(labels, "TNF")
	}

	// Monday Night Football
	if weekday == time.Monday {
		labels = append(labels, "MNF")
	}

	// Sunday Night Football
	if weekday == time.Sunday && hour >= 17 {
		labels = append(labels, "SNF")
	}

	// Friday game
	if weekday == time.Friday {
		labels = append(labels, "Friday")
	}

	// Saturday game
	if weekday == time.Saturday {
		labels = append(labels, "Saturday")
	}

	// Christmas game
	if month == time.December && day == 25 {
		labels = append(labels, "Christmas")
	}

	// Thanksgiving game
	if month == time.November && weekday == time.Thursday {
        // Find the first Thursday of the month
		firstDayOfMonth := time.Date(gameTime.Year(), time.November, 1, 0, 0, 0, 0, gameTime.Location())
		daysUntilThursday := (int(time.Thursday) - int(firstDayOfMonth.Weekday()) + 7) % 7
		firstThursday := 1 + daysUntilThursday
		
		// 4th Thursday is 3 weeks (21 days) after the first Thursday
		fourthThursday := firstThursday + 21
		
		if day == fourthThursday {
			labels = append(labels, "Thanksgiving")
		}
    }

	// International games
	if !strings.HasSuffix(location, "USA") {
		labels = append(labels, "International")
	}

	if len(labels) == 0 {
		return ""
	}

	// Join with commas
	result := ""
	for i, label := range labels {
		if i > 0 {
			result += ","
		}
		result += label
	}

	return result
}

func parseBroadcast(broadcast []string) string {
    // Define the contains logic as a closure
    contains := func(slice []string, item string) bool {
        for _, s := range slice {
            if s == item {
                return true
            }
        }
        return false
    }

    network := ""
    priorities := []string{"NBC", "Prime Video", "ESPN", "ABC", "NFL Net", "CBS", "FOX", "Peacock"}
    
    for _, priority := range priorities {
        if contains(broadcast, priority) {
            network = priority
            break
        }
    }

    return network
}