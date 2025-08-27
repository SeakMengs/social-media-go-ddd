package service

import (
	"context"
	"encoding/json"
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
	cacheKey := s.cacheKeys.User(id)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var user entity.User
		if json.Unmarshal([]byte(val), &user) == nil {
			return &user, nil
		}
	}

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
	}

	return user, nil
}

func (s *UserService) GetByName(ctx context.Context, name string) (*entity.User, error) {
	cacheKey := s.cacheKeys.UserByName(name)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		user, err := entity.UserUnmarshalJson([]byte(val))
		if err == nil {
			return user, nil
		}
	}

	user, err := s.repository.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	data, err := user.MarshalJson()
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
	}

	return user, nil
}
