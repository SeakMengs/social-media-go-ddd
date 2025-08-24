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

func NewRedisCache(ctx context.Context, addr, password string, db int, keyPrefix string) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisCache{client: rdb, keyPrefix: keyPrefix}, nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
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

func (r *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	fullPattern := fmt.Sprintf("%s:%s", r.keyPrefix, pattern)
	// cursor is the position to start scanning from, 0 means start from the beginning
	// match is the pattern to match keys against, eg "user:feed:userId:*"
	// count is the number of keys to return per scan iteration in this case 0 redis will default it to 10
	// docs: https://redis.io/docs/latest/commands/scan/
	// in conclusion, this function will delete all keys matching the given pattern
	iter := r.client.Scan(ctx, 0, fullPattern, 0).Iterator()
	for iter.Next(ctx) {
		_ = r.client.Del(ctx, iter.Val()) // ignore errors
	}
	return iter.Err()
}
