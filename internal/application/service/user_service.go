package service

import (
	"context"
	"encoding/json"
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

func (s *UserService) GetByID(ctx context.Context, id string, currentUserID *string) (*aggregate.User, error) {
	cacheKey := s.cacheKeys.User(id)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var user aggregate.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	user, err := s.repository.FindByID(ctx, id, currentUserID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
	}

	return user, nil
}

func (s *UserService) GetByName(ctx context.Context, name string, currentUserID *string) (*aggregate.User, error) {
	cacheKey := s.cacheKeys.UserByName(name)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var user aggregate.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	user, err := s.repository.FindByName(ctx, name, currentUserID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
	}

	return user, nil
}

func (s *UserService) GetManyByName(ctx context.Context, name string, currentUserID *string) ([]*aggregate.User, error) {
	users, err := s.repository.SearchManyByName(ctx, name, currentUserID)
	if err != nil {
		return nil, err
	}
	return users, nil
}
