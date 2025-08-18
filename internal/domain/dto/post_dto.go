package dto

import "github.com/google/uuid"

type (
	NewPost struct {
		UserID  uuid.UUID `json:"user_id"`
		Content string    `json:"content"`
	}
)
