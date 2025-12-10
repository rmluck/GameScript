package scheduler

import (
	"fmt"
	"log"
	"time"

	"gamescript/internal/models"
	"gamescript/internal/services/espn"
)

func (s *Scheduler) startNFLScheduler() {
	log.Println("Starting NFL scheduler...")

	// Calculate next midnight PST
	ticker := s.getNextMidnightPSTTicker()

	// Optional: Run immediately on startup
	// s.updateNFLSchedule()

	for {
		select {
		case <-ticker.C:
			s.updateNFLSchedule()

			// Reset ticker for next midnight PST
			ticker.Stop()
			ticker = s.getNextMidnightPSTTicker()
		
		case <-s.quit:
			ticker.Stop()
			log.Println("NFL scheduler stopped.")
			return
		}
	}
}

func (s *Scheduler) getNextMidnightPSTTicker() *time.Ticker {
	pstLocation, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Printf("Error loading PST timezone: %v, using UTC", err)
		pstLocation = time.UTC
	}

	now := time.Now().In(pstLocation)
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, pstLocation)
	durationUntilMidnight := time.Until(nextMidnight)

	log.Printf("Next NFL update scheduled for: %v (in %v)", nextMidnight, durationUntilMidnight)

	return time.NewTicker(durationUntilMidnight)
}

func (s *Scheduler) updateNFLSchedule() {
	log.Println("Starting NFL schedule update...")
	startTime := time.Now()

	client := espn.NewClient()
	currentYear := time.Now().Year()

	// Fetch entire NFL season
	games, err := client.FetchEntireNFLSeason(currentYear)
	if err != nil {
		log.Printf("Error fetching NFL schedule: %v", err)
		return
	}

	// Update games in database
	updated := 0
	errors := 0

	for _, game := range games {
		if err := s.updateNFLGame(game); err != nil {
			log.Printf("Error updating NFL game %s: %v", game.ESPNID, err)
			errors++
			continue
		}
		updated++
	}

	duration := time.Since(startTime)
	log.Printf("NFL schedule update completed in %v: %d games updated, %d errors", duration, updated, errors)
}

func (s *Scheduler) updateNFLGame(game models.Game) error {
    stmt := `
        INSERT INTO games (
            season_id, espn_id, home_team_id, away_team_id, start_time,
            day_of_week, week, location, primetime, network,
            home_score, away_score, status, is_postseason
        ) VALUES (
         	$1, $2,
            (SELECT id FROM teams WHERE season_id = $1 AND espn_id = $3),
            (SELECT id FROM teams WHERE season_id = $1 AND espn_id = $4),
            $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )
        ON CONFLICT (season_id, espn_id) DO UPDATE SET
            start_time = EXCLUDED.start_time,
            day_of_week = EXCLUDED.day_of_week,
            week = EXCLUDED.week,
            location = EXCLUDED.location,
            primetime = EXCLUDED.primetime,
            network = EXCLUDED.network,
            home_score = EXCLUDED.home_score,
            away_score = EXCLUDED.away_score,
            status = EXCLUDED.status
    `

    // Only set scores if status is "final"
    var homeScore, awayScore *int
    if game.Status != nil && *game.Status == "final" {
        homeScore = game.HomeScore
        awayScore = game.AwayScore
    }

    result, err := s.db.Conn.Exec(
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
        homeScore,   // Will be NULL for upcoming games
        awayScore,   // Will be NULL for upcoming games
        game.Status,
        game.IsPostseason,
    )

    if err != nil {
        return fmt.Errorf("database error: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error checking rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no game found with espn_id %s", game.ESPNID)
    }

    return nil
}

// Public method for manual triggering
func (s *Scheduler) UpdateNFLSchedule() {
	go s.updateNFLSchedule()
}