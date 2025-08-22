package entity

import (
	"social-media-go-ddd/internal/domain/dto"

	"github.com/google/uuid"
)

type Follow struct {
	BaseEntity
	FollowerID uuid.UUID `json:"follower_id"`
	FolloweeID uuid.UUID `json:"followee_id"`
}

func NewFollow(nf dto.NewFollow) (*Follow, error) {
	follow := &Follow{
		BaseEntity: NewBaseEntity(),
		FollowerID: nf.FollowerID,
		FolloweeID: nf.FolloweeID,
	}
	if err := follow.Validate(); err != nil {
		return nil, err
	}
	return follow, nil
}

func (f *Follow) Validate() error {
	if err := f.BaseEntity.Validate(); err != nil {
		return err
	}
	if f.FollowerID == uuid.Nil {
		return ErrFollowFollowerIDEmpty
	}
	if f.FolloweeID == uuid.Nil {
		return ErrFollowFolloweeIDEmpty
	}
	if f.FollowerID == f.FolloweeID {
		return ErrFollowSelfFollow
	}
	return nil
}
