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

type RepostService struct {
	baseService
	repository repository.RepostRepository
}

func NewRepostService(repo repository.RepostRepository, c cache.Cache) *RepostService {
	return &RepostService{
		baseService: NewBaseService(c),
		repository:  repo,
	}
}

func (s *RepostService) Create(ctx context.Context, nf dto.NewRepost) (*entity.Repost, error) {
	repost, err := entity.NewRepost(nf)
	if err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, repost); err != nil {
		return nil, err
	}

	// Invalidate post cache since repost count changed
	s.cache.Delete(ctx, s.cacheKeys.Post(repost.PostID.String()))
	// Invalidate user reposts cache
	s.cache.Delete(ctx, s.cacheKeys.UserPosts(repost.UserID.String()))
	s.cache.Delete(ctx, s.cacheKeys.UserReposts(repost.UserID.String()))
	s.cache.DeleteByPattern(ctx, s.cacheKeys.UserFeedPattern(repost.UserID.String()))

	return repost, nil
}

func (s *RepostService) Delete(ctx context.Context, dl dto.DeleteRepost) error {
	// Invalidate post cache since repost count changed
	s.cache.Delete(ctx, s.cacheKeys.Post(dl.PostID.String()))
	// Invalidate user reposts cache
	s.cache.Delete(ctx, s.cacheKeys.UserPosts(dl.UserID.String()))
	s.cache.Delete(ctx, s.cacheKeys.UserReposts(dl.UserID.String()))
	s.cache.DeleteByPattern(ctx, s.cacheKeys.UserFeedPattern(dl.UserID.String()))
	return s.repository.Delete(ctx, dl.UserID.String(), dl.PostID.String())
}

func (s *RepostService) GetByID(ctx context.Context, id string) (*entity.Repost, error) {
	cacheKey := s.cacheKeys.Repost(id)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var repost entity.Repost
		if json.Unmarshal([]byte(val), &repost) == nil {
			return &repost, nil
		}
	}

	repost, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(repost)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
	}

	return repost, nil
}

func (s *RepostService) GetByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	cacheKey := s.cacheKeys.UserReposts(userID)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var posts []*aggregate.Post
		if json.Unmarshal([]byte(val), &posts) == nil {
			return posts, nil
		}
	}

	reposts, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(reposts)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
	}

	return reposts, nil
}
