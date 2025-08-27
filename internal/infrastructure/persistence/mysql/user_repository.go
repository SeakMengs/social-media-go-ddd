package mysql

import (
	"context"
	"database/sql"
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

func (r *MySQLUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u.ToEntity()
}

func (r *MySQLUserRepository) FindByName(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = ?`

	row := r.db.QueryRowContext(ctx, query, username)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u.ToEntity()
}
