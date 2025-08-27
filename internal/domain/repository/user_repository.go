package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, u *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByName(ctx context.Context, name string) (*entity.User, error)
}
