package entity

import (
	"social-media-go-ddd/internal/domain/dto"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repost struct {
	BaseEntity
	UserID  uuid.UUID `json:"userId"`
	PostID  uuid.UUID `json:"postId"`
	Comment string    `json:"comment"`
}

func NewRepost(nr dto.NewRepost) (*Repost, error) {
	repost := &Repost{
		BaseEntity: NewBaseEntity(),
		UserID:     nr.UserID,
		PostID:     nr.PostID,
		Comment:    nr.Comment,
	}
	if err := repost.Validate(); err != nil {
		return nil, err
	}
	return repost, nil
}

func (r *Repost) Validate() error {
	if err := r.BaseEntity.Validate(); err != nil {
		return err
	}
	if r.UserID == uuid.Nil {
		return ErrRepostUserIDEmpty
	}
	if r.PostID == uuid.Nil {
		return ErrRepostPostIDEmpty
	}
	if len(r.Comment) > 1000 {
		return ErrRepostCommentTooLong
	}
	return nil
}

func (r *Repost) UpdateComment(comment string) error {
	r.Comment = strings.TrimSpace(comment)
	r.UpdatedAt = time.Now()
	return r.Validate()
}
