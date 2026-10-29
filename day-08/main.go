package main

import (
	"context"
	"fmt"
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

	// remove data if exists
	if err := client.Del(
		ctx,
		"day08:user:1",
		"day08:user:2",
		"day08:user:3",
	).Err(); err != nil {
		log.Fatalf("Cleanup failed: %v", err)
	}

	// create pipeline

	pipe := client.Pipeline()

	pipe.Set(ctx, "day08:user:1", "farid", 0)
	pipe.Set(ctx, "day08:user:2", "asma", 0)
	pipe.Set(ctx, "day08:user:3", "mobina", 0)

	// execute all queued command
	cmds, err := pipe.Exec(ctx)

	if err != nil {
		log.Fatalf("Pipeline executation failed: %v", err)
	}

	log.Printf("Executed %d commands.\n", len(cmds))

	// read ad verify the results
	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("day08:user:%d", i)

		value, err := client.Get(ctx, key).Result()

		if err != nil {
			log.Fatalf("Failed to read %s: %v", key, err)
		}

		log.Printf("%s = %s\n", key, value)
	}
}
