package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Extract user ID from JWT if present, otherwise generate/retrieve session token
func OptionalAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader != "" {
		// Try to authenticate with JWT
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_jwt_secret"), nil // Use env variable in production
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Locals("user_id", int(claims["user_id"].(float64)))
				c.Locals("is_authenticated", true)
				return c.Next()
			}
		}
	}

	// No valid JWT check for session token
	sessionToken := c.Cookies("session_token")

	if sessionToken == "" {
		// Generate new session token
		sessionToken = generateSessionToken()
		c.Cookie(&fiber.Cookie{
			Name: "session_token",
			Value: sessionToken,
			HTTPOnly: true,
			SameSite: "Lax",
			MaxAge: 60 * 60 * 24 * 30, // 30 days
		})
	}

	c.Locals("session_token", sessionToken)
	c.Locals("is_authenticated", false)

	return c.Next()
}

// Require authentication (for user-specific routes)
func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization required"})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_jwt_secret"), nil // Use env variable in production
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Locals("user_id", int(claims["user_id"].(float64)))
		c.Locals("is_authenticated", true)
		return c.Next()
	}

	return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
}

func generateSessionToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}