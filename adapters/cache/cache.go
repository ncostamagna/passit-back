package cache

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type (
	Cache interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value string, expiration time.Duration) error
		Delete(ctx context.Context, key string) error
	}

	cache struct {
		redis *redis.Client
	}
)

func NewCache(addr string, password string, db int) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &cache{redis: rdb}
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	return c.redis.Get(ctx, key).Result()
}

func (c *cache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.redis.Set(ctx, key, value, expiration).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}
