package cache

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache: key not found")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

func IsCacheMiss(err error) bool {
	return errors.Is(err, ErrCacheMiss)
}

func IsCacheError(err error) bool {
	return err != nil && !IsCacheMiss(err)
}

func DefaultTTL() time.Duration {
	return 5 * time.Minute
}
