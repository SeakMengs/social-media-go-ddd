package dto

import (
	"github.com/google/uuid"
)

type (
	NewFollow struct {
		FollowerID uuid.UUID `json:"follower_id"`
		FolloweeID uuid.UUID `json:"followee_id"`
	}

	DeleteFollow struct {
		FollowerID uuid.UUID `json:"follower_id"`
		FolloweeID uuid.UUID `json:"followee_id"`
	}
)
