package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepostRepository struct {
	basePgRepository
}

func NewPgRepostRepository(pool *pgxpool.Pool) *PgRepostRepository {
	return &PgRepostRepository{
		basePgRepository: NewBasePgRepository(pool),
	}
}

func (r *PgRepostRepository) FindByID(ctx context.Context, id string) (*entity.Repost, error) {
	query := `SELECT id, user_id, post_id, comment, created_at, updated_at FROM reposts WHERE id = $1`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rp, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Repost])
	if err != nil {
		return nil, err
	}

	return rp.ToEntity()
}

func (r *PgRepostRepository) Save(ctx context.Context, rp *entity.Repost) error {
	query := `INSERT INTO reposts (id, user_id, post_id, comment) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, rp.ID, rp.UserID, rp.PostID, rp.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgRepostRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reposts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
