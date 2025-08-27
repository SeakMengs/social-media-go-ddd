package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"time"

	"github.com/google/uuid"
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

func (r *PgPostRepository) getLikedStatus(ctx context.Context, postID string, userID string) (bool, error) {
	var liked bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE post_id=$1 AND user_id=$2)", postID, userID).Scan(&liked)
	if err != nil {
		return false, err
	}
	return liked, nil
}

func (r *PgPostRepository) getFavoritedStatus(ctx context.Context, postID string, userID string) (bool, error) {
	var favorited bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM favorites WHERE post_id=$1 AND user_id=$2)", postID, userID).Scan(&favorited)
	if err != nil {
		return false, err
	}
	return favorited, nil
}

func (r *PgPostRepository) FindByID(ctx context.Context, id string, currentUserID string) (*aggregate.Post, error) {
	query := `SELECT posts.id, 
       posts.user_id, 
       posts.content, 
       posts.created_at, 
       posts.updated_at,
       COALESCE(likes_count.count, 0) AS like_count,
       COALESCE(favorites_count.count, 0) AS favorite_count,
       COALESCE(reposts_count.count, 0) AS repost_count,
       users.id,
       users.username,
       users.email
	FROM posts
	    INNER JOIN users ON posts.user_id = users.id
		LEFT JOIN (
			SELECT post_id, 
				COUNT(*) AS count 
			FROM likes 
			GROUP BY post_id
		) likes_count ON likes_count.post_id = posts.id
		LEFT JOIN (
			SELECT post_id, 
				COUNT(*) AS count 
			FROM favorites 
			GROUP BY post_id
		) favorites_count ON favorites_count.post_id = posts.id
		LEFT JOIN (
			SELECT post_id, 
				COUNT(*) AS count 
			FROM reposts 
			GROUP BY post_id
		) reposts_count ON reposts_count.post_id = posts.id
	WHERE posts.id = $1`

	var post entity.Post
	var likeCount, favoriteCount, repostCount int
	var user entity.User
	var liked, favorited bool
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&likeCount,
		&favoriteCount,
		&repostCount,
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}

	if currentUserID != "" {
		liked, err = r.getLikedStatus(ctx, id, currentUserID)
		if err != nil {
			return nil, err
		}

		favorited, err = r.getFavoritedStatus(ctx, id, currentUserID)
		if err != nil {
			return nil, err
		}
	}

	return aggregate.NewPost(post, user, dto.CommonPostAggregate{
		LikeCount:     likeCount,
		FavoriteCount: favoriteCount,
		RepostCount:   repostCount,
		Liked:         liked,
		Favorited:     favorited,
	}), nil
}

