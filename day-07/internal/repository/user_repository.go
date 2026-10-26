package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deblinux123/Redis-DB/day-07/internal/model"
	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	Client *redis.Client
}

func NewUserRepository(client *redis.Client) *UserRepository {
	return &UserRepository{
		Client: client,
	}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	key := "user:" + user.ID

	return r.Client.HSet(ctx, key, map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
		"age":  user.Age,
	}).Err()
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (model.User, error) {
	key := "user:" + id

	data, err := r.Client.HGetAll(ctx, key).Result()

	if err != nil {
		return model.User{}, fmt.Errorf("Get user From Redis Failed: %v", err)
	}

	if len(data) == 0 {
		return model.User{}, redis.Nil
	}

	age, err := strconv.Atoi(data["age"])

	if err != nil {
		return model.User{}, fmt.Errorf("Invlaid age stored in Redis: %v", err)
	}

	return model.User{
		ID:   data["id"],
		Name: data["name"],
		Age:  age,
	}, nil
}
