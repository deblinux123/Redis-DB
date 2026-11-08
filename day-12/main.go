package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/deblinux123/Redis-DB/day-12/distributedlock"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Connected to Redis.")

	const workers = 5

	var wg sync.WaitGroup

	for i := 1; i <= workers; i++ {
		wg.Add(1)

		go func(workerId int) {
			defer wg.Done()

			// create new lock
			lock, err := distributedlock.NewLock(
				rdb,
				"lock:order:123",
				10*time.Second,
			)

			if err != nil {
				log.Printf("Worker %d: %v\n", workerId, err)
				return
			}

			acquire, err := lock.Acquire(ctx)

			if err != nil {
				log.Printf("Worker: %d acquire error: %v\n", workerId, err)
				return
			}

			if !acquire {
				fmt.Printf("Worker: %d lock is busy.\n", workerId)
				return
			}

			log.Printf("Worker: %d lock acquired\n", workerId)

			time.Sleep(3 * time.Second)

			released, err := lock.Release(ctx)

			if err != nil {
				log.Printf("Worker: %d release error: %v", workerId, err)
				return
			}

			log.Printf(
				"Worker: %d lock release = %v\n",
				workerId,
				released,
			)
		}(i)
	}

	wg.Wait()

	log.Println("All Workes finished.")

	// workerA, _ := distributedlock.NewLock(
	// 	rdb,
	// 	"lock:test",
	// 	2*time.Second,
	// )

	// ok, err := workerA.Acquire(ctx)

	// if err != nil || !ok {
	// 	log.Fatal("Worker A could not acquire lock:", err)
	// }

	// log.Println("Worker A acquired lock.")

	// time.Sleep(3 * time.Second)

	// workerB, _ := distributedlock.NewLock(
	// 	rdb,
	// 	"lock:test",
	// 	10*time.Second,
	// )

	// ok, err = workerB.Acquire(ctx)

	// if err != nil || !ok {
	// 	log.Fatal("Worker B could not acquire lock:", err)
	// }

	// log.Println("Worker B acquired lock.")

	// released, err := workerA.Release(ctx)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Worker A release result:", released)

	// renewed, err := workerB.Renew(ctx)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Worker B still owns lock:", renewed)

	// workerB.Release(ctx)
}
