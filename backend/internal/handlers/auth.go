// User registration, login, and profile management handlers

package handlers

import (
    "os"
    "regexp"
    "time"
    "unicode"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "gamescript/internal/database"
)


func validatePassword(password string) []string {
    var errors []string
    
    // Check length of password
    if len(password) < 8 {
        errors = append(errors, "Password must be at least 8 characters")
    }

    // Check for character types
	var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, char := range password {
        switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
        }
    }
    if !hasUpper {
        errors = append(errors, "Password must contain at least one uppercase letter")
    }
    if !hasLower {
        errors = append(errors, "Password must contain at least one lowercase letter")
    }
    if !hasNumber {
        errors = append(errors, "Password must contain at least one number")
    }
    if !hasSpecial {
        errors = append(errors, "Password must contain at least one special character")
    }
    
    return errors
}

func validateEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

func validateUsername(username string) []string {
    var errors []string
    
    if len(username) < 3 {
        errors = append(errors, "Username must be at least 3 characters")
    }
    if len(username) > 50 {
        errors = append(errors, "Username must be less than 50 characters")
    }
    
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    if !usernameRegex.MatchString(username) {
        errors = append(errors, "Username can only contain letters, numbers, hyphens, and underscores")
    }
    
    return errors
}

func RegisterUser(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        type RegisterRequest struct {
            Email    string `json:"email"`
            Username string `json:"username"`
            Password string `json:"password"`
        }

        var req RegisterRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }

        // Validate all fields
        var errors []string
        if req.Email == "" || req.Username == "" || req.Password == "" {
            return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
        }
        if !validateEmail(req.Email) {
            errors = append(errors, "Invalid email format")
        }
        errors = append(errors, validateUsername(req.Username)...)
        errors = append(errors, validatePassword(req.Password)...)
        if len(errors) > 0 {
            return c.Status(400).JSON(fiber.Map{"error": errors[0], "errors": errors})
        }

		// Hash password with higher cost for production
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
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
                "id":         id,
                "email":      email,
                "username":   username,
                "is_admin":   isAdmin,
                "created_at": createdAt,
            },
            "token": token,
        })
    }
}

func LoginUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

        // Validate input
		if req.Email == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Email and password are required"})
		}

		// Check for account lockout
		var userID int
		var failedAttempts int
		var lockedUntil *time.Time
		lockCheckQuery := `
			SELECT id, failed_login_attempts, locked_until
			FROM users
			WHERE email = $1
		`
		err := db.Conn.QueryRow(lockCheckQuery, req.Email).Scan(&userID, &failedAttempts, &lockedUntil)

		if err == nil && lockedUntil != nil && time.Now().Before(*lockedUntil) {
            remainingTime := time.Until(*lockedUntil).Minutes()
            return c.Status(423).JSON(fiber.Map{
                "error": "Account is temporarily locked due to multiple failed login attempts. Please try again later.",
                "locked_for_minutes": int(remainingTime) + 1,
            })
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

		err = db.Conn.QueryRow(query, req.Email).Scan(
			&id, &email, &username, &passwordHash, &isAdmin, &createdAt,
		)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
            // Increment failed attempts
            failedAttempts++
            var lockTime *time.Time
            
            if failedAttempts >= 5 {
                // Lock account for 15 minutes after 5 failed attempts
                lockDuration := time.Now().Add(15 * time.Minute)
                lockTime = &lockDuration
            }

            updateQuery := `
                UPDATE users
                SET failed_login_attempts = $1, locked_until = $2
                WHERE id = $3
            `
            db.Conn.Exec(updateQuery, failedAttempts, lockTime, id)

            if failedAttempts >= 5 {
                return c.Status(423).JSON(fiber.Map{
                    "error": "Too many failed attempts. Account locked for 15 minutes.",
                })
            }

            return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
        }

		// Reset failed attempts and update last login on successful login
        resetQuery := `
            UPDATE users
            SET failed_login_attempts = 0, locked_until = NULL, last_login = $1
            WHERE id = $2
        `
        db.Conn.Exec(resetQuery, time.Now(), id)

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

