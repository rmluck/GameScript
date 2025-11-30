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

// Cleanup old entries every 5 minutes
func init() {
    go cleanupVisitors()
}

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

// RateLimitAuth limits login/register attempts to prevent brute force
func RateLimitAuth(maxRequests int, window time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ip := c.IP()
        
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