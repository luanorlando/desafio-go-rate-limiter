package repository

import (
	"context"
	"time"
)

type RateLimiterRepository interface {
	Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)
	Block(ctx context.Context, key string, duration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}