func GetCurrentUser(db *database.DB) fiber.Handler {
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

func UpdateProfile(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(int)

        type UpdateProfileRequest struct {
            Username       *string `json:"username"`
            Email          *string `json:"email"`
            CurrentPassword *string `json:"current_password"`
            NewPassword    *string `json:"new_password"`
        }

        var request UpdateProfileRequest
        if err := c.BodyParser(&request); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        // If changing password, verify current password
        if request.NewPassword != nil && *request.NewPassword != "" {
            if request.CurrentPassword == nil || *request.CurrentPassword == "" {
                return c.Status(400).JSON(fiber.Map{"error": "Current password is required to change password"})
            }

            // Validate new password
            if errors := validatePassword(*request.NewPassword); len(errors) > 0 {
                return c.Status(400).JSON(fiber.Map{"error": errors[0], "errors": errors})
            }

            // Get current password hash
            var currentHash string
            err := db.Conn.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&currentHash)
            if err != nil {
                return c.Status(500).JSON(fiber.Map{"error": "Failed to verify password"})
            }

            // Verify current password
            if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(*request.CurrentPassword)); err != nil {
                return c.Status(401).JSON(fiber.Map{"error": "Current password is incorrect"})
            }

            // Hash new password
            hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.NewPassword), 12)
            if err != nil {
                return c.Status(500).JSON(fiber.Map{"error": "Failed to hash new password"})
            }

            // Update password
            _, err = db.Conn.Exec(`
                UPDATE users
                SET password_hash = $1, password_changed_at = $2, updated_at = $3
                WHERE id = $4
            `, string(hashedPassword), time.Now(), time.Now(), userID)
            if err != nil {
                return c.Status(500).JSON(fiber.Map{"error": "Failed to update password"})
            }
        }

        // Update username if provided
        if request.Username != nil && *request.Username != "" {
            if errors := validateUsername(*request.Username); len(errors) > 0 {
                return c.Status(400).JSON(fiber.Map{"error": errors[0], "errors": errors})
            }

            _, err := db.Conn.Exec("UPDATE users SET username = $1, updated_at = $2 WHERE id = $3", *request.Username, time.Now(), userID)
            if err != nil {
                if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
                    return c.Status(400).JSON(fiber.Map{"error": "Username already in use"})
                }
                return c.Status(500).JSON(fiber.Map{"error": "Failed to update username"})
            }
        }

        // Update email if provided
        if request.Email != nil && *request.Email != "" {
            if !validateEmail(*request.Email) {
                return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
            }

            _, err := db.Conn.Exec("UPDATE users SET email = $1, updated_at = $2 WHERE id = $3", *request.Email, time.Now(), userID)
            if err != nil {
                if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
                    return c.Status(400).JSON(fiber.Map{"error": "Email already in use"})
                }
                return c.Status(500).JSON(fiber.Map{"error": "Failed to update email"})
            }
        }

        // Return updated user
        var id int
        var email, username string
        var isAdmin bool
        var avatarURL *string
        var createdAt, updatedAt time.Time

        err := db.Conn.QueryRow(`
            SELECT id, email, username, is_admin, avatar_url, created_at, updated_at
            FROM users
            WHERE id = $1
        `, userID).Scan(&id, &email, &username, &isAdmin, &avatarURL, &createdAt, &updatedAt)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve user"})
        }

        return c.JSON(fiber.Map{
            "id":             id,
            "email":          email,
            "username":       username,
            "is_admin":      isAdmin,
            "avatar_url":    avatarURL,
            "created_at":    createdAt,
            "updated_at":    updatedAt,
        })
    }
}

// Helper function to generate JWT token
func generateJWT(userID int, email, username string) string {
    // Get JWT secret from environment variables
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "jwt_dev_secret_key"
    }

    // Set token claims
    claims := jwt.MapClaims{
        "user_id":  userID,
        "email":    email,
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(jwtSecret))

    return tokenString
}