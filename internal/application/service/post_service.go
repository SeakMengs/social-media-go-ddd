package service

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
)

type PostService struct {
	repository repository.PostRepository
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{
		repository: repo,
	}
}

func (s *PostService) Create(ctx context.Context, np dto.NewPost) (*entity.Post, error) {
	post, err := entity.NewPost(np)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Save(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetByID(ctx context.Context, id string) (*aggregate.Post, error) {
	post, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	posts, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
