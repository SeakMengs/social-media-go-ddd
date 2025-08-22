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

type MySQLUserRepository struct {
	baseMysqlRepository
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{baseMysqlRepository: NewBaseMysqlRepository(db)}
}

func (r *MySQLUserRepository) Save(ctx context.Context, u *entity.User) error {
	query := `INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Username, u.Email, u.Password.GetHash())
	return err
}

func (r *MySQLUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u.ToEntity()
}

func (r *MySQLUserRepository) FindByName(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = ?`

	row := r.db.QueryRowContext(ctx, query, username)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u.ToEntity()
}

func (r *MySQLUserRepository) getFeedTotalCount(ctx context.Context, userID string) (int, error) {
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

func (r *MySQLUserRepository) FindFeed(ctx context.Context, userID string, limit, offset int) ([]*aggregate.Post, int, error) {
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
		posts.created_at AS feed_time  -- Use original post time for sorting
	FROM posts
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
		reposts.created_at AS feed_time  -- Use repost time for sorting
	FROM reposts
	INNER JOIN posts ON reposts.post_id = posts.id
	LEFT JOIN follows ON reposts.user_id = follows.followee_id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = posts.id
	LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = posts.id
	WHERE follows.follower_id = ? OR reposts.user_id = ?

	ORDER BY feed_time DESC
	LIMIT ? OFFSET ?;
	`

	rows, err := r.db.QueryContext(ctx, query, userID, userID, userID, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feed []*aggregate.Post
	for rows.Next() {
		var post entity.Post
		var likeCount, favoriteCount, repostCount int
		var feedTime time.Time

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

			feed = append(feed, aggregate.NewRepost(post, &repost, commonAggregate))
		} else {
			// Regular post
			feed = append(feed, aggregate.NewPost(post, commonAggregate))
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
