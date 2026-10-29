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
	DB           int
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	Writetimeout time.Duration
}

func LoadConfig() (Config, error) {
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))

	if err != nil {
		return Config{}, fmt.Errorf("Invlaid Redis DB: %w", err)
	}

	cfg := Config{
		Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Password: getEnv("REDIS_PASS", ""),
		DB:       db,
	}

	cfg.DialTimeout, err = time.ParseDuration(
		getEnv("REDIS_DIAL_TIMEOUT", "5s"),
	)

	if err != nil || cfg.DialTimeout <= 0 {
		return Config{}, fmt.Errorf("Invlaid Redis Dial Timeout: %w", err)
	}

	cfg.ReadTimeout, err = time.ParseDuration(
		getEnv("REDIS_READ_TIMEOUT", "3s"),
	)

	if err != nil || cfg.ReadTimeout <= 0 {
		return Config{}, fmt.Errorf("Invlaid Redis Read Timeout: %w", err)
	}

	cfg.Writetimeout, err = time.ParseDuration(
		getEnv("REDIS_WRITE_TIMEOUT", "3s"),
	)

	if err != nil || cfg.Writetimeout <= 0 {
		return Config{}, fmt.Errorf("Invlaid Redis Write Timeout: %w", err)
	}

	if cfg.Addr == "" {
		return Config{}, fmt.Errorf("Redis Address Can't Be Empty.")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func newCleint(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DB:           cfg.DB,
		Password:     cfg.Password,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.Writetimeout,
	})
}

func Connect(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := newCleint(cfg)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("Failed to connect Redis: %w", err)
	}

	return client, nil
}
