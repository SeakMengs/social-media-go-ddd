package service

import (
	"context"
	"encoding/json"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/cache"
)

type FavoriteService struct {
	baseService
	repository repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository, c cache.Cache) *FavoriteService {
	return &FavoriteService{
		baseService: NewBaseService(c),
		repository:  repo,
	}
}

func (s *FavoriteService) Create(ctx context.Context, nf dto.NewFavorite) (*entity.Favorite, error) {
	favorite, err := entity.NewFavorite(nf)
	if err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, favorite); err != nil {
		return nil, err
	}

	// Invalidate post cache since favorite count changed
	s.cache.Delete(ctx, s.cacheKeys.Post(favorite.PostID.String()))
	return favorite, nil
}

func (s *FavoriteService) Delete(ctx context.Context, dl dto.DeleteFavorite) error {
	// Invalidate post cache since favorite count changed
	s.cache.Delete(ctx, s.cacheKeys.Post(dl.PostID.String()))
	return s.repository.Delete(ctx, dl.UserID.String(), dl.PostID.String())
}

func (s *FavoriteService) GetByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	cacheKey := s.cacheKeys.UserFavoritePosts(userID)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var posts []*aggregate.Post
		if json.Unmarshal([]byte(val), &posts) == nil {
			return posts, nil
		}
	}

	post, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(post)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
	}

	return post, nil
}
