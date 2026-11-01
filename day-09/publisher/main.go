package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type IngestionEvent struct {
	Event   string `json:"event"`
	PaperID string `json:"paper_id"`
	Status  string `json:"status"`
}

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	event := IngestionEvent{
		Event:   "paper.ingested",
		PaperID: "paper-123",
		Status:  "queued",
	}

	payload, err := json.Marshal(event)

	if err != nil {
		log.Fatalf("failed to encode event: %v", err)
	}

	channel := "rag:ingestion"

	subscribers, err := client.Publish(
		ctx,
		channel,
		payload,
	).Result()

	if err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}

	log.Printf("Published event: %s", payload)
	log.Printf("redis reported %d reciving subscriber(s)", subscribers)
}
