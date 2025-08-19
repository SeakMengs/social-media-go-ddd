package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type RepostRepository interface {
	Save(ctx context.Context, r *entity.Repost) error
	// Because a user can do many reposts, delete by id is enough
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.Repost, error)
}
