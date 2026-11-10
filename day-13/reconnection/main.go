package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6380",

		PoolSize:     10,
		MinIdleConns: 2,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,

		MaxRetries: -1,
	})

	defer rdb.Close()

	log.Println("Starting Redis connection test...")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := connectWithRetry(ctx, rdb); err != nil {
		fmt.Printf("Could not connect to Redis: %v", err)
		return
	}

	fmt.Println("Redis connection establlished.")

	if err := rdb.Set(
		ctx,
		"day13:retry:test",
		"connection-restored",
		time.Minute,
	).Err(); err != nil {
		log.Printf("SET failed: %v", err)
		return
	}

	value, err := rdb.Get(
		ctx,
		"day13:retry:test",
	).Result()

	if err != nil {
		log.Printf("GET failed: %v", err)
		return
	}

	log.Printf("Redis value: %s", value)

}

func connectWithRetry(ctx context.Context, rdb *redis.Client) error {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		log.Printf(
			"Connection attemt %d/%d...\n",
			attempt,
			maxAttempts,
		)

		attemptCtx, cancel := context.WithTimeout(
			ctx,
			2*time.Second,
		)

		err := rdb.Ping(attemptCtx).Err()

		cancel()

		if err == nil {
			log.Printf("Redis Connection successfully.")
			return nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Connection attemtp time out.")
		} else {
			fmt.Printf("Connection failed: %v\n", err)
		}

		if attempt == maxAttempts {
			return fmt.Errorf(
				"All %d connection attemepts failed: %v",
				maxAttempts,
				err,
			)
		}

		delay := time.Duration(attempt) * time.Second

		fmt.Printf("Retrying in %v...", delay)

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	return fmt.Errorf("connection attempts axhausted.")
}
