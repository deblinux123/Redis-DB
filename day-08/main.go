package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-08/internal/redisclient"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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

	benchMarkCommands(ctx, client, 100000)

	if err := runTransaction(ctx, client); err != nil {
		log.Fatal(err)
	}

	if err := runTransacionERR(ctx, client); err != nil {
		log.Fatal(err)
	}

	if err := runWatchTest(ctx, client); err != nil {
		log.Printf("Watch Test: %v\n", err)
	}

	if err := runWathConflictTest(ctx, client); err != nil {
		log.Printf("Watch confilct test: %v", err)
	}
}

func benchMarkCommands(ctx context.Context, client *redis.Client, count int) {
	start := time.Now()

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("day08:bench:%d", i)

		if err := client.Set(ctx, key, i, 0).Err(); err != nil {
			log.Fatalf("Individual SET failed: %v", err)
		}
	}

	individualDuration := time.Since(start)

	pipe := client.Pipeline()

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("day08:bench:%d", i)
		pipe.Del(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Fatalf("cleaninup failed: %v", err)
	}

	start = time.Now()

	pipe = client.Pipeline()
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("day08:bench:%d", i)
		pipe.Set(ctx, key, i, 0)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Fatalf("Pipeline SET failed: %v", err)
	}

	pipeLineDuration := time.Since(start)

	log.Printf("Individual commands: %v\n", individualDuration)
	log.Printf("Pipelined command: %v\n", pipeLineDuration)

	if pipeLineDuration > 0 {
		log.Printf("Speedup: %.2fx\n", float64(individualDuration)/float64(pipeLineDuration))
	}

}

func runTransaction(ctx context.Context, client *redis.Client) error {
	faridKey := "day08:account:farid"
	asmaKey := "day08:account:asma"

	if err := client.Set(ctx, faridKey, 100, 0).Err(); err != nil {
		return err
	}

	if err := client.Set(ctx, asmaKey, 50, 0).Err(); err != nil {
		return err
	}

	// execute both
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.DecrBy(ctx, faridKey, 20)
		pipe.IncrBy(ctx, asmaKey, 20)

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	faridBalance, err := client.Get(ctx, faridKey).Int()

	if err != nil {
		return err
	}

	asmaBalance, err := client.Get(ctx, asmaKey).Int()

	if err != nil {
		return err
	}

	log.Printf("Farid's balance: %d\n", faridBalance)
	log.Printf("Asma's balance: %d\n", asmaBalance)

	return nil
}

func runTransacionERR(ctx context.Context, client *redis.Client) error {
	counterKey := "day08:transaction:counter"
	nameKey := "day08:transaction:name"

	if err := client.Set(ctx, counterKey, 10, 0).Err(); err != nil {
		return err
	}

	if err := client.Set(ctx, nameKey, "farid", 0).Err(); err != nil {
		return err
	}

	cmds, err := client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Incr(ctx, counterKey)
		p.Incr(ctx, nameKey)
		p.Incr(ctx, counterKey)

		return nil
	})

	log.Println("-\n--- Transaction Error Test ---")

	if err != nil {
		log.Printf("transaction returned an error: %v\n", err)
	} else {
		log.Println("transactoin returned no error.")
	}

	for i, cmd := range cmds {
		// log.Printf("Command %d: %v", i+1, cmd.Err())
		if cmd.Err() != nil {
			log.Printf("Command %d failed: %v\n", i+1, err)
			continue
		}

		log.Printf("Command %d succeeded.\n", i+1)
	}

	counter, getErr := client.Get(ctx, counterKey).Int()
	if getErr != nil {
		return err
	}

	log.Printf("Final counter value: %v\n", counter)

	return nil
}

func runWatchTest(ctx context.Context, client *redis.Client) error {

	key := "day08:watch:balance"

	if err := client.Set(ctx, key, 100, 0).Err(); err != nil {
		return err
	}

	err := client.Watch(ctx, func(tx *redis.Tx) error {
		balance, err := tx.Get(ctx, key).Int()

		if err != nil {
			return err
		}

		log.Printf("Balance befor transaction: %d\n", balance)

		if balance < 80 {
			return fmt.Errorf("insufficient balance.")
		}

		_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			p.DecrBy(ctx, key, 80)
			return nil
		})

		return nil
	}, key)

	if err != nil {
		return fmt.Errorf("Watch transaction failed: %v", err)
	}

	balance, err := client.Get(ctx, key).Int()

	if err != nil {
		return err
	}

	log.Printf("Balance after transaction: %d\n", balance)

	return nil
}

func runWathConflictTest(ctx context.Context, client *redis.Client) error {
	key := "day08:watch:conflict"

	if err := client.Set(ctx, key, 100, 0).Err(); err != nil {
		return err
	}

	err := client.Watch(ctx, func(tx *redis.Tx) error {
		balance, err := client.Get(ctx, key).Int()

		if err != nil {
			return err
		}

		log.Printf("Client A redis Balance: %d\n", balance)

		if err := client.Set(ctx, key, 200, 0).Err(); err != nil {
			return err
		}

		_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			p.DecrBy(ctx, key, 50)
			return nil
		})

		return err
	}, key)

	if errors.Is(err, redis.TxFailedErr) {
		log.Println("Transaction aboted: watch key changed.")
	} else if err != nil {
		log.Printf("Watch test failed: %v", err)
	} else {
		log.Println("Transaction complited.")
	}

	finalBalance, err := client.Get(ctx, key).Int()

	if err != nil {
		return err
	}

	log.Printf("Final Balance: %d\n", finalBalance)

	return nil
}
