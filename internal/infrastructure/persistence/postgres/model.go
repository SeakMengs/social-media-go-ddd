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

func (b *BaseModel) ToEntity() (entity.BaseEntity, error) {
	if b == nil {
		return entity.BaseEntity{}, nil
	}
	return entity.BaseEntity{
		ID:        b.ID.Bytes,
		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
	}, nil
}

type User struct {
	BaseModel
	Username pgtype.Text `db:"username"`
	Password pgtype.Text `db:"password"`
	Email    pgtype.Text `db:"email"`
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
		BaseEntity: baseEntity,
		Username:   u.Username.String,
		Password:   valueobject.PasswordFromHash(u.Password.String),
		Email:      u.Email.String,
	}, nil
}

type Post struct {
	BaseModel
	UserID  pgtype.UUID `db:"user_id"`
	Content pgtype.Text `db:"content"`
}

func (p *Post) ToEntity() (*entity.Post, error) {
	if p == nil {
		return nil, nil
	}
	baseEntity, err := p.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		BaseEntity: baseEntity,
		UserID:     p.UserID.Bytes,
		Content:    p.Content.String,
	}, nil
}

type Like struct {
	BaseModel
	UserID pgtype.UUID `db:"user_id"`
	PostID pgtype.UUID `db:"post_id"`
}

func (l *Like) ToEntity() (*entity.Like, error) {
	if l == nil {
		return nil, nil
	}
	baseEntity, err := l.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	return &entity.Like{
		BaseEntity: baseEntity,
		UserID:     l.UserID.Bytes,
		PostID:     l.PostID.Bytes,
	}, nil
}

type Favorite struct {
	BaseModel
	UserID pgtype.UUID `db:"user_id"`
	PostID pgtype.UUID `db:"post_id"`
}

func (f *Favorite) ToEntity() (*entity.Favorite, error) {
	if f == nil {
		return nil, nil
	}
	baseEntity, err := f.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	return &entity.Favorite{
		BaseEntity: baseEntity,
		UserID:     f.UserID.Bytes,
		PostID:     f.PostID.Bytes,
	}, nil
}

type Repost struct {
	BaseModel
	UserID  pgtype.UUID `db:"user_id"`
	PostID  pgtype.UUID `db:"post_id"`
	Comment pgtype.Text `db:"comment"`
}

func (r *Repost) ToEntity() (*entity.Repost, error) {
	if r == nil {
		return nil, nil
	}
	baseEntity, err := r.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	return &entity.Repost{
		BaseEntity: baseEntity,
		UserID:     r.UserID.Bytes,
		PostID:     r.PostID.Bytes,
		Comment:    r.Comment.String,
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
	userID, err := entity.StringToUUID(s.UserID.String)
	if err != nil {
		return nil, err
	}
	baseEntity, err := s.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	return &entity.Session{
		BaseEntity: baseEntity,
		UserID:     userID,
		ExpireAt:   s.ExpireAt.Time,
	}, nil
}
