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
)
