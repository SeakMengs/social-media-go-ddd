package mysql

import (
	"social-media-go-ddd/internal/domain/entity"
	"social-media-go-ddd/internal/domain/valueobject"
	"time"
)

type BaseModel struct {
	ID        string    `db:"id"` // UUID stored as CHAR(36)
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (b *BaseModel) ToEntity() (entity.BaseEntity, error) {
	if b == nil {
		return entity.BaseEntity{}, nil
	}
	uuidID, err := entity.StringToUUID(b.ID)
	if err != nil {
		return entity.BaseEntity{}, err
	}
	return entity.BaseEntity{
		ID:        uuidID,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}, nil
}

type User struct {
	BaseModel
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
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
		Username:   u.Username,
		Password:   valueobject.PasswordFromHash(u.Password),
		Email:      u.Email,
	}, nil
}

type Post struct {
	BaseModel
	UserID  string `db:"user_id"`
	Content string `db:"content"`
}

func (p *Post) ToEntity() (*entity.Post, error) {
	if p == nil {
		return nil, nil
	}
	baseEntity, err := p.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	userID, err := entity.StringToUUID(p.UserID)
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		BaseEntity: baseEntity,
		UserID:     userID,
		Content:    p.Content,
	}, nil
}

type Like struct {
	BaseModel
	UserID string `db:"user_id"`
	PostID string `db:"post_id"`
}

func (l *Like) ToEntity() (*entity.Like, error) {
	if l == nil {
		return nil, nil
	}
	baseEntity, err := l.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	userID, err := entity.StringToUUID(l.UserID)
	if err != nil {
		return nil, err
	}
	postID, err := entity.StringToUUID(l.PostID)
	if err != nil {
		return nil, err
	}
	return &entity.Like{
		BaseEntity: baseEntity,
		UserID:     userID,
		PostID:     postID,
	}, nil
}

type Favorite struct {
	BaseModel
	UserID string `db:"user_id"`
	PostID string `db:"post_id"`
}

func (f *Favorite) ToEntity() (*entity.Favorite, error) {
	if f == nil {
		return nil, nil
	}
	baseEntity, err := f.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	userID, err := entity.StringToUUID(f.UserID)
	if err != nil {
		return nil, err
	}
	postID, err := entity.StringToUUID(f.PostID)
	if err != nil {
		return nil, err
	}
	return &entity.Favorite{
		BaseEntity: baseEntity,
		UserID:     userID,
		PostID:     postID,
	}, nil
}

type Repost struct {
	BaseModel
	UserID  string `db:"user_id"`
	PostID  string `db:"post_id"`
	Comment string `db:"comment"`
}

func (r *Repost) ToEntity() (*entity.Repost, error) {
	if r == nil {
		return nil, nil
	}
	baseEntity, err := r.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	userID, err := entity.StringToUUID(r.UserID)
	if err != nil {
		return nil, err
	}
	postID, err := entity.StringToUUID(r.PostID)
	if err != nil {
		return nil, err
	}
	return &entity.Repost{
		BaseEntity: baseEntity,
		UserID:     userID,
		PostID:     postID,
		Comment:    r.Comment,
	}, nil
}

type Session struct {
	BaseModel
	UserID   string    `db:"user_id"`
	ExpireAt time.Time `db:"expire_at"`
}

func (s *Session) ToEntity() (*entity.Session, error) {
	if s == nil {
		return nil, nil
	}
	userID, err := entity.StringToUUID(s.UserID)
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
		ExpireAt:   s.ExpireAt,
	}, nil
}

type Follow struct {
	BaseModel
	FollowerID string `db:"follower_id"`
	FolloweeID string `db:"followee_id"`
}

func (f *Follow) ToEntity() (*entity.Follow, error) {
	if f == nil {
		return nil, nil
	}
	baseEntity, err := f.BaseModel.ToEntity()
	if err != nil {
		return nil, err
	}
	followerID, err := entity.StringToUUID(f.FollowerID)
	if err != nil {
		return nil, err
	}
	followeeID, err := entity.StringToUUID(f.FolloweeID)
	if err != nil {
		return nil, err
	}
	return &entity.Follow{
		BaseEntity: baseEntity,
		FollowerID: followerID,
		FolloweeID: followeeID,
	}, nil
}
