package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type LikeRepository interface {
	Save(ctx context.Context, l *entity.Like) error
	Delete(ctx context.Context, userID, postID string) error
	CountByPost(ctx context.Context, postID string) (int, error)
	FindByUserAndPost(ctx context.Context, userID, postID string) (*entity.Like, error)
}
