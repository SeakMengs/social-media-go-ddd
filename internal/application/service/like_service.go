package service

import "social-media-go-ddd/internal/domain/repository"

type LikeService struct {
	repository repository.LikeRepository
}

func NewLikeService(repo repository.LikeRepository) *LikeService {
	return &LikeService{
		repository: repo,
	}
}
