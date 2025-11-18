package handlers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gamescript/internal/database"
)

func registerUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type RegisterRequest struct {
			Email   	string `json:"email"`
			Username 	string `json:"username"`
			Password 	string `json:"password"`
		}

		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		if req.Email == "" || req.Username == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
		}

		if len(req.Password) < 8 {
			return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		query := `
			INSERT INTO users (email, username, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id, email, username, is_admin, created_at
		`

		var id int
		var email, username string
		var isAdmin bool
		var createdAt time.Time

		err = db.Conn.QueryRow(query, req.Email, req.Username, string(hashedPassword)).Scan(
			&id, &email, &username, &isAdmin, &createdAt,
		)
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
				return c.Status(400).JSON(fiber.Map{"error": "Email already in use"})
			}
			if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
				return c.Status(400).JSON(fiber.Map{"error": "Username already in use"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
		}

		token := generateJWT(id, email, username)

		return c.Status(201).JSON(fiber.Map{
			"user": map[string]interface{}{
				"id": id,
				"email": email,
				"username": username,
				"is_admin": isAdmin,
				"created_at": createdAt,
			},
			"token": token,
		})
	}
}

func loginUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if req.Email == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Email and password are required"})
		}

		query := `
			SELECT id, email, username, password_hash, is_admin, created_at
			FROM users
			WHERE email = $1
		`

		var id int
		var email, username, passwordHash string
		var isAdmin bool
		var createdAt time.Time

		err := db.Conn.QueryRow(query, req.Email).Scan(
			&id, &email, &username, &passwordHash, &isAdmin, &createdAt,
		)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
		}

		token := generateJWT(id, email, username)

		return c.JSON(fiber.Map{
			"user": map[string]interface{}{
				"id": id,
				"email": email,
				"username": username,
				"is_admin": isAdmin,
				"created_at": createdAt,
			},
			"token": token,
		})
	}
}

func getCurrentUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int)

		query := `
			SELECT id, email, username, is_admin, avatar_url, created_at, updated_at
			FROM users
			WHERE id = $1
		`

		var id int
		var email, username string
		var isAdmin bool
		var avatarURL *string
		var createdAt, updatedAt time.Time

		err := db.Conn.QueryRow(query, userID).Scan(
			&id, &email, &username, &isAdmin, &avatarURL, &createdAt, &updatedAt,
		)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}

		return c.JSON(map[string]interface{}{
			"id": id,
			"email": email,
			"username": username,
			"is_admin": isAdmin,
			"avatar_url": avatarURL,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}
}

func generateJWT(userID int, email, username string) string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "jwt_dev_secret_key"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email": email,
		"username": username,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtSecret))

	return tokenString
}