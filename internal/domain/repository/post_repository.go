package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/entity"
)

type PostRepository interface {
	Save(ctx context.Context, p *entity.Post) error
	FindByID(ctx context.Context, id string, currentUserID *string) (*aggregate.Post, error)
	FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error)
	Delete(ctx context.Context, id string, userID string) error
	Update(ctx context.Context, p *entity.Post) error
	// Get followed users or own posts, reposts sort by created_at desc
	FindFeed(ctx context.Context, userID string, limit, offset int) ([]*aggregate.Post, int, error)
}