func (r *PgPostRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	var user entity.User
	err := r.pool.QueryRow(ctx, "SELECT id, username, email FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	query := `SELECT 
				posts.id, 
				posts.user_id, 
				posts.content, 
				posts.created_at, 
				posts.updated_at, 
				COALESCE(likes_count.count, 0) AS like_count, 
				COALESCE(favorites_count.count, 0) AS favorite_count, 
				COALESCE(reposts_count.count, 0) AS repost_count
			FROM posts
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS count 
				FROM likes 
				GROUP BY post_id
			) likes_count ON likes_count.post_id = posts.id
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS count 
				FROM favorites 
				GROUP BY post_id
			) favorites_count ON favorites_count.post_id = posts.id
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS count 
				FROM reposts 
				GROUP BY post_id
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

		liked, err := r.getLikedStatus(ctx, post.ID.String(), userID)
		if err != nil {
			return nil, err
		}

		favorited, err := r.getFavoritedStatus(ctx, post.ID.String(), userID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, aggregate.NewPost(post, user, dto.CommonPostAggregate{
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			RepostCount:   repostCount,
			Liked:         liked,
			Favorited:     favorited,
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

func (r *PgPostRepository) getFeedTotalCount(ctx context.Context, userID string) (int, error) {
	countQuery := `
		SELECT COUNT(*) FROM (
			-- Count original posts from followed users or self
			SELECT posts.id
			FROM posts
			LEFT JOIN follows ON posts.user_id = follows.followee_id
			WHERE follows.follower_id = $1 OR posts.user_id = $1
			
			UNION ALL
			
			-- Count reposts from followed users or self
			SELECT posts.id
			FROM reposts
			INNER JOIN posts ON reposts.post_id = posts.id
			LEFT JOIN follows ON reposts.user_id = follows.followee_id
			WHERE follows.follower_id = $1 OR reposts.user_id = $1
		) AS feed_count
	`

	var total int
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *PgPostRepository) FindFeed(ctx context.Context, userID string, limit, offset int) ([]*aggregate.Post, int, error) {
	// Must return exactly the same rows, and column type of both queries to avoid sql err
	query := `
	SELECT 
		posts.id,
		posts.user_id,
		posts.content,
		posts.created_at,
		posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count,
		NULL::uuid AS repost_id,
		NULL::uuid AS repost_user_id,
		NULL::uuid AS repost_post_id,
		NULL::text AS repost_comment,
		NULL::timestamptz AS repost_created_at,
		NULL::timestamptz AS repost_updated_at,
		posts.created_at AS feed_time,  -- Use original post time for sorting
		users.id, users.username, users.email -- post owner
	FROM posts
	INNER JOIN users ON posts.user_id = users.id
	LEFT JOIN follows ON posts.user_id = follows.followee_id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = posts.id
	WHERE follows.follower_id = $1 OR posts.user_id = $1

	UNION ALL

	-- Reposts from followed users or self
	SELECT 
		posts.id,
		posts.user_id,
		posts.content,
		posts.created_at,
		posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count,
		reposts.id AS repost_id,
		reposts.user_id AS repost_user_id,
		reposts.post_id AS repost_post_id,
		reposts.comment AS repost_comment,
		reposts.created_at AS repost_created_at,
		reposts.updated_at AS repost_updated_at,
		reposts.created_at AS feed_time,  -- Use repost time for sorting
		users.id, users.username, users.email -- post owner
	FROM reposts
	INNER JOIN posts ON reposts.post_id = posts.id
	INNER JOIN users ON users.id = posts.user_id
	LEFT JOIN follows ON reposts.user_id = follows.followee_id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = posts.id
	WHERE follows.follower_id = $1 OR reposts.user_id = $1

	ORDER BY feed_time DESC
	LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feed []*aggregate.Post
	for rows.Next() {
		var post entity.Post
		var likeCount, favoriteCount, repostCount int
		var feedTime time.Time
		var postUser entity.User

		// Nullable repost fields,
		var repostID uuid.UUID
		var repostUserID uuid.UUID
		var repostPostID uuid.UUID
		// use pointer to avoid "can't scan into dest[11] (col: repost_comment): cannot scan NULL into *string"
		var repostComment *string
		var repostCreatedAt, repostUpdatedAt *time.Time

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&likeCount,
			&favoriteCount,
			&repostCount,
			&repostID,
			&repostUserID,
			&repostPostID,
			&repostComment,
			&repostCreatedAt,
			&repostUpdatedAt,
			&feedTime,
			&postUser.ID,
			&postUser.Username,
			&postUser.Email,
		)
		if err != nil {
			return nil, 0, err
		}

		commonAggregate := dto.CommonPostAggregate{
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			RepostCount:   repostCount,
		}

		// Check if this is a repost (repost_id is not null)
		if repostID != uuid.Nil {
			repost := entity.Repost{
				BaseEntity: entity.BaseEntity{
					ID:        repostID,
					CreatedAt: *repostCreatedAt,
					UpdatedAt: *repostUpdatedAt,
				},
				UserID:  repostUserID,
				PostID:  repostPostID,
				Comment: "",
			}
			if repostComment != nil {
				repost.Comment = *repostComment
			}

			var repostUser entity.User
			err = r.pool.QueryRow(ctx, "SELECT id, username, email FROM users WHERE id = $1", repostUserID).Scan(&repostUser.ID, &repostUser.Username, &repostUser.Email)
			if err != nil {
				return nil, 0, err
			}

			feed = append(feed, aggregate.NewRepost(post, &repost, postUser, &repostUser, commonAggregate))
		} else {
			// Regular post
			feed = append(feed, aggregate.NewPost(post, postUser, commonAggregate))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.getFeedTotalCount(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return feed, total, nil
}
