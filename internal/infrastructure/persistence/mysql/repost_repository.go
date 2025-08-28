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

// One user can only repost one, if already exist, update comment
func (r *MySQLRepostRepository) Save(ctx context.Context, rp *entity.Repost) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	checkQuery := `SELECT id FROM reposts WHERE user_id = ? AND post_id = ?`
	row := tx.QueryRowContext(ctx, checkQuery, rp.UserID, rp.PostID)

	var existingID string
	err = row.Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing repost found, insert new repost
			insertQuery := `INSERT INTO reposts (id, user_id, post_id, comment) VALUES (?, ?, ?, ?)`
			_, err := tx.ExecContext(ctx, insertQuery, rp.ID, rp.UserID, rp.PostID, rp.Comment)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Existing repost found, update comment
		updateQuery := `UPDATE reposts SET comment = ? WHERE id = ?`
		_, err := tx.ExecContext(ctx, updateQuery, rp.Comment, existingID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLRepostRepository) Delete(ctx context.Context, userID string, postID string) error {
	query := `DELETE FROM reposts WHERE user_id = ? AND post_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, postID)
	return err
}

func (r *MySQLRepostRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	var repostUser User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email FROM users WHERE id=?", userID).
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
			users.id, users.username, users.email
		FROM reposts
		INNER JOIN users ON reposts.user_id = users.id
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

		liked, err = getLikedStatus(ctx, r.db, post.ID, userID)
		if err != nil {
			return nil, err
		}

		favorited, err = getFavoritedStatus(ctx, r.db, post.ID, userID)
		if err != nil {
			return nil, err
		}

		reposted, err = getRepostedStatus(ctx, r.db, post.ID, userID)
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
