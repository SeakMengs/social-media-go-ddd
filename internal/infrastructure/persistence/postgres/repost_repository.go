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

func (r *PgRepostRepository) Save(ctx context.Context, rp *entity.Repost) error {
	query := `INSERT INTO reposts (id, user_id, post_id, comment) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, rp.ID, rp.UserID, rp.PostID, rp.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgRepostRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reposts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
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
			users.id, users.username, users.email
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

		reposts = append(reposts, aggregate.NewRepost(*ePost, eRepost, *eUser, eRepostUser, dto.CommonPostAggregate{
			LikeCount:     likeCount,
			FavoriteCount: favoriteCount,
			RepostCount:   repostCount,
		}))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reposts, nil
}
