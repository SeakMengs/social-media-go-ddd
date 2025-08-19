package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type FavoriteRepository interface {
	Save(ctx context.Context, f *entity.Favorite) error
	Delete(ctx context.Context, userID, postID string) error
	FindByUser(ctx context.Context, userID string) ([]*entity.Favorite, error)
}
