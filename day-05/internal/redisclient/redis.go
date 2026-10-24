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
		return Config{}, fmt.Errorf("Invalid Redis DB")
	}

	cfg := Config{
		Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Password: os.Getenv("REDIDS_PASSWORD"),
		DB:       db,
	}

	cfg.DialTimeout, err = time.ParseDuration(
		getEnv("REDIS_DIAL_TIMEOUT", "5S"),
	)

	if err != nil || cfg.DialTimeout <= 0 {
		return Config{}, fmt.Errorf("Invalid Redis Dial Timeout")
	}

	cfg.ReadTimeout, err = time.ParseDuration(
		getEnv("REDIS_READ_TIMEOUT", "3S"),
	)

	if err != nil || cfg.ReadTimeout <= 0 {
		return Config{}, fmt.Errorf("Invalid Redis Read Timeout")
	}

	cfg.WriteTimeout, err = time.ParseDuration(
		getEnv("REDIS_WRITE_TIMEOUT", "3S"),
	)

	if err != nil || cfg.WriteTimeout <= 0 {
		return Config{}, fmt.Errorf("Invalid Redis Write Timeout")
	}

	if cfg.Addr == "" {
		return Config{}, fmt.Errorf("Redis Address connot be empty")
	}

	return cfg, nil

}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
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

		return nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	return client, nil
}
