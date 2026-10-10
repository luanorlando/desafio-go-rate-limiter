package limiter

import (
	"context"
	"fmt"
	"time"

	"rate-limiter/internal/repository"
)

type RateLimiter struct {
	repo              repository.RateLimiterRepository
	defaultLimitIP    int
	defaultLimitToken int
	blockDuration     time.Duration
}

func NewRateLimiter(repo repository.RateLimiterRepository, defaultLimitIP int, defaultLimitToken int, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		repo:              repo,
		defaultLimitIP:    defaultLimitIP,
		defaultLimitToken: defaultLimitToken,
		blockDuration:     blockDuration,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, ip string, token string) (bool, error) {
	var identifier string
	var limit int

	if token != "" {
		identifier = fmt.Sprintf("token:%s", token)
		limit = rl.defaultLimitToken
	} else {
		identifier = fmt.Sprintf("ip:%s", ip)
		limit = rl.defaultLimitIP
	}

	blocked, err := rl.repo.IsBlocked(ctx, identifier)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	currentSecond := time.Now().Unix()
	key := fmt.Sprintf("counter:%s:%d", identifier, currentSecond)

	count, err := rl.repo.Increment(ctx, key, 2*time.Second)
	if err != nil {
		return false, err
	}

	if int(count) > limit {
		err := rl.repo.Block(ctx, identifier, rl.blockDuration)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
