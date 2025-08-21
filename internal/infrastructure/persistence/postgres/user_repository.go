package postgres

import (
	"context"
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

func (r *PgUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return u.ToEntity()
}

func (r *PgUserRepository) FindByName(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1`
	rows, err := r.pool.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return u.ToEntity()
}
