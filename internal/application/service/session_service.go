package service

import (
	"context"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/repository"
)

type SessionService struct {
	repository repository.SessionRepository
}

func NewSessionService(repo repository.SessionRepository) *SessionService {
	return &SessionService{
		repository: repo,
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
	session, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
