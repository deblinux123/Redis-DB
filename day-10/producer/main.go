package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	// connect to redis client
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	defer client.Close()

	// check connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	// create stream
	stream := "rag:jobs"

	for i := 1; i <= 5; i++ {
		job := map[string]interface{}{
			"event":    "paper.ingest",
			"paper_id": "paper-" + time.Now().Format("150405") + "-" + string(rune('0'+i)),
			"status":   "queued",
		}

		id, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			Values: job,
		}).Result()

		if err != nil {
			log.Fatalf("Failed to publish job: %v", err)
			continue
		}

		log.Printf("Job published successfully. ID: %s", id)
	}
}
