package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPostRepository struct {
	basePgRepository
}

func NewPgPostRepository(pool *pgxpool.Pool) *PgPostRepository {
	return &PgPostRepository{basePgRepository: NewBasePgRepository(pool)}
}

func (r *PgPostRepository) Save(ctx context.Context, p *entity.Post) error {
	query := `INSERT INTO posts (id, user_id, content) VALUES ($1, $2, $3)`

	_, err := r.pool.Exec(ctx, query, p.ID, p.UserID, p.Content)
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

	var post entity.Post
	var likeCount, favoriteCount, repostCount int
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&likeCount,
		&favoriteCount,
		&repostCount,
	)
	if err != nil {
		return nil, err
	}

	return aggregate.NewPost(post, dto.CommonPostAggregate{
		LikeCount:     likeCount,
		FavoriteCount: favoriteCount,
		RepostCount:   repostCount,
	}), nil
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
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*aggregate.Post
	for rows.Next() {
		var post entity.Post
		var likeCount, favoriteCount, repostCount int

		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&likeCount,
			&favoriteCount,
			&repostCount,
		); err != nil {
			return nil, err
		}

		posts = append(posts, aggregate.NewPost(post, dto.CommonPostAggregate{
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			RepostCount:   repostCount,
		}))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PgPostRepository) Delete(ctx context.Context, id string, userId string) error {
	query := `DELETE FROM posts WHERE id = $1 AND user_id = $2`
	_, err := r.pool.Exec(ctx, query, id, userId)
	return err
}

func (r *PgPostRepository) Update(ctx context.Context, p *entity.Post) error {
	query := `UPDATE posts SET content = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`
	_, err := r.pool.Exec(ctx, query, p.Content, p.UpdatedAt, p.ID, p.UserID)
	return err
}
