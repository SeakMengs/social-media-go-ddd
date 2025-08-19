package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type FavoriteRepository interface {
	Save(ctx context.Context, f *entity.Favorite) error
	// Should delete all favorite in db by user id and post id because one person should be able to favorite only one post
	Delete(ctx context.Context, userID, postID string) error
}
