package aggregate

import (
	"encoding/json"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

type User struct {
	// cannot embed entity user struct becaue it implement it own marshal json
	ID             uuid.UUID            `json:"id"`
	Username       string               `json:"username"`
	Email          string               `json:"email"`
	Password       valueobject.Password `json:"-"`
	CreatedAt      time.Time            `json:"createdAt,omitempty"`
	UpdatedAt      time.Time            `json:"updatedAt,omitempty"`
	Followed       bool                 `json:"followed"`
	FollowerCount  int                  `json:"followerCount"`
	FollowingCount int                  `json:"followingCount"`
}

func NewUser(u entity.User, cgu dto.CommonUserAggregate) *User {
	return &User{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		Password:       u.Password,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Followed:       cgu.Followed,
		FollowerCount:  cgu.FollowerCount,
		FollowingCount: cgu.FollowingCount,
	}
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
