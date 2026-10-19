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
		return Config{}, fmt.Errorf("Invalid DB.")
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
		return Config{}, fmt.Errorf("Invalid REDIS_DIAL_TIMEOUT")
	}

	cfg.ReadTimeout, durationErr = time.ParseDuration(
		getEnv("REDIS_READ_TIMEOUT", "3S"),
	)

	if durationErr != nil || cfg.ReadTimeout <= 0 {
		return Config{}, fmt.Errorf("Invalid REDIS_READ_TIMEOUT")
	}

	cfg.WriteTimeout, durationErr = time.ParseDuration(
		getEnv("REDIS_WRITE_TIMEOUT", "3S"),
	)

	if durationErr != nil || cfg.WriteTimeout <= 0 {
		return Config{}, fmt.Errorf("Invalid READIS_WRITE_TIMEOUT")
	}

	if cfg.Addr == "" {
		return Config{}, fmt.Errorf("Redis address cannot be empty.")
	}

	return cfg, nil
}

func getEnv(key, fallbak string) string {
	if value := os.Getenv(key); key != "" {
		return value
	}

	return fallbak
}

func NewClient(cfg Config) *redis.Client {
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
	client := NewClient(cfg)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()

		return nil, fmt.Errorf("Redis Connection failed: %w", err)
	}

	return client, nil
}
