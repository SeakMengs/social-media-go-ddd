package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type SessionRepository interface {
	Save(ctx context.Context, s *entity.Session) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.Session, error)
	UpdateExpireAt(ctx context.Context, s *entity.Session) error
}
