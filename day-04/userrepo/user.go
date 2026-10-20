package userrepo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	City  string `json:"city"`
}

type Repository struct {
	client *redis.Client
}

func New(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func userKey(id string) string {
	return "user:" + id
}

func (r *Repository) Create(ctx context.Context, user User) error {
	if user.ID == "" {
		return errors.New("User ID cannot be empty")
	}

	if user.Age < 0 {
		return errors.New("Age cannot be negative")
	}

	key := userKey(user.ID)

	created, err := r.client.HSetNX(
		ctx,
		key,
		"_initialized",
		"1",
	).Result()

	if err != nil {
		return fmt.Errorf("Initialize user: %w", err)
	}

	if !created {
		return fmt.Errorf("User %s Already exists", user.ID)
	}

	fields := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
		"city":  user.City,
	}

	if err := r.client.HSet(ctx, key, fields).Err(); err != nil {
		_ = r.client.Del(ctx, key).Err()

		return fmt.Errorf("create User: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (User, error) {
	key := userKey(id)

	fields, err := r.client.HGetAll(
		ctx,
		key,
	).Result()

	if err != nil {
		return User{}, fmt.Errorf("Get User: %w", err)
	}

	if len(fields) == 0 {
		return User{}, ErrUserNotFound
	}

	age, err := strconv.Atoi(fields["age"])

	if err != nil {
		return User{}, fmt.Errorf("Invalid age stored")
	}

	return User{
		ID:    fields["id"],
		Name:  fields["name"],
		Email: fields["email"],
		Age:   age,
		City:  fields["city"],
	}, nil
}

func (r *Repository) UpdateField(ctx context.Context, id string, field string, value interface{}) error {
	key := userKey(id)

	exists, err := r.client.Exists(ctx, key).Result()

	if err != nil {
		return fmt.Errorf("Check user: %w", err)
	}

	if exists == 0 {
		return ErrUserNotFound
	}

	switch field {
	case "name", "email", "city":
	case "age":
		switch v := value.(type) {
		case int:
			if v < 0 {
				return errors.New("age cannot be negative")
			}
		case int64:
			if v < 0 {
				return errors.New("age cannot be negative")
			}
		default:
			return errors.New("age must be an integer.")
		}
	default:
		return fmt.Errorf("unsupported filed: %s", field)
	}

	return r.client.HSet(ctx, key, field, value).Err()
}

func (r *Repository) DeleteField(ctx context.Context, id string, field string) error {
	key := userKey(id)

	exists, err := r.client.Exists(ctx, key).Result()

	if err != nil {
		return fmt.Errorf("Check user: %w", err)
	}

	if exists == 0 {
		return ErrUserNotFound
	}

	removed, err := r.client.HDel(ctx, key, field).Result()

	if err != nil {
		return fmt.Errorf("Delete field: %w", err)
	}

	if removed == 0 {
		return fmt.Errorf("Filed %q not found", field)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	removed, err := r.client.Del(ctx, userKey(id)).Result()

	if err != nil {
		return fmt.Errorf("Delete user: %w", err)
	}

	if removed == 0 {
		return ErrUserNotFound
	}

	return nil
}
