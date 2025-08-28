package service

import "fmt"

type cacheKeys struct{}

func NewCacheKeys() cacheKeys {
	return cacheKeys{}
}

func (c *cacheKeys) User(id string) string {
	return fmt.Sprintf("user:%s", id)
}

func (c *cacheKeys) UserByName(name string) string {
	return fmt.Sprintf("user:name:%s", name)
}

func (c *cacheKeys) Post(id string) string {
	return fmt.Sprintf("post:%s", id)
}

func (c *cacheKeys) Session(id string) string {
	return fmt.Sprintf("session:%s", id)
}
