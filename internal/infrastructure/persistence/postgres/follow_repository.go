package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgFollowRepository struct {
	basePgRepository
}

func NewPgFollowRepository(pool *pgxpool.Pool) *PgFollowRepository {
	return &PgFollowRepository{
		basePgRepository: NewBasePgRepository(pool),
	}
}

// If already follow, does nothing
func (r *PgFollowRepository) Save(ctx context.Context, f *entity.Follow) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	checkQuery := `SELECT id FROM follows WHERE follower_id = $1 AND followee_id = $2`
	row := tx.QueryRow(ctx, checkQuery, f.FollowerID, f.FolloweeID)

	var existingID string
	if err := row.Scan(&existingID); err != nil {
		if err == pgx.ErrNoRows {
			// No existing follow found, insert new follow
			insertQuery := `INSERT INTO follows (id, follower_id, followee_id) VALUES ($1, $2, $3)`
			_, err := tx.Exec(ctx, insertQuery, f.ID, f.FollowerID, f.FolloweeID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PgFollowRepository) Delete(ctx context.Context, followerID, followeeID string) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2`
	_, err := r.pool.Exec(ctx, query, followerID, followeeID)
	return err
}
