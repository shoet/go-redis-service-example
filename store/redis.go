package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shoet/go-redis-service-example/config"
)

type RedisClient struct {
	redis     *redis.Client
	expireSec int
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	res := r.redis.Get(ctx, key)
	if err := res.Err(); err != nil {
		return "", fmt.Errorf("failed to get: %w", err)
	}
	return res.Val(), nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	res := r.redis.Set(ctx, key, value, time.Second*time.Duration(r.expireSec))
	if res.Err() != nil {
		return fmt.Errorf("failed to set: %w", res.Err())
	}
	return nil
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.KVSHost, cfg.KVSPort),
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	c := &RedisClient{
		redis:     r,
		expireSec: cfg.KVSExpireSec,
	}
	return c
}

type KVStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}
