package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/entity"
)

type PostRepository interface {
	Save(ctx context.Context, p *entity.Post) error
	FindByID(ctx context.Context, id string) (*aggregate.Post, error)
	FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error)
	Delete(ctx context.Context, id string, userId string) error
	Update(ctx context.Context, p *entity.Post) error

	// FindFeedByUserID(ctx context.Context, userID string, limit, offset int) ([]*aggregate.Post, error)
}
