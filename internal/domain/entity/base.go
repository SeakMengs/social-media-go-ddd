package entity

import (
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
		return ErrIDEmpty
	}
	if b.CreatedAt.IsZero() {
		return ErrCreatedAtEmpty
	}
	if b.UpdatedAt.IsZero() {
		return ErrUpdatedAtEmpty
	}
	if b.CreatedAt.After(b.UpdatedAt) {
		return ErrCreatedAtAfterUpdatedAt
	}
	return nil
}
