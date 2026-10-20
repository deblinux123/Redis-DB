package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/deblinux123/Redis-DB/day-04/internal/redisclient"
	"github.com/deblinux123/Redis-DB/day-04/userrepo"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	cfg, err := redisclient.LoadConfig()
	if err != nil {
		log.Fatalf("load Redis config: %v", err)
	}

	client, err := redisclient.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("connect to Redis: %v", err)
	}
	defer client.Close()

	repo := userrepo.New(client)

	user := userrepo.User{
		ID:    "1001",
		Name:  "Farid",
		Email: "farid@example.com",
		Age:   28,
		City:  "Baku",
	}

	// Clean up the previous demo user so the example is repeatable.
	if err := repo.Delete(ctx, user.ID); err != nil &&
		!errors.Is(err, userrepo.ErrUserNotFound) {
		log.Fatalf("clean previous user: %v", err)
	}

	// 1. Create user
	if err := repo.Create(ctx, user); err != nil {
		log.Fatalf("create user: %v", err)
	}
	fmt.Println("User created successfully")

	// 2. Get complete profile
	got, err := repo.Get(ctx, user.ID)
	if err != nil {
		log.Fatalf("get user: %v", err)
	}
	fmt.Printf("User: %+v\n", got)

	// 3. Update one field
	if err := repo.UpdateField(
		ctx,
		user.ID,
		"city",
		"Tbilisi",
	); err != nil {
		log.Fatalf("update city: %v", err)
	}
	fmt.Println("City updated successfully")

	// 4. Read the updated profile
	got, err = repo.Get(ctx, user.ID)
	if err != nil {
		log.Fatalf("get updated user: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", got)

	// 5. Delete one field
	if err := repo.DeleteField(ctx, user.ID, "email"); err != nil {
		log.Fatalf("delete email: %v", err)
	}
	fmt.Println("Email field deleted")

	// 6. Check the remaining Hash in Redis
	fields, err := client.HGetAll(ctx, "user:"+user.ID).Result()
	if err != nil {
		log.Fatalf("read raw hash: %v", err)
	}
	fmt.Printf("Remaining Redis fields: %v\n", fields)

	// 7. Delete complete profile
	if err := repo.Delete(ctx, user.ID); err != nil {
		log.Fatalf("delete user: %v", err)
	}
	fmt.Println("User deleted successfully")
}
