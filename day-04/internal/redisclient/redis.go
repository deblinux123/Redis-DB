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
		return Config{}, fmt.Errorf("invalid REDIS_DB: must be a non-negative integer")
	}

	cfg := Config{
		Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}

	cfg.DialTimeout, err = time.ParseDuration(
		getEnv("REDIS_DIAL_TIMEOUT", "5s"),
	)
	if err != nil || cfg.DialTimeout <= 0 {
		return Config{}, fmt.Errorf("invalid REDIS_DIAL_TIMEOUT")
	}

	cfg.ReadTimeout, err = time.ParseDuration(
		getEnv("REDIS_READ_TIMEOUT", "3s"),
	)
	if err != nil || cfg.ReadTimeout <= 0 {
		return Config{}, fmt.Errorf("invalid REDIS_READ_TIMEOUT")
	}

	cfg.WriteTimeout, err = time.ParseDuration(
		getEnv("REDIS_WRITE_TIMEOUT", "3s"),
	)
	if err != nil || cfg.WriteTimeout <= 0 {
		return Config{}, fmt.Errorf("invalid REDIS_WRITE_TIMEOUT")
	}

	if cfg.Addr == "" {
		return Config{}, fmt.Errorf("Redis address cannot be empty")
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
