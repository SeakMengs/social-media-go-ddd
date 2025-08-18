package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewBaseEntity() BaseEntity {
	return BaseEntity{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (b *BaseEntity) UpdateTimestamp() {
	b.UpdatedAt = time.Now()
}

func (b *BaseEntity) Validate() error {
	if b.ID == uuid.Nil {
		return errors.New(ErrIDEmpty)
	}
	if b.CreatedAt.IsZero() {
		return errors.New(ErrCreatedAtEmpty)
	}
	if b.UpdatedAt.IsZero() {
		return errors.New(ErrUpdatedAtEmpty)
	}
	if b.CreatedAt.After(b.UpdatedAt) {
		return errors.New(ErrCreatedAtAfterUpdatedAt)
	}
	return nil
}
