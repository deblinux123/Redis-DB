package main

import (
	"context"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-05/internal/redisclient"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system envirnoment.")
	}

	cfg, err := redisclient.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	client, err := redisclient.Connect(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	defer client.Close()

	log.Println("✅Redis Connected Successfully.")

	if err := client.Set(
		ctx,
		"day05:practice",
		"hello redis",
		0,
	).Err(); err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	log.Println("Value stored successfully.")

	value, err := client.Get(
		ctx,
		"day05:practice",
	).Result()

	if err != nil {
		log.Fatalf("Failed to get value")
	}

	log.Println("Retrive value:", value)
}
