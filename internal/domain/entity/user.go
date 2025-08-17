package entity

import (
	"errors"
	"social-media-go-ddd/internal/domain/valueobject"
	"strings"
	"time"
)

type User struct {
	BaseEntity
	Name     string               `json:"name"`
	Password valueobject.Password `json:"-"`
}

func NewUser(name, password string) (*User, error) {
	hashPw, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		BaseEntity: *NewBaseEntity(),
		Name:       name,
		Password:   hashPw,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *User) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New(ErrNameEmpty)
	}
	if s.CreatedAt.After(s.UpdatedAt) {
		return errors.New(ErrCreatedAtAfterUpdatedAt)
	}

	return nil
}

func (s *User) UpdatePassword(password string) error {
	hashPw, err := valueobject.NewPassword(password)
	if err != nil {
		return err
	}

	s.Password = hashPw
	s.UpdatedAt = time.Now()

	return s.Validate()
}

func (s *User) UpdateName(name string) error {
	s.Name = name
	s.UpdatedAt = time.Now()

	return s.Validate()
}
