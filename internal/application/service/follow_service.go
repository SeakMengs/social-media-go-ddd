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

	// Update the cache of both users involved in the follow relationship
	// since their follower/following counts changed
	s.cache.Delete(ctx, s.cacheKeys.User(follow.FollowerID.String()))
	s.cache.Delete(ctx, s.cacheKeys.User(follow.FolloweeID.String()))
	return follow, nil
}

func (s *FollowService) Delete(ctx context.Context, dl dto.DeleteFollow) error {
	// Update the cache of both users involved in the follow relationship
	// since their follower/following counts changed
	s.cache.Delete(ctx, s.cacheKeys.User(dl.FolloweeID.String()))
	s.cache.Delete(ctx, s.cacheKeys.User(dl.FollowerID.String()))
	return s.repository.Delete(ctx, dl.FollowerID.String(), dl.FolloweeID.String())
}
