package entity

import (
	"social-media-go-ddd/internal/domain/dto"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	BaseEntity
	UserID   uuid.UUID `json:"userId"`
	ExpireAt time.Time `json:"expireAt"`
}

func DefaultSessionExpireAt() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
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

func NewSessionForUpdate(oldSession *Session, up dto.UpdateSessionExpireAt) (*Session, error) {
	session := &Session{
		BaseEntity: oldSession.BaseEntity,
		UserID:     oldSession.UserID,
		ExpireAt:   up.ExpireAt,
	}
	session.UpdateTimestamp()
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
		return ErrSessionUserIDEmpty
	}
	if s.IsExpired() {
		return ErrSessionExpired
	}
	return nil
}

func (s *Session) IsExpired() bool {
	return s.ExpireAt.Before(time.Now())
}

func (s *Session) UpdateExpireAt(expireAt time.Time) error {
	s.ExpireAt = expireAt
	s.UpdateTimestamp()
	return s.Validate()
}
