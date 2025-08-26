package service

import (
	"context"
	"encoding/json"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/cache"
)

type SessionService struct {
	baseService
	repository repository.SessionRepository
}

func NewSessionService(repo repository.SessionRepository, c cache.Cache) *SessionService {
	return &SessionService{
		repository:  repo,
		baseService: NewBaseService(c),
	}
}

func (s *SessionService) Create(ctx context.Context, ns dto.NewSession) (*entity.Session, error) {
	session, err := entity.NewSession(ns)
	if err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) GetByID(ctx context.Context, id string) (*entity.Session, error) {
	cacheKey := s.cacheKeys.Session(id)
	val, err := s.cache.Get(ctx, cacheKey)
	if !cache.IsCacheError(err) {
		var session entity.Session
		if json.Unmarshal([]byte(val), &session) == nil {
			if !session.IsExpired() {
				return &session, nil
			}
			// Delete expired session from cache
			s.cache.Delete(ctx, cacheKey)
		}
	}

	session, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !session.IsExpired() {
		data, err := json.Marshal(session)
		if err == nil {
			s.cache.Set(ctx, cacheKey, data, cache.DefaultTTL())
		}
	}

	return session, nil
}

func (s *SessionService) Delete(ctx context.Context, ds dto.DeleteSession) error {
	s.cache.Delete(ctx, s.cacheKeys.Session(ds.ID))
	return s.repository.Delete(ctx, ds.ID)
}

func (s *SessionService) UpdateExpireAt(ctx context.Context, old *entity.Session, up dto.UpdateSessionExpireAt) (*entity.Session, error) {
	session, err := entity.NewSessionForUpdate(old, up)
	if err != nil {
		return nil, err
	}
	if err = s.repository.UpdateExpireAt(ctx, session); err != nil {
		return nil, err
	}
	s.cache.Delete(ctx, s.cacheKeys.Session(session.ID.String()))
	return session, nil
}
