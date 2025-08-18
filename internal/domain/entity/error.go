package entity

const (
	// Common errors
	ErrIDEmpty                 = "id cannot be empty"
	ErrCreatedAtEmpty          = "created_at cannot be zero"
	ErrUpdatedAtEmpty          = "updated_at cannot be zero"
	ErrCreatedAtAfterUpdatedAt = "created_at must be before updated_at"

	// User errors
	ErrUsernameEmpty = "username cannot be empty"
	ErrEmailEmpty    = "email cannot be empty"
	ErrEmailInvalid  = "invalid email format"

	// Post errors
	ErrUserIDEmpty    = "user_id cannot be null"
	ErrContentEmpty   = "content cannot be empty"
	ErrContentTooLong = "content exceeds maximum length"

	// Like errors
	ErrLikeUserIDEmpty = "user_id cannot be null"
	ErrLikePostIDEmpty = "post_id cannot be null"

	// Favorite errors
	ErrFavoriteUserIDEmpty = "user_id cannot be null"
	ErrFavoritePostIDEmpty = "post_id cannot be null"

	// Repost errors
	ErrRepostUserIDEmpty    = "user_id cannot be null"
	ErrRepostPostIDEmpty    = "post_id cannot be null"
	ErrRepostCommentTooLong = "comment exceeds maximum length"

	// Session errors
	ErrSessionUserIDEmpty   = "user_id cannot be null"
	ErrSessionExpiredInPast = "session cannot expire in the past"
)
