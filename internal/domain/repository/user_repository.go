package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, u *entity.User) error
	FindByID(ctx context.Context, id string, currentUserID string) (*aggregate.User, error)
	FindByName(ctx context.Context, name string, currentUserID string) (*aggregate.User, error)
}
