package redis

import (
	"context"
	"fmt"
	"social-media-go-ddd/internal/infrastructure/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	// key prefix for database driver for example: "mysql:user:id", "postgres:user:id"
	keyPrefix string
}

func NewRedisCache(addr, password string, db int, keyPrefix string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: rdb, keyPrefix: keyPrefix}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", r.keyPrefix, key)).Result()
	if err == redis.Nil {
		return "", cache.ErrCacheMiss
	}
	return val, err
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, fmt.Sprintf("%s:%s", r.keyPrefix, key), value, expiration).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, fmt.Sprintf("%s:%s", r.keyPrefix, key)).Err()
}
