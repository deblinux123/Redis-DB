package main

import (
	"context"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-11/middleware"
	"github.com/deblinux123/Redis-DB/day-11/ratelimiter"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	log.Println("Connected to Redis.")

	limiter := ratelimiter.NewLimiter(
		rdb,
		5,
		60*time.Second,
	)

	app := fiber.New()

	app.Use(middleware.RateLimit(limiter))

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Request allowed",
			"ip":      c.IP(),
		})
	})

	log.Fatal(app.Listen(":8080"))
}
