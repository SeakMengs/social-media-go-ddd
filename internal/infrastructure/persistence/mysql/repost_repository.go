package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
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

func (r *MySQLRepostRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.content, p.created_at, p.updated_at,
			COALESCE(likes_count.count, 0) AS like_count,
			COALESCE(favorites_count.count, 0) AS favorite_count,
			COALESCE(reposts_count.count, 0) AS repost_count,
			reposts.id, reposts.user_id, reposts.post_id, reposts.comment, reposts.created_at, reposts.updated_at
		FROM reposts
		INNER JOIN posts p ON reposts.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = p.id
		WHERE reposts.user_id = ?
		ORDER BY reposts.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reposts []*aggregate.Post
	for rows.Next() {
		var post entity.Post
		var likeCount, favoriteCount, repostCount int
		var repost entity.Repost

		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&likeCount,
			&favoriteCount,
			&repostCount,
			&repost.ID,
			&repost.UserID,
			&repost.PostID,
			&repost.Comment,
			&repost.CreatedAt,
			&repost.UpdatedAt,
		); err != nil {
			return nil, err
		}

		reposts = append(reposts, aggregate.NewRepost(post, &repost, dto.CommonPostAggregate{
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			RepostCount:   repostCount,
		}))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reposts, nil
}
