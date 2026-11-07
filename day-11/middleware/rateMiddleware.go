package middleware

import (
	"context"

	"github.com/deblinux123/Redis-DB/day-11/ratelimiter"
	"github.com/gofiber/fiber/v3"
)

func RateLimit(limiter *ratelimiter.Limiter) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()

		allowed, err := limiter.Allow(
			context.Background(),
			ip,
		)

		// fail open
		if err != nil {
			return c.Next()
		}

		// reject request
		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":    "Too Many Requests",
				"messages": "Rate Limit Exceeded",
			})
		}

		return c.Next()
	}
}
