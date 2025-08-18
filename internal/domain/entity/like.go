package entity

import (
	"errors"
	"social-media-go-ddd/internal/domain/dto"

	"github.com/google/uuid"
)

type Like struct {
	BaseEntity
	UserID uuid.UUID `json:"user_id"`
	PostID uuid.UUID `json:"post_id"`
}

func NewLike(nl dto.NewLike) (*Like, error) {
	like := &Like{
		BaseEntity: NewBaseEntity(),
		UserID:     nl.UserID,
		PostID:     nl.PostID,
	}
	if err := like.Validate(); err != nil {
		return nil, err
	}
	return like, nil
}

func (l *Like) Validate() error {
	if err := l.BaseEntity.Validate(); err != nil {
		return err
	}
	if l.UserID == uuid.Nil {
		return errors.New(ErrLikeUserIDEmpty)
	}
	if l.PostID == uuid.Nil {
		return errors.New(ErrLikePostIDEmpty)
	}
	return nil
}
