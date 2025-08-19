package dto

import "github.com/google/uuid"

type (
	NewFavorite struct {
		UserID uuid.UUID `json:"user_id"`
		PostID uuid.UUID `json:"post_id"`
	}

	DeleteFavorite struct {
		UserID uuid.UUID `json:"user_id"`
		PostID uuid.UUID `json:"post_id"`
	}
)
