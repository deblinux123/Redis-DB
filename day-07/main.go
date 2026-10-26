package main

import (
	"context"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-07/internal/redisclient"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environmet")
	}

	cfg, err := redisclient.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load Redis config: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	client, err := redisclient.Connect(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	defer client.Close()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Status":  "ok",
			"Message": "Server Is Running",
		})
	})

	app.Get("/health/redis", func(c *fiber.Ctx) error {
		if err := client.Ping(c.UserContext()).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"Status": "error",
				"Redis":  "Unavailable",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Staus": "Ok",
			"Redis": "Connected.",
		})
	})

	log.Fatal(app.Listen(":8080"))
}
