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

	if err := client.Del(ctx, "day05:tasks").Err(); err != nil {
		log.Fatalf("Failed to delete: %v", err)
	}

	log.Println("Tasks Deleted.")

	err = client.RPush(ctx, "day05:tasks", "generate-report", "send-email", "resize-image").Err()

	if err != nil {
		log.Fatalf("Failed to create tasks")
	}

	log.Println("Tasks are created.")

	log.Println("Tasks:")

	var tasks []string

	tasks, err = client.LRange(ctx, "day05:tasks", 0, -1).Result()

	if err != nil {
		log.Fatalf("Failed to get tasks: %v", err)
	}

	for i, task := range tasks {
		log.Printf("%d- %s", i, task)
	}

	processTask, err := client.LPop(ctx, "day05:tasks").Result()

	if err != nil {
		log.Fatalf("There is no task: %v", err)
	}

	log.Printf("Processed task: %s", processTask)

	remain, err := client.LLen(ctx, "day05:tasks").Result()

	if err != nil {
		log.Fatalf("failed to get tasks")
	}

	log.Println("Remaining tasks:", remain)

	var remainTasks []string
	for _, task := range tasks {
		remainTasks = append(remainTasks, task)
	}

	log.Println("Tasks left:", remainTasks)
}
