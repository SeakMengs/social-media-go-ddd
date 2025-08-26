package entity

import (
	"encoding/json"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/valueobject"
	"strings"
)

type User struct {
	BaseEntity
	Username string               `json:"username"`
	Email    string               `json:"email"`
	Password valueobject.Password `json:"-"`
}

// MarshalJSON includes the password field.
func (u *User) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		User
		Password string `json:"password"`
	}{
		User:     *u,
		Password: u.Password.GetHash(),
	})
}

// UnmarshalJSON includes the password field.
func UserUnmarshalJson(data []byte) (*User, error) {
	u := struct {
		User
		Password string `json:"password"`
	}{}
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	u.User.Password = valueobject.PasswordFromHash(u.Password)
	return &u.User, nil
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
		return ErrUsernameEmpty
	}
	if strings.TrimSpace(u.Email) == "" {
		return ErrEmailEmpty
	}
	if !strings.Contains(u.Email, "@") {
		return ErrEmailInvalid
	}

	return nil
}

func (u *User) UpdatePassword(password string) error {
	hashPw, err := valueobject.NewPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashPw
	u.UpdateTimestamp()
	return u.Validate()
}

func (u *User) UpdateUsername(username string) error {
	u.Username = strings.TrimSpace(username)
	u.UpdateTimestamp()
	return u.Validate()
}

func (u *User) UpdateEmail(email string) error {
	u.Email = strings.TrimSpace(strings.ToLower(email))
	u.UpdateTimestamp()
	return u.Validate()
}
