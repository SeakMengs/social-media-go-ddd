package entity

import (
	"social-media-go-ddd/internal/domain/dto"

	"github.com/google/uuid"
)

type Like struct {
	BaseEntity
	UserID uuid.UUID `json:"userId"`
	PostID uuid.UUID `json:"postId"`
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
		return ErrLikeUserIDEmpty
	}
	if l.PostID == uuid.Nil {
		return ErrLikePostIDEmpty
	}
	return nil
}
