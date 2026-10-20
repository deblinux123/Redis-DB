package redisclient

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadConfig() (Config, error) {
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))

	if err != nil || db < 0 {
		return Config{}, fmt.Errorf("Invalid DB")
	}

	cfg := Config{
		Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}

	var durationErr error

	cfg.DialTimeout, durationErr = time.ParseDuration(
		getEnv("REDIS_DIAL_TIMEOUT", "5S"),
	)

	if durationErr != nil || cfg.DialTimeout <= 0 {
		return Config{}, fmt.Errorf("Invalid Redis Dial Timeout")
	}

	cfg.ReadTimeout, durationErr = time.ParseDuration(
		getEnv("REDIS_READ_TIMEOUT", "3S"),
	)

	if durationErr != nil || cfg.ReadTimeout <= 0 {
		return Config{}, fmt.Errorf("Invalid Redis Read Timeout")
	}

	cfg.WriteTimeout, durationErr = time.ParseDuration(
		getEnv("REDIS_WRITE_TIMEOUT", "3S"),
	)

	if durationErr != nil || cfg.WriteTimeout <= 0 {
		return Config{}, fmt.Errorf("Invlaid Redis Write Timeout")
	}

	if cfg.Addr == "" {
		return Config{}, fmt.Errorf("Redis Address cannot be empty.")
	}

	return Config{}, nil
}

func getEnv(key, fallbakc string) string {
	if value := os.Getenv(key); key != "" {
		return value
	}

	return fallbakc
}

func newClient(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}

func Connect(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := newClient(cfg)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()

		return nil, fmt.Errorf("Failed to connect redis:%w", err)
	}

	return client, nil
}
