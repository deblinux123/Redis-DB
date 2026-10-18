package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-02/internal/redisclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found; using system environment")
	}

	cfg, err := redisclient.LoadConfig()

	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	client, err := redisclient.Connect(ctx, cfg)

	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	defer client.Close()

	fmt.Println("Redis connection successfull")

	opCtx, opCancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)

	defer opCancel()

	if err := client.Set(
		opCtx,
		"day02:status",
		"connected",
		30*time.Second,
	).Err(); err != nil {
		log.Fatalf("set faile: %v", err)
	}

	value, err := client.Get(ctx, "day02:status").Result()

	if err != nil {
		log.Fatalf("get failed: %v", err)
	}

	fmt.Println("Redis Value:", value)
}
