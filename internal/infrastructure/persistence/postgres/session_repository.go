package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgSessionRepository struct {
	basePgRepository
}

func NewPgSessionRepository(pool *pgxpool.Pool) *PgSessionRepository {
	return &PgSessionRepository{basePgRepository: NewBasePgRepository(pool)}
}

func (r *PgSessionRepository) Save(ctx context.Context, s *entity.Session) error {
	query := `INSERT INTO sessions (id, user_id, expire_at) VALUES ($1, $2, $3)`

	_, err := r.pool.Exec(ctx, query, s.ID, s.UserID, s.ExpireAt)
	return err
}

func (r *PgSessionRepository) FindByID(ctx context.Context, id string) (*entity.Session, error) {
	query := `SELECT id, user_id, expire_at, created_at, updated_at FROM sessions WHERE id = $1`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Session])
	if err != nil {
		return nil, err
	}

	return s.ToEntity()
}

func (r *PgSessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
