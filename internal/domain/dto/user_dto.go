package dto

import (
	"time"

	"github.com/google/uuid"
)

type (
	NewUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	UserLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	CommonUserAggregate struct {
		Followed       bool `json:"followed"`
		FollowerCount  int  `json:"followerCount"`
		FollowingCount int  `json:"followingCount"`
	}

	UserResponse struct {
		ID        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt,omitempty"`
		UpdatedAt time.Time `json:"updatedAt,omitempty"`
	}

	UserAggregateResponse struct {
		ID             uuid.UUID `json:"id"`
		Username       string    `json:"username"`
		Email          string    `json:"email"`
		CreatedAt      time.Time `json:"createdAt,omitempty"`
		UpdatedAt      time.Time `json:"updatedAt,omitempty"`
		Followed       bool      `json:"followed"`
		FollowerCount  int       `json:"followerCount"`
		FollowingCount int       `json:"followingCount"`
	}
)
