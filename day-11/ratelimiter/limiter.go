package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewLimiter(client *redis.Client, limit int, window time.Duration) *Limiter {
	return &Limiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (l *Limiter) Allow(
	ctx context.Context,
	ip string,
) (bool, error) {
	key := "rate_limit:ip:" + ip

	result, err := l.client.Eval(
		ctx,
		RateLimitScript,
		[]string{key},
		int(l.window.Seconds()),
		l.limit,
	).Int()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}
