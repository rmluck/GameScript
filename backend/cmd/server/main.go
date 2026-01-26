// Main application entry point

package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/joho/godotenv"

    "gamescript/internal/database"
    "gamescript/internal/handlers"
    "gamescript/internal/middleware"
    "gamescript/internal/scheduler"
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
        AppName:   "GameScript API",
        BodyLimit: 4 * 1024 * 1024, // 4 MB max body size
    })

    // Recovery middleware (must be first)
    app.Use(recover.New())

    // Logger middleware
    app.Use(logger.New())

    // CORS configuration
    allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
    if allowedOrigins == "" {
        allowedOrigins = "http://localhost:5173,http://localhost:3000"
    }

    // CORS middleware
    app.Use(cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // API group
    api := app.Group("/api")

    // Health check
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })

    // Auth routes with rate limiting
    auth := api.Group("/auth")
    auth.Post("/register", middleware.RateLimitAuth(5, 15*time.Minute), handlers.RegisterUser(db))
    auth.Post("/login", middleware.RateLimitAuth(5, 15*time.Minute), handlers.LoginUser(db))
    auth.Get("/me", middleware.AuthMiddleware, handlers.GetCurrentUser(db))
    auth.Put("/profile", middleware.AuthMiddleware, handlers.UpdateProfile(db))

    // Setup other routes
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

    // Start server
    log.Printf("Server starting on http://localhost:%s", port)
    log.Fatal(app.Listen(":" + port))
}