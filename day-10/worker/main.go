package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	streamName = "rag:jobs"
	groupName  = "rag-workers"
)

func main() {
	workerID := "worker-1"

	if len(os.Args) > 1 {
		workerID = os.Args[1]
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	err := client.XGroupCreate(
		ctx,
		streamName,
		groupName,
		"0",
	).Err()

	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("failed to create consumer group: %v", err)
	}

	log.Printf("%s started successfully", workerID)

	for {
		if ctx.Err() != nil {
			log.Println("Worker shutting donw...")
			return
		}

		streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: workerID,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    2 * time.Second,
		}).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}

			if ctx.Err() != nil {
				log.Println("Worker shutting down...")
				return
			}

			log.Printf("Failed to read jobs: %v", err)
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				log.Printf(
					"%s recived job %s: %v",
					workerID,
					message.ID,
					message.Values,
				)

				if err := processJob(message.Values); err != nil {
					log.Printf("Job processing faile: %v", err)
					continue
				}

				// acknowledge successfull processing
				ack, err := client.XAck(
					ctx,
					streamName,
					groupName,
					message.ID,
				).Result()

				if err != nil {
					log.Fatalf("Failed to acknowledge job %s : %v", message.ID, err)
					continue
				}

				log.Printf(
					"%s complited job %s(ACk: %d)",
					workerID,
					message.ID,
					ack,
				)
			}
		}
	}
}

func processJob(values map[string]interface{}) error {
	paperID, ok := values["paper_id"].(string)

	if !ok || paperID == "" {
		return fmt.Errorf("Invalid paper id")
	}

	log.Printf("Processing paper: %s", paperID)

	time.Sleep(5 * time.Second)

	log.Printf("Paper processed successfully: %s", paperID)

	return nil
}
