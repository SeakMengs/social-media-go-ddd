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
	return s.repository.FindByID(ctx, id)
}

func (s *PostService) GetByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	posts, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) Delete(ctx context.Context, dp dto.DeletePost) error {
	return s.repository.Delete(ctx, dp.ID, dp.UserID.String())
}

func (s *PostService) Update(ctx context.Context, old *entity.Post, up dto.UpdatePost) (*entity.Post, error) {
	post, err := entity.NewPostForUpdate(old, up)
	if err != nil {
		return nil, err
	}
	err = s.repository.Update(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}
