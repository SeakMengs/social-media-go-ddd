package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/entity"
)

type MySQLPostRepository struct {
	baseMysqlRepository
}

func NewMySQLPostRepository(db *sql.DB) *MySQLPostRepository {
	return &MySQLPostRepository{
		baseMysqlRepository: NewBaseMysqlRepository(db),
	}
}

func (r *MySQLPostRepository) Save(ctx context.Context, p *entity.Post) error {
	query := `INSERT INTO posts (id, user_id, content) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.UserID, p.Content)
	return err
}

func (r *MySQLPostRepository) FindByID(ctx context.Context, id string) (*aggregate.Post, error) {
	query := `SELECT posts.id, posts.user_id, posts.content, posts.created_at, posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count
	FROM posts
	LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id
	) likes_count ON likes_count.post_id = posts.id
	LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id
	) favorites_count ON favorites_count.post_id = posts.id
	LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id
	) reposts_count ON reposts_count.post_id = posts.id
	WHERE posts.id = ?`

	var post aggregate.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.LikeCount,
		&post.FavoriteCount,
		&post.RepostCount,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *MySQLPostRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	query := `SELECT posts.id, posts.user_id, posts.content, posts.created_at, posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count
	FROM posts
	LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id
	) likes_count ON likes_count.post_id = posts.id
	LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id
	) favorites_count ON favorites_count.post_id = posts.id
	LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id
	) reposts_count ON reposts_count.post_id = posts.id
	WHERE posts.user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*aggregate.Post
	for rows.Next() {
		var p aggregate.Post
		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.LikeCount,
			&p.FavoriteCount,
			&p.RepostCount,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *MySQLPostRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM posts WHERE id = ? AND user_id = ?`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}

func (r *MySQLPostRepository) Update(ctx context.Context, p *entity.Post) error {
	query := `UPDATE posts SET content = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	_, err := r.db.ExecContext(ctx, query, p.Content, p.UpdatedAt, p.ID, p.UserID)
	return err
}
