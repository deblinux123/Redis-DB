package distributedlock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Lock struct {
	client *redis.Client
	key    string
	token  string
	ttl    time.Duration
}

func NewLock(
	client *redis.Client,
	key string,
	ttl time.Duration,
) (*Lock, error) {
	if client == nil {
		return nil, errors.New("Redis client coannot be nil")
	}

	if key == "" {
		return nil, errors.New("Lock key cannot be empty")
	}

	if ttl <= 0 {
		return nil, errors.New("Lock TTL must be grater than zero")
	}

	token, err := generateToken()

	if err != nil {
		return nil, fmt.Errorf("Generate lock token erro: %w", err)
	}
	return &Lock{
		client: client,
		key:    key,
		token:  token,
		ttl:    ttl,
	}, nil

}

func generateToken() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (l *Lock) Acquire(ctx context.Context) (bool, error) {
	acquired, err := l.client.SetNX(ctx, l.key, l.token, l.ttl).Result()

	if err != nil {
		return false, fmt.Errorf("acquire lock erro: %w", err)
	}

	return acquired, nil
}

func (l *Lock) Release(ctx context.Context) (bool, error) {
	result, err := l.client.Eval(
		ctx,
		releaseScript,
		[]string{l.key},
		l.token,
	).Int()

	if err != nil {
		return false, fmt.Errorf("Release lock: %w", err)
	}

	return result == 1, nil
}

func (l *Lock) Renew(ctx context.Context) (bool, error) {
	result, err := l.client.Eval(
		ctx,
		renewScript,
		[]string{l.key},
		l.token,
		l.ttl,
	).Int()

	if err != nil {
		return false, fmt.Errorf("Renew lock: %w", err)
	}

	return result == 1, nil
}
