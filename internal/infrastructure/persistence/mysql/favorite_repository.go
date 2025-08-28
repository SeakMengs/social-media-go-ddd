package mysql

import (
	"context"
	"database/sql"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
)

type MySQLFavoriteRepository struct {
	baseMysqlRepository
}

func NewMySQLFavoriteRepository(db *sql.DB) *MySQLFavoriteRepository {
	return &MySQLFavoriteRepository{
		baseMysqlRepository: NewBaseMysqlRepository(db),
	}
}

// If already favorited, does nothing
func (r *MySQLFavoriteRepository) Save(ctx context.Context, f *entity.Favorite) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	checkQuery := `SELECT id FROM favorites WHERE user_id = ? AND post_id = ?`
	row := tx.QueryRowContext(ctx, checkQuery, f.UserID, f.PostID)

	var existingID string
	err = row.Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing favorite found, insert new favorite
			insertQuery := `INSERT INTO favorites (id, user_id, post_id) VALUES (?, ?, ?)`
			_, err := tx.ExecContext(ctx, insertQuery, f.ID, f.UserID, f.PostID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLFavoriteRepository) Delete(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM favorites WHERE user_id = ? AND post_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, postID)
	return err
}

func (r *MySQLFavoriteRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.content, p.created_at, p.updated_at,
			COALESCE(likes_count.count, 0) AS like_count,
			COALESCE(favorites_count.count, 0) AS favorite_count,
			COALESCE(reposts_count.count, 0) AS repost_count,
			users.id, users.username, users.email
		FROM favorites
		INNER JOIN users ON favorites.user_id = users.id
		INNER JOIN posts p ON favorites.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM likes GROUP BY post_id) likes_count ON likes_count.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM favorites GROUP BY post_id) favorites_count ON favorites_count.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS count FROM reposts GROUP BY post_id) reposts_count ON reposts_count.post_id = p.id
		WHERE favorites.user_id = ?
		ORDER BY favorites.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*aggregate.Post
	for rows.Next() {
		var post Post
		var likeCount, favoriteCount, repostCount int
		var liked, favorited, reposted bool
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
