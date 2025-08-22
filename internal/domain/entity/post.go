package entity

import (
	"social-media-go-ddd/internal/domain/dto"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	BaseEntity
	UserID  uuid.UUID `json:"userId"`
	Content string    `json:"content"`
}

func NewPost(np dto.NewPost) (*Post, error) {
	post := &Post{
		BaseEntity: NewBaseEntity(),
		UserID:     np.UserID,
		Content:    strings.TrimSpace(np.Content),
	}
	if err := post.Validate(); err != nil {
		return nil, err
	}
	return post, nil
}

func NewPostForUpdate(oldPost *Post, up dto.UpdatePost) (*Post, error) {
	post := &Post{
		BaseEntity: oldPost.BaseEntity,
		UserID:     oldPost.UserID,
		Content:    strings.TrimSpace(up.Content),
	}
	if err := post.Validate(); err != nil {
		return nil, err
	}
	post.UpdatedAt = time.Now()
	return post, nil
}

func (p *Post) Validate() error {
	if err := p.BaseEntity.Validate(); err != nil {
		return err
	}
	if p.UserID == uuid.Nil {
		return ErrUserIDEmpty
	}
	if strings.TrimSpace(p.Content) == "" {
		return ErrContentEmpty
	}
	if len(p.Content) > 5000 {
		return ErrContentTooLong
	}
	if p.CreatedAt.After(p.UpdatedAt) {
		return ErrCreatedAtAfterUpdatedAt
	}
	return nil
}

func (p *Post) UpdateContent(content string) error {
	p.Content = strings.TrimSpace(content)
	p.UpdatedAt = time.Now()
	return p.Validate()
}
