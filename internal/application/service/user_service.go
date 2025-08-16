package service

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (s *UserService) Create(ctx context.Context, name string) (*entity.User, error) {
	user, err := entity.NewUser(name)
	if err != nil {
		return nil, err
	}

	if err = s.repository.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
