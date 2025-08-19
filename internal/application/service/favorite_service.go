package service

import "social-media-go-ddd/internal/domain/repository"

type FavoriteService struct {
	repository repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{
		repository: repo,
	}
}
