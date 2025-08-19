package dto

import (
	"github.com/google/uuid"
)

type (
	NewRepost struct {
		UserID  uuid.UUID `json:"user_id"`
		PostID  uuid.UUID `json:"post_id"`
		Comment string    `json:"comment"`
	}

	DeleteRepost struct {
		ID string `json:"id"`
	}
)
