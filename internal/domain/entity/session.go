package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	BaseEntity
	UserID   uuid.UUID `json:"userId"`
	ExpireAt time.Time `json:"expireAt"`
}

func NewSession(userID uuid.UUID, expireAt time.Time) (*Session, error) {
	session := &Session{
		BaseEntity: *NewBaseEntity(),
		UserID:     userID,
		ExpireAt:   expireAt,
	}

	if err := session.Validate(); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Session) Validate() error {
	return nil
}

func (s *Session) IsExpired() bool {
	return s.ExpireAt.Before(time.Now())
}
