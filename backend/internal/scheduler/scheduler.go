// For scheduler to run automatically, the backend needs to be deployed to a hosting service (Render, Fly.io, Railway, etc.). Can also be run on local server.

package scheduler

import (
	"log"

	"gamescript/internal/database"
)

type Scheduler struct {
	db *database.DB
	quit chan bool
}

func NewScheduler(db *database.DB) *Scheduler {
	return &Scheduler {
		db: db,
		quit: make(chan bool),
	}
}

func (s *Scheduler) Start() {
	log.Println("Starting sports schedulers...")

	// Start NFL scheduler
	go s.startNFLScheduler()

	// TODO: NBA scheduler
	// go s.startNBAScheduler()

	// TODO: CFB scheduler
	// go s.startCFBScheduler()
}

func (s *Scheduler) Stop() {
	log.Println("Stopping all schedulers...")
	close(s.quit)
}