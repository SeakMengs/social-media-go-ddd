package aggregate

import (
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
)

type PostType string

const (
	// When a user posts a text
	PostTypeText PostType = "text"
	// When a user reposts another post
	PostTypeRepost PostType = "repost"
)

// Post join with like, repost, favorite
type Post struct {
	entity.Post
	User          entity.User `json:"user"`
	Liked         bool        `json:"liked"`
	Favorited     bool        `json:"favorited"`
	LikeCount     int         `json:"likeCount"`
	RepostCount   int         `json:"repostCount"`
	FavoriteCount int         `json:"favoriteCount"`
	Type          PostType    `json:"type"`
	// If this post is a repost, Repost refers to the original post
	Repost     *entity.Repost `json:"repost,omitempty"`
	RepostUser *entity.User   `json:"repostUser,omitempty"`
}

func NewPost(post entity.Post, user entity.User, cpa dto.CommonPostAggregate) *Post {
	return &Post{
		Post:          post,
		User:          user,
		Liked:         cpa.Liked,
		Favorited:     cpa.Favorited,
		LikeCount:     cpa.LikeCount,
		RepostCount:   cpa.RepostCount,
		FavoriteCount: cpa.FavoriteCount,
		Type:          PostTypeText,
		// post type text which mean repost is null
		Repost: nil,
	}
}

func NewRepost(post entity.Post, repost *entity.Repost, user entity.User, repostUser *entity.User, cpa dto.CommonPostAggregate) *Post {
	return &Post{
		Post:          post,
		User:          user,
		Liked:         cpa.Liked,
		Favorited:     cpa.Favorited,
		LikeCount:     cpa.LikeCount,
		RepostCount:   cpa.RepostCount,
		FavoriteCount: cpa.FavoriteCount,
		Type:          PostTypeRepost,
		Repost:        repost,
		RepostUser:    repostUser,
	}
}
