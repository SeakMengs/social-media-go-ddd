package dto

import (
	"time"

	"github.com/google/uuid"
)

type (
	NewSession struct {
		UserID   uuid.UUID `json:"user_id"`
		ExpireAt time.Time `json:"expire_at"`
	}

	UpdateSessionExpireAt struct {
		ID       string    `json:"id"`
		ExpireAt time.Time `json:"expire_at"`
	}

	DeleteSession struct {
		ID string `json:"id"`
	}
)
