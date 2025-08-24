package service

import (
	"context"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/cache"
)

type LikeService struct {
	baseService
	repository repository.LikeRepository
}

func NewLikeService(repo repository.LikeRepository, c cache.Cache) *LikeService {
	return &LikeService{
		baseService: NewBaseService(c),
		repository:  repo,
	}
}

func (s *LikeService) Create(ctx context.Context, nl dto.NewLike) (*entity.Like, error) {
	like, err := entity.NewLike(nl)
	if err != nil {
		return nil, err
	}

	if err = s.repository.Save(ctx, like); err != nil {
		return nil, err
	}

	return like, nil
}

func (s *LikeService) Delete(ctx context.Context, dl dto.DeleteLike) error {
	return s.repository.Delete(ctx, dl.UserID.String(), dl.PostID.String())
}
