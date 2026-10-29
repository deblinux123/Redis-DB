package main

import (
	"context"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-08/internal/redisclient"
	"github.com/joho/godotenv"
)

func main() {
	// loading config
	if err := godotenv.Load("./.env"); err != nil {
		log.Printf("No .env file found; using system envirnoment.")
	}

	cfg, err := redisclient.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load redis config: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	client, err := redisclient.Connect(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer client.Close()

	log.Println("Connected to Redis successfully.")
}
