package entity

import "errors"

var (
	// Common errors
	ErrIDEmpty                 = errors.New("id cannot be empty")
	ErrCreatedAtEmpty          = errors.New("created_at cannot be zero")
	ErrUpdatedAtEmpty          = errors.New("updated_at cannot be zero")
	ErrCreatedAtAfterUpdatedAt = errors.New("created_at must be before updated_at")

	// User errors
	ErrUsernameEmpty = errors.New("username cannot be empty")
	ErrEmailEmpty    = errors.New("email cannot be empty")
	ErrEmailInvalid  = errors.New("invalid email format")

	// Post errors
	ErrUserIDEmpty    = errors.New("user_id cannot be null")
	ErrContentEmpty   = errors.New("content cannot be empty")
	ErrContentTooLong = errors.New("content exceeds maximum length")

	// Like errors
	ErrLikeUserIDEmpty = errors.New("user_id cannot be null")
	ErrLikePostIDEmpty = errors.New("post_id cannot be null")

	// Favorite errors
	ErrFavoriteUserIDEmpty = errors.New("user_id cannot be null")
	ErrFavoritePostIDEmpty = errors.New("post_id cannot be null")

	// Repost errors
	ErrRepostUserIDEmpty    = errors.New("user_id cannot be null")
	ErrRepostPostIDEmpty    = errors.New("post_id cannot be null")
	ErrRepostCommentTooLong = errors.New("comment exceeds maximum length")

	// Session errors
	ErrSessionUserIDEmpty = errors.New("user_id cannot be null")
	ErrSessionExpired     = errors.New("session expired")

	// Follow errors
	ErrFollowFollowerIDEmpty = errors.New("follower_id cannot be empty")
	ErrFollowFolloweeIDEmpty = errors.New("followee_id cannot be empty")
	ErrFollowSelfFollow      = errors.New("a user cannot follow themselves")
)
