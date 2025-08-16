package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgUserRepository struct {
	db *pgxpool.Pool
}

func NewPgUserRepository(db *pgxpool.Pool) *PgUserRepository {
	return &PgUserRepository{db: db}
}

func (r *PgUserRepository) Save(ctx context.Context, u *entity.User) error {
	query := `INSERT INTO users (id, name) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query, u.ID, u.Name)
	return err
}

func (r *PgUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, name, created_at, updated_at FROM users WHERE id = $1`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return u.ToEntity(), nil
}
