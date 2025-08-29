package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"time"

	"github.com/google/uuid"
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

func (r *MySQLPostRepository) FindByID(ctx context.Context, id string, currentUserID *string) (*aggregate.Post, error) {
	query := `SELECT posts.id, posts.user_id, posts.content, posts.created_at, posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count,
		 users.id,
       users.username,
       users.email,
	   -- Check if the current user has liked, favorited, or reposted the post
	   EXISTS (SELECT 1 FROM likes l WHERE l.post_id = posts.id AND l.user_id = ?) AS liked,
	   EXISTS (SELECT 1 FROM favorites f WHERE f.post_id = posts.id AND f.user_id = ?) AS favorited,
	   EXISTS (SELECT 1 FROM reposts r WHERE r.post_id = posts.id AND r.user_id = ?) AS reposted
	FROM posts
	INNER JOIN users ON posts.user_id = users.id
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

	var user User
	var post Post
	var likeCount, favoriteCount, repostCount int
	var liked, favorited, reposted bool
	err := r.db.QueryRowContext(ctx, query, currentUserID, currentUserID, currentUserID, id).Scan(
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
		&liked,
		&favorited,
		&reposted,
	)
	if err != nil {
		return nil, err
	}

	ePost, err := post.ToEntity()
	if err != nil {
		return nil, err
	}

	eUser, err := user.ToEntity()
	if err != nil {
		return nil, err
	}

	return aggregate.NewPost(*ePost, *eUser, dto.CommonPostAggregate{
		LikeCount:     likeCount,
		FavoriteCount: favoriteCount,
		RepostCount:   repostCount,
		Liked:         liked,
		Favorited:     favorited,
		Reposted:      reposted,
	}), nil
}

func (r *MySQLPostRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	var user User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email FROM users WHERE id=?", userID).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	query := `SELECT posts.id, posts.user_id, posts.content, posts.created_at, posts.updated_at,
		COALESCE(likes_count.count, 0) AS like_count,
		COALESCE(favorites_count.count, 0) AS favorite_count,
		COALESCE(reposts_count.count, 0) AS repost_count,
		-- Check if the current user has liked, favorited, or reposted the post
		EXISTS (SELECT 1 FROM likes l WHERE l.post_id = posts.id AND l.user_id = ?) AS liked,
		EXISTS (SELECT 1 FROM favorites f WHERE f.post_id = posts.id AND f.user_id = ?) AS favorited,
		EXISTS (SELECT 1 FROM reposts r WHERE r.post_id = posts.id AND r.user_id = ?) AS reposted
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

	rows, err := r.db.QueryContext(ctx, query, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*aggregate.Post
	for rows.Next() {
		var post Post
		var liked, favorited, reposted bool
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
			&liked,
			&favorited,
			&reposted,
		); err != nil {
			return nil, err
		}

		ePost, err := post.ToEntity()
		if err != nil {
			return nil, err
		}

		eUser, err := user.ToEntity()
		if err != nil {
			return nil, err
		}

		posts = append(posts, aggregate.NewPost(*ePost, *eUser, dto.CommonPostAggregate{
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			RepostCount:   repostCount,
			Liked:         liked,
			Favorited:     favorited,
			Reposted:      reposted,
		}))
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

func (r *MySQLPostRepository) getFeedTotalCount(ctx context.Context, userID string) (int, error) {
	countQuery := `
		SELECT COUNT(*) FROM (
			-- Count original posts from followed users or self
			SELECT posts.id
			FROM posts
			LEFT JOIN follows ON posts.user_id = follows.followee_id
			WHERE follows.follower_id = ? OR posts.user_id = ?

			UNION ALL
			
			-- Count reposts from followed users or self
			SELECT posts.id
			FROM reposts
			INNER JOIN posts ON reposts.post_id = posts.id
			LEFT JOIN follows ON reposts.user_id = follows.followee_id
			WHERE follows.follower_id = ? OR reposts.user_id = ?
		) AS feed_count
	`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID, userID, userID, userID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *MySQLPostRepository) FindFeed(ctx context.Context, userID string, limit, offset int) ([]*aggregate.Post, int, error) {
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
		NULL AS repost_id,
		NULL AS repost_user_id,
		NULL AS repost_post_id,
		NULL AS repost_comment,
		NULL AS repost_created_at,
		NULL AS repost_updated_at,
		posts.created_at AS feed_time,  -- Use original post time for sorting
		users.id, users.username, users.email, -- post owner
		-- Check if the current user has liked, favorited, or reposted the original post
		EXISTS (SELECT 1 FROM likes l WHERE l.post_id = posts.id AND l.user_id = ?) AS liked,
		EXISTS (SELECT 1 FROM favorites f WHERE f.post_id = posts.id AND f.user_id = ?) AS favorited,
		EXISTS (SELECT 1 FROM reposts r WHERE r.post_id = posts.id AND r.user_id = ?) AS reposted
	FROM posts
	INNER JOIN users ON posts.user_id = users.id
	LEFT JOIN follows ON posts.user_id = follows.followee_id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = posts.id
	WHERE follows.follower_id = ? OR posts.user_id = ?

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
		users.id, users.username, users.email, -- post owner
		-- Check if the current user has liked, favorited, or reposted the original post
		EXISTS (SELECT 1 FROM likes l WHERE l.post_id = posts.id AND l.user_id = ?) AS liked,
		EXISTS (SELECT 1 FROM favorites f WHERE f.post_id = posts.id AND f.user_id = ?) AS favorited,
		EXISTS (SELECT 1 FROM reposts r WHERE r.post_id = posts.id AND r.user_id = ?) AS reposted
	FROM reposts
	INNER JOIN posts ON reposts.post_id = posts.id
	INNER JOIN users ON users.id = posts.user_id
	LEFT JOIN follows ON reposts.user_id = follows.followee_id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = posts.id
	WHERE follows.follower_id = ? OR reposts.user_id = ?

	ORDER BY feed_time DESC
	LIMIT ? OFFSET ?;
	`

	rows, err := r.db.QueryContext(ctx, query, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feed []*aggregate.Post
	for rows.Next() {
		var post Post
		var likeCount, favoriteCount, repostCount int
		var liked, favorited, reposted bool
		var feedTime time.Time
		var postUser User

		// Nullable repost fields
		var repostID uuid.UUID
		var repostUserID uuid.UUID
		var repostPostID uuid.UUID
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
			&liked,
			&favorited,
			&reposted,
		)
		if err != nil {
			return nil, 0, err
		}

		commonAggregate := dto.CommonPostAggregate{
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			RepostCount:   repostCount,
			Liked:         liked,
			Favorited:     favorited,
			Reposted:      reposted,
		}

		ePostUser, err := postUser.ToEntity()
		if err != nil {
			return nil, 0, err
		}

		ePost, err := post.ToEntity()
		if err != nil {
			return nil, 0, err
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

			var repostUser User
			err = r.db.QueryRowContext(ctx, "SELECT id, username, email FROM users WHERE id = ?", repostUserID).Scan(&repostUser.ID, &repostUser.Username, &repostUser.Email)
			if err != nil {
				return nil, 0, err
			}

			eRepostUser, err := repostUser.ToEntity()
			if err != nil {
				return nil, 0, err
			}

			feed = append(feed, aggregate.NewRepost(*ePost, &repost, *ePostUser, eRepostUser, commonAggregate))
		} else {
			// Regular post
			feed = append(feed, aggregate.NewPost(*ePost, *ePostUser, commonAggregate))
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
