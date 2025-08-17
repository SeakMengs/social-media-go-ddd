package http

import (
	"errors"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/gofiber/fiber/v2"
)

// User context set by auth middleware handler
func GetUserFromCtx(c *fiber.Ctx) (*entity.User, error) {
	user := c.Locals("user")
	if user == nil {
		return nil, errors.New("user not found in context")
	}

	u, ok := user.(*entity.User)
	if !ok {
		return nil, errors.New("user in context has wrong type")
	}

	return u, nil
}

// Session context set by auth middleware handler
func GetSessionFromCtx(c *fiber.Ctx) (*entity.Session, error) {
	session := c.Locals("session")
	if session == nil {
		return nil, errors.New("session not found in context")
	}

	s, ok := session.(*entity.Session)
	if !ok {
		return nil, errors.New("session in context has wrong type")
	}

	return s, nil
}
