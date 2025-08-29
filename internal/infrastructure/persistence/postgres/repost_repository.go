package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepostRepository struct {
	basePgRepository
}

func NewPgRepostRepository(pool *pgxpool.Pool) *PgRepostRepository {
	return &PgRepostRepository{
		basePgRepository: NewBasePgRepository(pool),
	}
}

func (r *PgRepostRepository) FindByID(ctx context.Context, id string) (*entity.Repost, error) {
	query := `SELECT id, user_id, post_id, comment, created_at, updated_at FROM reposts WHERE id = $1`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rp, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Repost])
	if err != nil {
		return nil, err
	}

	return rp.ToEntity()
}

// One user can only repost one, if already exist, update comment
func (r *PgRepostRepository) Save(ctx context.Context, rp *entity.Repost) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	checkQuery := `SELECT id FROM reposts WHERE user_id = $1 AND post_id = $2`
	row := tx.QueryRow(ctx, checkQuery, rp.UserID, rp.PostID)

	var existingID string
	if err := row.Scan(&existingID); err != nil {
		if err == pgx.ErrNoRows {
			// No existing repost found, insert new repost
			insertQuery := `INSERT INTO reposts (id, user_id, post_id, comment) VALUES ($1, $2, $3, $4)`
			_, err := tx.Exec(ctx, insertQuery, rp.ID, rp.UserID, rp.PostID, rp.Comment)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if existingID != "" {
		// Existing repost found, update comment
		updateQuery := `UPDATE reposts SET comment = $1 WHERE id = $2`
		_, err := tx.Exec(ctx, updateQuery, rp.Comment, existingID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PgRepostRepository) Delete(ctx context.Context, userID string, postID string) error {
	query := `DELETE FROM reposts WHERE user_id = $1 AND post_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, postID)
	return err
}

func (r *PgRepostRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	var repostUser User
	err := r.pool.QueryRow(ctx, "SELECT id, username, email FROM users WHERE id=$1", userID).
		Scan(&repostUser.ID, &repostUser.Username, &repostUser.Email)
	if err != nil {
		return nil, err
	}

	eRepostUser, err := repostUser.ToEntity()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT p.id, p.user_id, p.content, p.created_at, p.updated_at,
			COALESCE(likes_count.count, 0) AS like_count,
			COALESCE(favorites_count.count, 0) AS favorite_count,
			COALESCE(reposts_count.count, 0) AS repost_count,
			reposts.id, reposts.user_id, reposts.post_id, reposts.comment, reposts.created_at, reposts.updated_at,
			users.id, users.username, users.email,
			-- Check if the current user has liked, favorited, or reposted the original post
			EXISTS (SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $1) AS liked,
			EXISTS (SELECT 1 FROM favorites f WHERE f.post_id = p.id AND f.user_id = $1) AS favorited,
			EXISTS (SELECT 1 FROM reposts r WHERE r.post_id = p.id AND r.user_id = $1) AS reposted
		FROM reposts
		INNER JOIN users ON reposts.user_id = users.id
		INNER JOIN posts p ON reposts.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = p.id
		WHERE reposts.user_id = $1
		ORDER BY reposts.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reposts []*aggregate.Post
	for rows.Next() {
		var post Post
		var likeCount, favoriteCount, repostCount int
		var liked, favorited, reposted bool
		var repost Repost
		var user User

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
			&user.ID,
			&user.Username,
			&user.Email,
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

		eRepost, err := repost.ToEntity()
		if err != nil {
			return nil, err
		}

		eUser, err := user.ToEntity()
		if err != nil {
			return nil, err
		}

		reposts = append(reposts, aggregate.NewRepost(*ePost, eRepost, *eUser, eRepostUser, dto.CommonPostAggregate{
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

	return reposts, nil
}
