package dto

import "github.com/google/uuid"

type (
	NewLike struct {
		UserID uuid.UUID `json:"user_id"`
		PostID uuid.UUID `json:"post_id"`
	}

	DeleteLike struct {
		UserID uuid.UUID `json:"user_id"`
		PostID uuid.UUID `json:"post_id"`
	}
)
