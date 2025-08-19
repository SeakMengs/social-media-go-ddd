package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type RepostRepository interface {
	Save(ctx context.Context, r *entity.Repost) error
	Delete(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userID string) ([]*entity.Repost, error)
}
