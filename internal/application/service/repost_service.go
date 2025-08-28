package service

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/cache"
)

type RepostService struct {
	baseService
	repository repository.RepostRepository
}

func NewRepostService(repo repository.RepostRepository, c cache.Cache) *RepostService {
	return &RepostService{
		baseService: NewBaseService(c),
		repository:  repo,
	}
}

func (s *RepostService) Create(ctx context.Context, nf dto.NewRepost) (*entity.Repost, error) {
	repost, err := entity.NewRepost(nf)
	if err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, repost); err != nil {
		return nil, err
	}

	// Invalidate post cache since repost count changed
	s.cache.Delete(ctx, s.cacheKeys.Post(repost.PostID.String()))
	return repost, nil
}

func (s *RepostService) Delete(ctx context.Context, dl dto.DeleteRepost) error {
	// Invalidate post cache since repost count changed
	s.cache.Delete(ctx, s.cacheKeys.Post(dl.PostID.String()))
	return s.repository.Delete(ctx, dl.UserID.String(), dl.PostID.String())
}

func (s *RepostService) GetByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	reposts, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return reposts, nil
}
