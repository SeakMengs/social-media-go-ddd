package aggregate

import "social-media-go-ddd/internal/domain/entity"

// Post join with like, repost, favorite
type Post struct {
	entity.Post
	LikeCount     int `json:"likeCount"`
	RepostCount   int `json:"repostCount"`
	FavoriteCount int `json:"favoriteCount"`
}
