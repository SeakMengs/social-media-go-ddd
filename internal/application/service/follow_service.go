package service

import (
	"context"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
)

type FollowService struct {
	repository repository.FollowRepository
}

func NewFollowService(repo repository.FollowRepository) *FollowService {
	return &FollowService{
		repository: repo,
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

	return follow, nil
}

func (s *FollowService) Delete(ctx context.Context, dl dto.DeleteFollow) error {
	return s.repository.Delete(ctx, dl.FollowerID.String(), dl.FolloweeID.String())
}
