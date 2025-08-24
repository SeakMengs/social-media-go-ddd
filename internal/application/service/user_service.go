package service

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/cache"
)

type UserService struct {
	baseService
	repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository, c cache.Cache) *UserService {
	return &UserService{
		baseService: NewBaseService(c),
		repository:  repo,
	}
}

func (s *UserService) Create(ctx context.Context, nu dto.NewUser) (*entity.User, error) {
	user, err := entity.NewUser(nu)
	if err != nil {
		return nil, err
	}

	if err = s.repository.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *UserService) GetByName(ctx context.Context, name string) (*entity.User, error) {
	return s.repository.FindByName(ctx, name)
}

// return posts, total, error
func (s *UserService) GetFeed(ctx context.Context, userID string, limit, offset int) ([]*aggregate.Post, int, error) {
	return s.repository.FindFeed(ctx, userID, limit, offset)
}
