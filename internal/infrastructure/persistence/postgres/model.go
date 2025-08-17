package postgres

import (
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/valueobject"

	"github.com/jackc/pgx/v5/pgtype"
)

type BaseModel struct {
	ID        pgtype.UUID        `db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

func (b *BaseModel) ToEntity() (*entity.BaseEntity, error) {
	if b == nil {
		return nil, nil
	}

	return &entity.BaseEntity{
		ID:        b.ID.Bytes,
		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
	}, nil
}

type User struct {
	BaseModel
	Name     pgtype.Text `db:"name"`
	Password pgtype.Text `db:"password"`
}

func (u *User) ToEntity() (*entity.User, error) {
	if u == nil {
		return nil, nil
	}

	baseEntity, err := u.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}

	return &entity.User{
		BaseEntity: *baseEntity,
		Name:       u.Name.String,
		Password:   valueobject.PasswordFromHash(u.Password.String),
	}, nil
}

type Session struct {
	BaseModel
	UserID   pgtype.Text        `db:"user_id"`
	ExpireAt pgtype.Timestamptz `db:"expire_at"`
}

func (s *Session) ToEntity() (*entity.Session, error) {
	if s == nil {
		return nil, nil
	}

	userId, err := entity.StringToUUID(s.UserID.String)
	if err != nil {
		return nil, err
	}

	baseEntity, err := s.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}

	return &entity.Session{
		BaseEntity: *baseEntity,
		UserID:     userId,
		ExpireAt:   s.ExpireAt.Time,
	}, nil
}
