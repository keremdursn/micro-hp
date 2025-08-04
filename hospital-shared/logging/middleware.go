package logging

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CorrelationIDMiddleware adds correlation ID to request context
func CorrelationIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get or generate correlation ID
		correlationID := c.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Set correlation ID in response header
		c.Set("X-Correlation-ID", correlationID)

		// Add correlation ID to context
		ctx := context.WithValue(c.Context(), CorrelationIDKey, correlationID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// RequestLoggingMiddleware logs all HTTP requests
func RequestLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Extract user ID if available (from JWT)
		var userID *uint
		if user := c.Locals("user"); user != nil {
			if userMap, ok := user.(map[string]interface{}); ok {
				if id, exists := userMap["user_id"]; exists {
					if uid, ok := id.(uint); ok {
						userID = &uid
					}
				}
			}
		}

		// Log the request
		if GlobalLogger != nil {
			GlobalLogger.LogRequest(
				c.UserContext(),
				c.Method(),
				c.Path(),
				c.Response().StatusCode(),
				duration,
				userID,
			)
		}

		return err
	}
}

// ErrorLoggingMiddleware logs application errors
func ErrorLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil && GlobalLogger != nil {
			GlobalLogger.LogError(
				c.UserContext(),
				err,
				"Request error occurred",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Int("status", c.Response().StatusCode()),
			)
		}

		return err
	}
}
