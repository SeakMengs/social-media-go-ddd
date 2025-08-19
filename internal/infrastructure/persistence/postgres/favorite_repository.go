package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgFavoriteRepository struct {
	basePgRepository
}

func NewPgFavoriteRepository(pool *pgxpool.Pool) *PgFavoriteRepository {
	return &PgFavoriteRepository{
		basePgRepository: NewBasePgRepository(pool),
	}
}

// If already favorite, does nothing
func (r *PgFavoriteRepository) Save(ctx context.Context, f *entity.Favorite) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	checkQuery := `SELECT id FROM favorites WHERE user_id = $1 AND post_id = $2`
	row := tx.QueryRow(ctx, checkQuery, f.UserID, f.PostID)

	var existingID string
	if err := row.Scan(&existingID); err != nil {
		if err == pgx.ErrNoRows {
			// No existing favorite found, insert new favorite
			insertQuery := `INSERT INTO favorites (id, user_id, post_id) VALUES ($1, $2, $3)`
			_, err := tx.Exec(ctx, insertQuery, f.ID, f.UserID, f.PostID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PgFavoriteRepository) Delete(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND post_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, postID)
	return err
}
