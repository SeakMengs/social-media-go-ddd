package service

import "fmt"

type cacheKeys struct{}

func NewCacheKeys() cacheKeys {
	return cacheKeys{}
}

// User
func (c *cacheKeys) User(id string) string {
	return fmt.Sprintf("user:%s", id)
}

func (c *cacheKeys) UserByName(name string) string {
	return fmt.Sprintf("user:name:%s", name)
}

func (c *cacheKeys) UserFeed(userID string, limit, offset int) string {
	return fmt.Sprintf("user:feed:%s:%d:%d", userID, limit, offset)
}

func (c *cacheKeys) UserFeedPattern(userId string) string {
	return fmt.Sprintf("user:feed:%s:*", userId)
}

func (c *cacheKeys) UserPosts(userID string) string {
	return fmt.Sprintf("user:posts:%s", userID)
}

func (c *cacheKeys) UserReposts(userID string) string {
	return fmt.Sprintf("user:reposts:%s", userID)
}

// Post
func (c *cacheKeys) Post(id string) string {
	return fmt.Sprintf("post:%s", id)
}

// Session
func (c *cacheKeys) Session(id string) string {
	return fmt.Sprintf("session:%s", id)
}

// Repost
func (c *cacheKeys) Repost(id string) string {
	return fmt.Sprintf("repost:%s", id)
}
