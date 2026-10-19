package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-03/internal/redisclient"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file fund; using system envirnoment")
	}

	cfg, err := redisclient.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	client, err := redisclient.Connect(ctx, cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	// create user with set

	if err := client.Set(
		ctx,
		"user:1001:name",
		"farid",
		0,
	).Err(); err != nil {
		log.Fatal(err)
	}

	// set users email
	if err := client.Set(
		ctx,
		"user:1001:email",
		"farid@exmaple.com",
		0,
	).Err(); err != nil {
		log.Fatal(err)
	}

	// retrive the user's name

	name, err := client.Get(
		ctx,
		"user:1001:name",
	).Result()

	if errors.Is(err, redis.Nil) {
		fmt.Println("User not found")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("User:", name)
	}

	// check the email exists
	exists, err := client.Exists(
		ctx,
		"user:1001:email",
	).Result()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email exists:", exists)

	// updateing
	if err := client.Set(
		ctx,
		"user:1001:name",
		"farid kaki",
		0,
	).Err(); err != nil {
		log.Fatal(err)
	}

	// counter
	views, err := client.Incr(
		ctx,
		"page:home:views",
	).Result()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Page view:", views)

	// expiration
	if err := client.Set(
		ctx,
		"demo:temporary",
		"expire soon",
		10*time.Second,
	).Err(); err != nil {
		log.Fatal(err)
	}

	ttl, err := client.TTL(
		ctx,
		"demo:temporary",
	).Result()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Temporary key TTL:", ttl)

	// delete
	deleted, err := client.Del(
		ctx,
		"user:1001:email",
	).Result()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted email keys:", deleted)

}
