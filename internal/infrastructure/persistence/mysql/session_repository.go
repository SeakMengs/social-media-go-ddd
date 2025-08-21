package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/entity"
)

type MySQLSessionRepository struct {
	baseMysqlRepository
}

func NewMySQLSessionRepository(db *sql.DB) *MySQLSessionRepository {
	return &MySQLSessionRepository{
		baseMysqlRepository: NewBaseMysqlRepository(db),
	}
}

func (r *MySQLSessionRepository) Save(ctx context.Context, s *entity.Session) error {
	query := `INSERT INTO sessions (id, user_id, expire_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.UserID, s.ExpireAt)
	return err
}

func (r *MySQLSessionRepository) FindByID(ctx context.Context, id string) (*entity.Session, error) {
	query := `SELECT id, user_id, expire_at, created_at, updated_at FROM sessions WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var sess Session
	err := row.Scan(&sess.ID, &sess.UserID, &sess.ExpireAt, &sess.CreatedAt, &sess.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return sess.ToEntity()
}

func (r *MySQLSessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
