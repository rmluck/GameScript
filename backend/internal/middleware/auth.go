// Middleware for JWT authentication in a Fiber web application
// TODO: Update to claim started scenario from non-user when logging in/signing up

package middleware

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   int    `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func AuthMiddleware(c *fiber.Ctx) error {
    // Get token from Authorization header
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{"error": "Missing authorization token"})
    }

    // Extract token from "Bearer <token>"
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization format. Use 'Bearer <token>'"})
    }

    // Get JWT secret
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "jwt_dev_secret_key"
    }

    // Parse and validate token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jwtSecret), nil
    })
    if err != nil {
        log.Printf("Token parsing error: %v", err)
        return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
    }

    // Extract claims
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        log.Printf("Invalid token claims or token not valid")
        return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
    }

    // Check if token is expired
    // if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
    //     return c.Status(401).JSON(fiber.Map{"error": "Token has expired"})
    // }

    // Store user info in context
    c.Locals("user_id", claims.UserID)
    c.Locals("email", claims.Email)
    c.Locals("username", claims.Username)
    c.Locals("is_authenticated", true)

    return c.Next()
}

// Allows requests with or without authentication
func OptionalAuth(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    
    // If no auth header, check for session token in cookie
    if authHeader == "" {
        sessionToken := c.Cookies("session_token")
        if sessionToken != "" {
            c.Locals("session_token", sessionToken)
        }
        c.Locals("is_authenticated", false)
        return c.Next()
    }

    // If auth header exists, validate it
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        c.Locals("is_authenticated", false)
        return c.Next() // Invalid format, but continue as guest
    }

    // Get JWT secret
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "jwt_dev_secret_key"
    }

    // Parse and validate token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jwtSecret), nil
    })

    // If token is valid, extract claims
    if err == nil && token.Valid {
        claims, ok := token.Claims.(*Claims)
        if ok {
            c.Locals("user_id", claims.UserID)
            c.Locals("email", claims.Email)
            c.Locals("username", claims.Username)
            c.Locals("is_authenticated", true)
            return c.Next()
        }
    }

    c.Locals("is_authenticated", false)
    return c.Next()
}