package postgres

import (
	"context"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgFavoriteRepository struct {
	basePgRepository
}

func NewPgFavoriteRepository(pool *pgxpool.Pool) *PgFavoriteRepository {
	return &PgFavoriteRepository{
		basePgRepository: NewBasePgRepository(pool),
	}
}

// If already favorite, does nothing
func (r *PgFavoriteRepository) Save(ctx context.Context, f *entity.Favorite) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	checkQuery := `SELECT id FROM favorites WHERE user_id = $1 AND post_id = $2`
	row := tx.QueryRow(ctx, checkQuery, f.UserID, f.PostID)

	var existingID string
	if err := row.Scan(&existingID); err != nil {
		if err == pgx.ErrNoRows {
			// No existing favorite found, insert new favorite
			insertQuery := `INSERT INTO favorites (id, user_id, post_id) VALUES ($1, $2, $3)`
			_, err := tx.Exec(ctx, insertQuery, f.ID, f.UserID, f.PostID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PgFavoriteRepository) Delete(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND post_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, postID)
	return err
}

func (r *PgFavoriteRepository) FindByUserID(ctx context.Context, userID string) ([]*aggregate.Post, error) {
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
		WHERE favorites.user_id = $1
		ORDER BY favorites.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
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

		liked, err = getLikedStatus(ctx, r.pool, post.ID.String(), userID)
		if err != nil {
			return nil, err
		}

		favorited, err = getFavoritedStatus(ctx, r.pool, post.ID.String(), userID)
		if err != nil {
			return nil, err
		}

		reposted, err = getRepostedStatus(ctx, r.pool, post.ID.String(), userID)
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
