package entity

import (
	"errors"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/valueobject"
	"strings"
	"time"
)

type User struct {
	BaseEntity
	Username string               `json:"username"`
	Email    string               `json:"email"`
	Password valueobject.Password `json:"-"`
}

func NewUser(nu dto.NewUser) (*User, error) {
	hashPw, err := valueobject.NewPassword(nu.Password)
	if err != nil {
		return nil, err
	}
	user := &User{
		BaseEntity: NewBaseEntity(),
		Username:   strings.TrimSpace(nu.Username),
		Email:      strings.TrimSpace(strings.ToLower(nu.Email)),
		Password:   hashPw,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate() error {
	if err := u.BaseEntity.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(u.Username) == "" {
		return errors.New(ErrUsernameEmpty)
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New(ErrEmailEmpty)
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New(ErrEmailInvalid)
	}

	return nil
}

func (u *User) UpdatePassword(password string) error {
	hashPw, err := valueobject.NewPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashPw
	u.UpdatedAt = time.Now()
	return u.Validate()
}

func (u *User) UpdateUsername(username string) error {
	u.Username = strings.TrimSpace(username)
	u.UpdatedAt = time.Now()
	return u.Validate()
}

func (u *User) UpdateEmail(email string) error {
	u.Email = strings.TrimSpace(strings.ToLower(email))
	u.UpdatedAt = time.Now()
	return u.Validate()
}
