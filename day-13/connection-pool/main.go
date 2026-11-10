package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",

		PoolSize:     10,
		MinIdleConns: 2,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,

		MaxRetries: 2,
	})

	defer rdb.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	log.Println("Connected to Redis successfully.")

	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			key := fmt.Sprintf("day13:pool:%d", id)

			value := fmt.Sprintf("worker-%d", id)

			err := rdb.Set(
				ctx,
				key,
				value,
				time.Minute,
			).Err()

			if err != nil {
				log.Printf("Worker %d failed: %v", id, err)
			}

			log.Printf("Worker %d complited.", id)
		}(i)
	}

	wg.Wait()

	stats := rdb.PoolStats()

	fmt.Println("\n--- Connection Pool Stats ---")
	fmt.Printf("Total Connections: %d\n", stats.TotalConns)
	fmt.Printf("Idle Connections: %d\n", stats.IdleConns)
	fmt.Printf("Stale Connections: %d\n", stats.StaleConns)
	fmt.Printf("Hits: %d\n", stats.Hits)
	fmt.Printf("Misses: %d\n", stats.Misses)
	fmt.Printf("Timeouts: %d\n", stats.Timeouts)

	if err := rdb.Set(
		ctx,
		"day13:error:test",
		"hello from redis",
		time.Minute,
	).Err(); err != nil {
		log.Fatal(err)
	}

	value, err := rdb.Get(
		ctx,
		"day13:error:test",
	).Result()

	if err != nil {
		log.Printf("GET failed: %v", err)
	} else {
		log.Println("Existing key:", value)
	}

	value, err = rdb.Get(
		ctx,
		"day13:error:missing",
	).Result()

	if errors.Is(err, redis.Nil) {
		log.Println("Key does not exist.")
	} else if err != nil {
		log.Printf("Redis error: %v", err)
	} else {
		log.Println("Missing key value:", value)
	}
}
