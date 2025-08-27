package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgUserRepository struct {
	basePgRepository
}

func NewPgUserRepository(pool *pgxpool.Pool) *PgUserRepository {
	return &PgUserRepository{
		basePgRepository: NewBasePgRepository(pool),
	}
}

func (r *PgUserRepository) Save(ctx context.Context, u *entity.User) error {
	query := `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(ctx, query, u.ID, u.Username, u.Email, u.Password.GetHash())
	return err
}

type UserWithFollow struct {
	User
	Followed       bool `db:"followed"`
	FollowingCount int  `db:"following_count"`
	FollowerCount  int  `db:"follower_count"`
}

func (r *PgUserRepository) FindByID(ctx context.Context, id string, currentUserID string) (*aggregate.User, error) {
	query := `SELECT 
	    users.id,
	    users.username,
	    users.email,
	    users.password,
	    users.created_at,
	    users.updated_at,
	    EXISTS (
	        SELECT 1 FROM follows WHERE follows.follower_id = $1 AND follows.followee_id = users.id
	    ) AS followed,
	    (SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count,
	    (SELECT COUNT(*) FROM follows WHERE followee_id = users.id) AS follower_count
	FROM users
	WHERE users.id = $2`
	rows, err := r.pool.Query(ctx, query, currentUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserWithFollow])
	if err != nil {
		return nil, err
	}

	entityUser, err := u.User.ToEntity()
	if err != nil {
		return nil, err
	}
	return aggregate.NewUser(*entityUser, dto.CommonUserAggregate{
		Followed:       u.Followed,
		FollowerCount:  u.FollowerCount,
		FollowingCount: u.FollowingCount,
	}), nil
}

func (r *PgUserRepository) FindByName(ctx context.Context, username string, currentUserID string) (*aggregate.User, error) {
	query := `SELECT 
	    users.id,
	    users.username,
	    users.email,
	    users.password,
	    users.created_at,
	    users.updated_at,
	    EXISTS (
	        SELECT 1 FROM follows WHERE follows.follower_id = $1 AND follows.followee_id = users.id
	    ) AS followed,
	    (SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count,
	    (SELECT COUNT(*) FROM follows WHERE followee_id = users.id) AS follower_count
	FROM users
	WHERE users.username = $2`
	rows, err := r.pool.Query(ctx, query, currentUserID, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserWithFollow])
	if err != nil {
		return nil, err
	}

	entityUser, err := u.User.ToEntity()
	if err != nil {
		return nil, err
	}
	return aggregate.NewUser(*entityUser, dto.CommonUserAggregate{
		Followed:       u.Followed,
		FollowerCount:  u.FollowerCount,
		FollowingCount: u.FollowingCount,
	}), nil
}
