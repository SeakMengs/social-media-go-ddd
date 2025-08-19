package entity

import (
	"errors"
	"social-media-go-ddd/internal/domain/dto"

	"github.com/google/uuid"
)

type Favorite struct {
	BaseEntity
	UserID uuid.UUID `json:"userId"`
	PostID uuid.UUID `json:"postId"`
}

func NewFavorite(nf dto.NewFavorite) (*Favorite, error) {
	favorite := &Favorite{
		BaseEntity: NewBaseEntity(),
		UserID:     nf.UserID,
		PostID:     nf.PostID,
	}
	if err := favorite.Validate(); err != nil {
		return nil, err
	}
	return favorite, nil
}

func (f *Favorite) Validate() error {
	if err := f.BaseEntity.Validate(); err != nil {
		return err
	}
	if f.UserID == uuid.Nil {
		return errors.New(ErrFavoriteUserIDEmpty)
	}
	if f.PostID == uuid.Nil {
		return errors.New(ErrFavoritePostIDEmpty)
	}
	return nil
}
