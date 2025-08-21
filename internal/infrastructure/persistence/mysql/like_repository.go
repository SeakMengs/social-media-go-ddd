package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/entity"
)

type MySQLLikeRepository struct {
	baseMysqlRepository
}

func NewMySQLLikeRepository(db *sql.DB) *MySQLLikeRepository {
	return &MySQLLikeRepository{
		baseMysqlRepository: NewBaseMysqlRepository(db),
	}
}

// If already liked, does nothing
func (r *MySQLLikeRepository) Save(ctx context.Context, l *entity.Like) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	checkQuery := `SELECT id FROM likes WHERE user_id = ? AND post_id = ?`
	row := tx.QueryRowContext(ctx, checkQuery, l.UserID, l.PostID)

	var existingID string
	err = row.Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing like found, insert new like
			insertQuery := `INSERT INTO likes (id, user_id, post_id) VALUES (?, ?, ?)`
			_, err := tx.ExecContext(ctx, insertQuery, l.ID, l.UserID, l.PostID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLLikeRepository) Delete(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM likes WHERE user_id = ? AND post_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, postID)
	return err
}
