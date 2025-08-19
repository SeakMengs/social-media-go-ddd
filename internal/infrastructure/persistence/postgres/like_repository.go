package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgLikeRepository struct {
	basePgRepository
}

func NewPgLikeRepository(pool *pgxpool.Pool) *PgLikeRepository {
	return &PgLikeRepository{
		basePgRepository: NewBasePgRepository(pool),
	}
}

// If already like, does nothing
func (r *PgLikeRepository) Save(ctx context.Context, l *entity.Like) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	checkQuery := `SELECT id FROM likes WHERE user_id = $1 AND post_id = $2`
	row := tx.QueryRow(ctx, checkQuery, l.UserID, l.PostID)

	var existingID string
	if err := row.Scan(&existingID); err != nil {
		if err == pgx.ErrNoRows {
			// No existing like found, insert new like
			insertQuery := `INSERT INTO likes (id, user_id, post_id) VALUES ($1, $2, $3)`
			_, err := tx.Exec(ctx, insertQuery, l.ID, l.UserID, l.PostID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PgLikeRepository) Delete(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, postID)
	return err
}
