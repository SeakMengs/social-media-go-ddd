package service

import (
	"context"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
)

type FavoriteService struct {
	repository repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{
		repository: repo,
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

	return favorite, nil
}

func (s *FavoriteService) Delete(ctx context.Context, dl dto.DeleteFavorite) error {
	return s.repository.Delete(ctx, dl.UserID.String(), dl.PostID.String())
}
