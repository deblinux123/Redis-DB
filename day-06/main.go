package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-06/internal/redisclient"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system envirnoment.")
	}

	cfg, err := redisclient.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	client, err := redisclient.Connect(ctx, cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	log.Println("Redis connected successfully.")

	key := "day06:cache:message"
	value := "Hello from go"

	err = client.Set(ctx, key, value, 10*time.Second).Err()

	if err != nil {
		log.Fatalf("Failed to set cache: %v", err)
	}

	log.Println("Cache stored successfully.")

	res, err := client.Get(ctx, key).Result()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
	time.Sleep(5 * time.Second)

	ttl, err := client.TTL(ctx, key).Result()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Remaining TTL:", ttl)

	userID := "1002"

	userData, err := getUser(ctx, client, userID)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User: Name=%s, Age=%s", userData.Name, userData.Age)

}

type userInformation struct {
	Name string
	Age  string
}

func fetchUserFromDB(userID string) userInformation {
	fmt.Printf("Fetching user %s from database...\n", userID)

	userInfo := userInformation{
		Name: "farid",
		Age:  "32",
	}

	return userInfo
}

func getUser(ctx context.Context, client *redis.Client, userID string) (userInformation, error) {

	key := "user:" + userID

	userData, err := client.HGetAll(ctx, key).Result()

	if err != nil {
		return userInformation{}, fmt.Errorf("failed to get user.")
	}

	if len(userData) > 0 {
		fmt.Println("Cache HIT")

		return userInformation{
			Name: userData["Name"],
			Age:  userData["Age"],
		}, nil
	}

	fmt.Println("Cache MISS")

	user := fetchUserFromDB(userID)

	err = client.HSet(ctx, key, map[string]interface{}{
		"Name": user.Name,
		"Age":  user.Age,
	}).Err()

	if err != nil {
		return userInformation{}, fmt.Errorf("failed to cache user: %w", err)
	}

	err = client.Expire(ctx, key, 60*time.Second).Err()

	if err != nil {
		return userInformation{}, fmt.Errorf("Failed to set cache expiration: %w", err)
	}

	return user, err

}
