package service

import (
	"context"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/cache"
)

type FollowService struct {
	baseService
	repository repository.FollowRepository
}

func NewFollowService(repo repository.FollowRepository, c cache.Cache) *FollowService {
	return &FollowService{
		baseService: NewBaseService(c),
		repository:  repo,
	}
}

func (s *FollowService) Create(ctx context.Context, nf dto.NewFollow) (*entity.Follow, error) {
	follow, err := entity.NewFollow(nf)
	if err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, follow); err != nil {
		return nil, err
	}

	// Since we follow someone, our feed might change, so we invalidate the cache
	s.cache.DeleteByPattern(ctx, s.cacheKeys.UserFeedPattern(follow.FollowerID.String()))
	return follow, nil
}

func (s *FollowService) Delete(ctx context.Context, dl dto.DeleteFollow) error {
	s.cache.DeleteByPattern(ctx, s.cacheKeys.UserFeedPattern(dl.FollowerID.String()))
	return s.repository.Delete(ctx, dl.FollowerID.String(), dl.FolloweeID.String())
}
