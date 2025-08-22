package mysql

import (
	"context"
	"social-media-go-ddd/internal/domain/entity"

	"database/sql"
)

type MySQLFollowRepository struct {
	baseMysqlRepository
}

func NewMySQLFollowRepository(db *sql.DB) *MySQLFollowRepository {
	return &MySQLFollowRepository{
		baseMysqlRepository: NewBaseMysqlRepository(db),
	}
}

// If already follow, does nothing
func (r *MySQLFollowRepository) Save(ctx context.Context, f *entity.Follow) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	checkQuery := `SELECT id FROM follows WHERE follower_id = ? AND followee_id = ?`
	row := tx.QueryRowContext(ctx, checkQuery, f.FollowerID, f.FolloweeID)

	var existingID string
	if err := row.Scan(&existingID); err != nil {
		if err == sql.ErrNoRows {
			// No existing follow found, insert new follow
			insertQuery := `INSERT INTO follows (id, follower_id, followee_id) VALUES (?, ?, ?)`
			_, err := tx.ExecContext(ctx, insertQuery, f.ID, f.FollowerID, f.FolloweeID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLFollowRepository) Delete(ctx context.Context, followerID, followeeID string) error {
	query := `DELETE FROM follows WHERE follower_id = ? AND followee_id = ?`
	_, err := r.db.ExecContext(ctx, query, followerID, followeeID)
	return err
}
