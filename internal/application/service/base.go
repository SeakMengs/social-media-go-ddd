package service

import "social-media-go-ddd/internal/infrastructure/cache"

type baseService struct {
	cache     cache.Cache
	cacheKeys cacheKeys
}

func NewBaseService(c cache.Cache) baseService {
	return baseService{
		cache:     c,
		cacheKeys: NewCacheKeys(),
	}
}
