package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

func main() {
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

	channel := "rag:ingestion"

	pubsub := client.Subscribe(ctx, channel)

	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
		log.Fatalf("Subscription failed: %v", err)
	}

	log.Printf("Subscribed to channel: %s", channel)

	messages := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Subscriber shutting down...")
			return
		case msg, ok := <-messages:
			if !ok {
				log.Println("Subscrption channel closed.")
				return
			}

			log.Printf(
				"Recived message on %s:%s",
				msg.Channel,
				msg.Payload,
			)
		}
	}
}
