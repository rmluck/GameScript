package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gamescript/internal/database"
	"gamescript/internal/handlers"
	"gamescript/internal/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Start background scheduler
	scheduler := scheduler.NewScheduler(db)
	scheduler.Start()
	defer scheduler.Stop()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "GameScript API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup routes
	handlers.SetupRoutes(app, db, scheduler)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")
		scheduler.Stop()
		_ = app.Shutdown()
	}()

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}