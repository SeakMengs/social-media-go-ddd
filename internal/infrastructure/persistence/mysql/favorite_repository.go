package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/entity"
)

type MySQLFavoriteRepository struct {
	baseMysqlRepository
}

func NewMySQLFavoriteRepository(db *sql.DB) *MySQLFavoriteRepository {
	return &MySQLFavoriteRepository{
		baseMysqlRepository: NewBaseMysqlRepository(db),
	}
}

// If already favorited, does nothing
func (r *MySQLFavoriteRepository) Save(ctx context.Context, f *entity.Favorite) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	checkQuery := `SELECT id FROM favorites WHERE user_id = ? AND post_id = ?`
	row := tx.QueryRowContext(ctx, checkQuery, f.UserID, f.PostID)

	var existingID string
	err = row.Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing favorite found, insert new favorite
			insertQuery := `INSERT INTO favorites (id, user_id, post_id) VALUES (?, ?, ?)`
			_, err := tx.ExecContext(ctx, insertQuery, f.ID, f.UserID, f.PostID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLFavoriteRepository) Delete(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM favorites WHERE user_id = ? AND post_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, postID)
	return err
}
