package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type LikeRepository interface {
	Save(ctx context.Context, l *entity.Like) error
	// Should delete all like in db by user id and post id because one person should be able to like only one post
	Delete(ctx context.Context, userID, postID string) error
}
