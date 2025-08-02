package middleware

import (
	"time"

	"hospital-service/pkg/metrics"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// AuthRateLimiter Auth endpointleri için rate limiter
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,               // 5 istek
		Expiration: 1 * time.Minute, // 1 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // IP bazlı limit
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("auth", c.IP()).Inc()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	})
}

// LoginRateLimiter Login endpointi için daha sıkı limit
func LoginRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        3,               // 3 istek
		Expiration: 5 * time.Minute, // 5 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":login" // IP + endpoint bazlı
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("login", c.IP()).Inc()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many login attempts. Please try again in 5 minutes.",
			})
		},
	})
}

// GeneralRateLimiter Genel endpointler için rate limiter
func GeneralRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // 100 istek
		Expiration: 1 * time.Minute, // 1 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("general", c.IP()).Inc()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}

// AdminRateLimiter Admin endpointleri için özel limit
func AdminRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,              // 20 istek
		Expiration: 1 * time.Minute, // 1 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":admin"
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("admin", c.IP()).Inc()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Admin rate limit exceeded.",
			})
		},
	})
}
