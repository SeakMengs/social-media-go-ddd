package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/entity"
)

type MySQLRepostRepository struct {
	baseMysqlRepository
}

func NewMySQLRepostRepository(db *sql.DB) *MySQLRepostRepository {
	return &MySQLRepostRepository{
		baseMysqlRepository: NewBaseMysqlRepository(db),
	}
}

func (r *MySQLRepostRepository) FindByID(ctx context.Context, id string) (*entity.Repost, error) {
	query := `SELECT id, user_id, post_id, comment, created_at, updated_at FROM reposts WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var rp Repost
	err := row.Scan(&rp.ID, &rp.UserID, &rp.PostID, &rp.Comment, &rp.CreatedAt, &rp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return rp.ToEntity()
}

func (r *MySQLRepostRepository) Save(ctx context.Context, rp *entity.Repost) error {
	query := `INSERT INTO reposts (id, user_id, post_id, comment) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, rp.ID, rp.UserID, rp.PostID, rp.Comment)
	return err
}

func (r *MySQLRepostRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reposts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
