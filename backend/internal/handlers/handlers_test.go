package handlers

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http/httptest"
    "testing"

    "gamescript/internal/database"
    "gamescript/internal/scheduler"

    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "github.com/stretchr/testify/assert"
)

// Load .env once before all tests
func init() {
    if err := godotenv.Load("../../.env"); err != nil {
        log.Println("Warning: No .env file found for tests, using environment variables")
    }
}

// Test helper to setup test app
func setupTestApp(t *testing.T) (*fiber.App, *database.DB) {
    // Connect to test database
    db, err := database.NewConnection()
    if err != nil {
        t.Fatal("Failed to connect to test database:", err)
    }

    app := fiber.New()
    scheduler := scheduler.NewScheduler(db)
    SetupRoutes(app, db, scheduler)

    return app, db
}

func TestHealthCheck(t *testing.T) {
    app, _ := setupTestApp(t)

    req := httptest.NewRequest("GET", "/api/health", nil)
    resp, _ := app.Test(req)

    assert.Equal(t, 200, resp.StatusCode)

    var body map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&body)
    assert.Equal(t, "ok", body["status"])
}

func TestGetSports(t *testing.T) {
    app, _ := setupTestApp(t)

    req := httptest.NewRequest("GET", "/api/sports", nil)
    resp, _ := app.Test(req)

    assert.Equal(t, 200, resp.StatusCode)

    var sports []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&sports)
    assert.Greater(t, len(sports), 0, "Should return at least one sport")
}

func TestRegisterUser(t *testing.T) {
    app, db := setupTestApp(t)
    defer db.Close()

    // Clean up test user if exists
    db.Conn.Exec("DELETE FROM users WHERE email = 'integration_test@example.com'")

    payload := map[string]interface{}{
        "email":    "integration_test@example.com",
        "username": "integrationtest",
        "password": "password123",
    }
    jsonPayload, _ := json.Marshal(payload)

    req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonPayload))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)

    assert.Equal(t, 201, resp.StatusCode)

    var body map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&body)
    assert.NotNil(t, body["token"])
    assert.NotNil(t, body["user"])
    
    // Cleanup
    db.Conn.Exec("DELETE FROM users WHERE email = 'integration_test@example.com'")
}

func TestRegisterUserDuplicateEmail(t *testing.T) {
    app, db := setupTestApp(t)
    defer db.Close()

    // Clean up first
    db.Conn.Exec("DELETE FROM users WHERE email = 'duplicate_test@example.com'")

    // Register first user
    payload := map[string]interface{}{
        "email":    "duplicate_test@example.com",
        "username": "duplicatetest1",
        "password": "password123",
    }
    jsonPayload, _ := json.Marshal(payload)
    req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonPayload))
    req.Header.Set("Content-Type", "application/json")
    app.Test(req)

    // Try to register with same email
    payload2 := map[string]interface{}{
        "email":    "duplicate_test@example.com",
        "username": "duplicatetest2",
        "password": "password123",
    }
    jsonPayload2, _ := json.Marshal(payload2)
    req2 := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonPayload2))
    req2.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req2)

    assert.Equal(t, 400, resp.StatusCode)

    var body map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&body)
    assert.Contains(t, body["error"], "Email already")

    // Cleanup
    db.Conn.Exec("DELETE FROM users WHERE email = 'duplicate_test@example.com'")
}

func TestLoginUser(t *testing.T) {
    app, db := setupTestApp(t)
    defer db.Close()

    // Clean up first
    db.Conn.Exec("DELETE FROM users WHERE email = 'login_test@example.com'")

    // Register user first
    registerPayload := map[string]interface{}{
        "email":    "login_test@example.com",
        "username": "logintest",
        "password": "password123",
    }
    jsonRegister, _ := json.Marshal(registerPayload)
    reqRegister := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonRegister))
    reqRegister.Header.Set("Content-Type", "application/json")
    app.Test(reqRegister)

    // Now login
    loginPayload := map[string]interface{}{
        "email":    "login_test@example.com",
        "password": "password123",
    }
    jsonLogin, _ := json.Marshal(loginPayload)
    reqLogin := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonLogin))
    reqLogin.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(reqLogin)

    assert.Equal(t, 200, resp.StatusCode)

    var body map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&body)
    assert.NotNil(t, body["token"])

    // Cleanup
    db.Conn.Exec("DELETE FROM users WHERE email = 'login_test@example.com'")
}

func TestCreateScenarioUnauthorized(t *testing.T) {
    app, _ := setupTestApp(t)

    payload := map[string]interface{}{
        "name":      "Test Scenario",
        "sport_id":  1,
        "season_id": 1,
        "is_public": true,
    }
    jsonPayload, _ := json.Marshal(payload)

    req := httptest.NewRequest("POST", "/api/scenarios", bytes.NewBuffer(jsonPayload))
    req.Header.Set("Content-Type", "application/json")
    // No Authorization header
    resp, _ := app.Test(req)

    // Should still work (guest scenario)
    assert.Equal(t, 201, resp.StatusCode)
}

func TestInvalidScenarioID(t *testing.T) {
    app, _ := setupTestApp(t)

    req := httptest.NewRequest("GET", "/api/scenarios/99999", nil)
    resp, _ := app.Test(req)

    assert.Equal(t, 404, resp.StatusCode)
}