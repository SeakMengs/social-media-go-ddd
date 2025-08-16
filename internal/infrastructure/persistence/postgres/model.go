package postgres

import (
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgtype"
)

type BaseModel struct {
	ID        pgtype.UUID        `db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

type User struct {
	BaseModel
	Name pgtype.Text `db:"name"`
}

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		BaseEntity: entity.BaseEntity{
			ID:        u.ID.Bytes,
			CreatedAt: u.CreatedAt.Time,
			UpdatedAt: u.UpdatedAt.Time,
		},
		Name: u.Name.String,
	}
}
