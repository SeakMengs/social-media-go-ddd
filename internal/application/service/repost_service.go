package service

import (
	"context"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
)

type RepostService struct {
	repository repository.RepostRepository
}

func NewRepostService(repo repository.RepostRepository) *RepostService {
	return &RepostService{
		repository: repo,
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

	return repost, nil
}

func (s *RepostService) Delete(ctx context.Context, dl dto.DeleteRepost) error {
	return s.repository.Delete(ctx, dl.ID)
}
func (s *RepostService) GetByID(ctx context.Context, id string) (*entity.Repost, error) {
	return s.repository.FindByID(ctx, id)
}
