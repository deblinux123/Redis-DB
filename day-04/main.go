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
		log.Fatalf("Load redis config: %v", err)
	}

	client, err := redisclient.Connect(ctx, cfg)

	if err != nil {
		log.Fatalf("Connect ro redis: %v", err)
	}

	defer client.Close()

	repo := userrepo.New(client)

	user := userrepo.User{
		ID:    "1003",
		Name:  "Asma",
		Email: "Asma@gmail.com",
		Age:   29,
		City:  "Karaj",
	}

	if err := repo.Delete(ctx, user.ID); err != nil &&
		!errors.Is(err, userrepo.ErrUserNotFound) {
		log.Fatalf("Clean previous user: %v", err)
	}

	// create user
	if err := repo.Create(ctx, user); err != nil {
		log.Fatalf("Create user: %v", err)
	}

	fmt.Println("User created successfully")

	// get profile
	got, err := repo.Get(ctx, user.ID)

	if err != nil {
		log.Fatalf("Ge user: %v", err)
	}

	fmt.Printf("User: %+v\n", got)

	// update one field
	if err := repo.UpdateField(
		ctx,
		user.ID,
		"city",
		"KARAJ",
	); err != nil {
		log.Fatalf("Update city: %v", err)
	}

	fmt.Println("City updated successfully.")

	// read the updated profile
	got, err = repo.Get(ctx, user.ID)

	if err != nil {
		log.Fatalf("Get update user: %v", err)
	}

	fmt.Printf("Updated user: %v\n", got)

	// delete one filed
	if err := repo.DeleteField(ctx, user.ID, "email"); err != nil {
		log.Fatalf("Delete email: %v", err)
	}

	fmt.Println("Email field deleted")

	// check the remining Hash in Redis
	fields, err := client.HGetAll(ctx, "user"+user.ID).Result()

	if err != nil {
		log.Fatalf("Read raw Hash: %v", err)
	}

	fmt.Printf("Remaining Redis fields: %v\n", fields)

	// delete complite user profile

	if err := repo.Delete(ctx, user.ID); err != nil {
		log.Fatalf("Delete user: %v", err)
	}

	fmt.Println("User deleted successfully.")

	user2 := userrepo.User{
		ID:    "1005",
		Name:  "Ali",
		Email: "ali123@gmail.com",
		Age:   33,
		City:  "Tehran",
	}

	if err = repo.Delete(ctx, user2.ID); err != nil && !errors.Is(err, userrepo.ErrUserNotFound) {
		log.Fatalf("Clean previous user: %v", err)
	}

	if err := repo.Create(ctx, user2); err != nil {
		log.Fatalf("User 2 create: %v", err)
	}

	fmt.Println("User 2 created.")

	// update the user 2 email
	err = repo.UpdateField(ctx, user2.ID, "email", "ali-123321@email.com")

	if err != nil {
		log.Fatalf("Update email user 2: %v", err)
	}

	fmt.Println("User 2 email updated.")

	count, err := client.HLen(ctx, "user:"+user2.ID).Result()

	if err != nil {
		log.Fatalf("Count user 2 fileds: %v", err)
	}

	fmt.Printf("User 2 fileds: %d\n", count)

	fmt.Println(client.HExists(ctx, "user:"+user2.ID, "email"))

}
