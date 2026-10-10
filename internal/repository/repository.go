package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RateLimiterRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := r.client.Pipeline()

	incrCmd := pipe.Incr(ctx, key)

	pipe.Expire(ctx, key, expiration)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incrCmd.Val(), nil
}

func (r *RedisRepository) Block(ctx context.Context, key string, duration time.Duration) error {
	blockKey := "blocked:" + key
	return r.client.Set(ctx, blockKey, "1", duration).Err()
}

func (r *RedisRepository) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := "blocked:" + key
	val, err := r.client.Get(ctx, blockKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
