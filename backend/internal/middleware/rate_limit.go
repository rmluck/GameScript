// Rrate limiting middleware to prevent brute force attacks on authentication endpoints

package middleware

import (
    "sync"
    "time"

    "github.com/gofiber/fiber/v2"
)


type visitor struct {
    lastSeen time.Time
    count    int
}

var (
    visitors = make(map[string]*visitor)
    mu       sync.RWMutex
)

func init() {
    go cleanupVisitors()
}

// Periodically removes old entries from the visitors map every 5 minutes
func cleanupVisitors() {
    for {
        time.Sleep(5 * time.Minute)
        mu.Lock()
        for ip, v := range visitors {
            if time.Since(v.lastSeen) > 10*time.Minute {
                delete(visitors, ip)
            }
        }
        mu.Unlock()
    }
}

// Limits login/register attempts to prevent brute force
func RateLimitAuth(maxRequests int, window time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get client IP
        ip := c.IP()
        
        // Check and update visitor info
        mu.Lock()
        v, exists := visitors[ip]
        if !exists {
            visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
            mu.Unlock()
            return c.Next()
        }

        // Reset counter if window has passed
        if time.Since(v.lastSeen) > window {
            v.count = 1
            v.lastSeen = time.Now()
            mu.Unlock()
            return c.Next()
        }

        // Increment counter
        v.count++
        v.lastSeen = time.Now()
        
        // Check if limit exceeded
        if v.count > maxRequests {
            mu.Unlock()
            return c.Status(429).JSON(fiber.Map{
                "error": "Too many requests. Please try again later.",
            })
        }
        
        mu.Unlock()
        return c.Next()
    }
}