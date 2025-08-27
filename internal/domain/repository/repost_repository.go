package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/entity"
)

type RepostRepository interface {
	Save(ctx context.Context, r *entity.Repost) error
	// Should delete all reposts in db by user id and post id because one person should be able to repost only one post
	Delete(ctx context.Context, userID string, postID string) error
	FindByID(ctx context.Context, id string) (*entity.Repost, error)
	FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error)
}
