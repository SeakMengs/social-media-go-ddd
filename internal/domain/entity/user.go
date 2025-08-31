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

func (u *User) MarshalJSON() ([]byte, error) {
	// use alias to avoid infinite recursion because when marshal User, it will call MarshalJSON again which causes infinite recursion
	type Alias User
	return json.Marshal(&struct {
		*Alias
		Password string `json:"password"`
	}{
		Alias:    (*Alias)(u),
		Password: u.Password.GetHash(),
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	usr := &struct {
		*Alias
		Password string `json:"password"`
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &usr); err != nil {
		return err
	}

	u.Password = valueobject.PasswordFromHash(usr.Password)
	return nil
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

func (u *User) ToResponse() dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
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
