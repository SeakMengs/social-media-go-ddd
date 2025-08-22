package repository

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"
)

type FollowRepository interface {
	Save(ctx context.Context, f *entity.Follow) error
	// Should delete all follow relationships in db by follower id and followee id because one person should be able to follow only one person
	Delete(ctx context.Context, followerID, followeeID string) error
}
