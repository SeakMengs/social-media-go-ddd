package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPostRepository struct {
	db *pgxpool.Pool
}

func NewPgPostRepository(db *pgxpool.Pool) *PgPostRepository {
	return &PgPostRepository{db: db}
}

func (r *PgPostRepository) Save(ctx context.Context, p *entity.Post) error {
	query := `INSERT INTO posts (id, user_id, content) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, p.ID, p.UserID, p.Content)
	return err
}

func (r *PgPostRepository) FindByID(ctx context.Context, id string) (*aggregate.Post, error) {
	query := `SELECT posts.id, posts.user_id, posts.content, posts.created_at, posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count FROM posts
		LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id
		) likes_count ON likes_count.post_id = posts.id
		LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id
		) favorites_count ON favorites_count.post_id = posts.id
		LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id
		) reposts_count ON reposts_count.post_id = posts.id
		WHERE posts.id = $1`
	var post aggregate.Post
	err := r.db.QueryRow(ctx, query, id).Scan(
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

func (r *PgPostRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	query := `SELECT posts.id, posts.user_id, posts.content, posts.created_at, posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count FROM posts
		LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id
		) likes_count ON likes_count.post_id = posts.id
		LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id
		) favorites_count ON favorites_count.post_id = posts.id
		LEFT JOIN (
		SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id
		) reposts_count ON reposts_count.post_id = posts.id
		WHERE posts.user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
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

func (r *PgPostRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
