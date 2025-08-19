package service

import "social-media-go-ddd/internal/domain/repository"

type RepostService struct {
	repository repository.RepostRepository
}

func NewRepostService(repo repository.RepostRepository) *RepostService {
	return &RepostService{
		repository: repo,
	}
}
