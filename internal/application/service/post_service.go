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

type PostService struct {
	baseService
	repository repository.PostRepository
}

func NewPostService(repo repository.PostRepository, c cache.Cache) *PostService {
	return &PostService{
		baseService: NewBaseService(c),
		repository:  repo,
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

func (s *PostService) GetByID(ctx context.Context, id string, currentUserID *string) (*aggregate.Post, error) {
	cacheKey := s.cacheKeys.Post(id)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var post aggregate.Post
		if json.Unmarshal([]byte(val), &post) == nil {
			return &post, nil
		}
	}

	post, err := s.repository.FindByID(ctx, id, currentUserID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(post)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
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

func (s *PostService) Delete(ctx context.Context, dp dto.DeletePost) error {
	// Only invalidate the specific post cache
	s.cache.Delete(ctx, s.cacheKeys.Post(dp.ID))
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

	// Only invalidate the specific post cache
	s.cache.Delete(ctx, s.cacheKeys.Post(post.ID.String()))
	return post, nil
}

// return posts, total, error
func (s *PostService) GetFeed(ctx context.Context, userID string, limit, offset int) ([]*aggregate.Post, int, error) {
	posts, total, err := s.repository.FindFeed(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
