package entity

import (
	"errors"
	"social-media-go-ddd/internal/domain/dto"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	BaseEntity
	UserID   uuid.UUID `json:"userId"`
	ExpireAt time.Time `json:"expireAt"`
}

func NewSession(ns dto.NewSession) (*Session, error) {
	session := &Session{
		BaseEntity: NewBaseEntity(),
		UserID:     ns.UserID,
		ExpireAt:   ns.ExpireAt,
	}
	if err := session.Validate(); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Session) Validate() error {
	if err := s.BaseEntity.Validate(); err != nil {
		return err
	}
	if s.UserID == uuid.Nil {
		return errors.New(ErrSessionUserIDEmpty)
	}
	if s.ExpireAt.Before(time.Now()) {
		return errors.New(ErrSessionExpiredInPast)
	}
	return nil
}

func (s *Session) IsExpired() bool {
	return s.ExpireAt.Before(time.Now())
}
