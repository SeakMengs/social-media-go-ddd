package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseEntity
	Name string
}

func NewUser(name string) (*User, error) {
	user := &User{
		BaseEntity: BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: name,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *User) Validate() error {
	if s.Name == "" {
		return errors.New(ErrNameEmpty)
	}
	if s.CreatedAt.After(s.UpdatedAt) {
		return errors.New(ErrCreatedAtAfterUpdatedAt)
	}

	return nil
}

func (s *User) UpdateName(name string) error {
	s.Name = name
	s.UpdatedAt = time.Now()

	return s.Validate()
}
