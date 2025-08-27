package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
)

type MySQLUserRepository struct {
	baseMysqlRepository
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{baseMysqlRepository: NewBaseMysqlRepository(db)}
}

func (r *MySQLUserRepository) Save(ctx context.Context, u *entity.User) error {
	query := `INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Username, u.Email, u.Password.GetHash())
	return err
}

func (r *MySQLUserRepository) FindByID(ctx context.Context, id string, currentUserID *string) (*aggregate.User, error) {
	query := `
		SELECT 
			users.id,
			users.username,
			users.email,
			users.password,
			users.created_at,
			users.updated_at,
			EXISTS (
				SELECT 1 FROM follows WHERE follows.follower_id = ? AND follows.followee_id = users.id
			) AS followed,
			(SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count,
			(SELECT COUNT(*) FROM follows WHERE followee_id = users.id) AS follower_count
		FROM users
		WHERE users.id = ?
	`

	row := r.db.QueryRowContext(ctx, query, currentUserID, id)

	var u User
	var followed bool
	var followingCount, followerCount int
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &followed, &followingCount, &followerCount)
	if err != nil {
		return nil, err
	}

	userEntity, err := u.ToEntity()
	if err != nil {
		return nil, err
	}

	return aggregate.NewUser(*userEntity, dto.CommonUserAggregate{
		Followed:       followed,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
	}), nil
}

func (r *MySQLUserRepository) FindByName(ctx context.Context, username string, currentUserID *string) (*aggregate.User, error) {
	query := `
		SELECT 
			users.id,
			users.username,
			users.email,
			users.password,
			users.created_at,
			users.updated_at,
			EXISTS (
				SELECT 1 FROM follows WHERE follows.follower_id = ? AND follows.followee_id = users.id
			) AS followed,
			(SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count,
			(SELECT COUNT(*) FROM follows WHERE followee_id = users.id) AS follower_count
		FROM users
		WHERE users.username = ?
	`

	row := r.db.QueryRowContext(ctx, query, currentUserID, username)

	var u User
	var followed bool
	var followingCount, followerCount int
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &followed, &followingCount, &followerCount)
	if err != nil {
		return nil, err
	}

	userEntity, err := u.ToEntity()
	if err != nil {
		return nil, err
	}

	return aggregate.NewUser(*userEntity, dto.CommonUserAggregate{
		Followed:       followed,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
	}), nil
}

func (r *MySQLUserRepository) SearchManyByName(ctx context.Context, username string, currentUserID *string) ([]*aggregate.User, error) {
	query := `
		SELECT 
			users.id,
			users.username,
			users.email,
			users.password,
			users.created_at,
			users.updated_at,
			EXISTS (
				SELECT 1 FROM follows WHERE follows.follower_id = ? AND follows.followee_id = users.id
			) AS followed,
			(SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count,
			(SELECT COUNT(*) FROM follows WHERE followee_id = users.id) AS follower_count
		FROM users
		WHERE users.username ILIKE CONCAT(?, '%')
	`

	rows, err := r.db.QueryContext(ctx, query, currentUserID, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*aggregate.User
	for rows.Next() {
		var u User
		var followed bool
		var followingCount, followerCount int
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &followed, &followingCount, &followerCount); err != nil {
			return nil, err
		}

		userEntity, err := u.ToEntity()
		if err != nil {
			return nil, err
		}

		users = append(users, aggregate.NewUser(*userEntity, dto.CommonUserAggregate{
			Followed:       followed,
			FollowingCount: followingCount,
			FollowerCount:  followerCount,
		}))
	}
	return users, nil
}
