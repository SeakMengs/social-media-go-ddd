package http

import (
	"errors"
	"social-media-go-ddd/internal/application/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func readAuthorizationHeader(ctx *fiber.Ctx) (string, string, error) {
	header := ctx.Get("Authorization")
	if header == "" {
		return "", "", errors.New("no authorization header specified")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", "", errors.New("wrong authorization header format")
	}

	tokenType := strings.ToUpper(parts[0])
	token := parts[1]

	if token == "" {
		return "", "", errors.New("token is empty")
	}

	return tokenType, token, nil
}

func readBearerToken(ctx *fiber.Ctx) (string, error) {
	tokenType, token, err := readAuthorizationHeader(ctx)
	if err != nil {
		return "", err
	}

	if !strings.EqualFold(tokenType, "BEARER") {
		return "", errors.New("invalid token type; expected 'Bearer'")
	}

	return token, nil
}

type AuthMiddlewareService struct {
	session *service.SessionService
	user    *service.UserService
}

func newAuthMiddlewareService(session *service.SessionService, user *service.UserService) *AuthMiddlewareService {
	return &AuthMiddlewareService{
		session: session,
		user:    user,
	}
}

type AuthMiddleware struct {
	service *AuthMiddlewareService
}

func NewAuthMiddleware(session *service.SessionService, user *service.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		service: newAuthMiddlewareService(session, user),
	}
}

func (a *AuthMiddleware) Handler(ctx *fiber.Ctx) error {
	token, err := readBearerToken(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	session, err := a.service.session.GetByID(ctx.Context(), token)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, errors.New("invalid session"))
	}

	if session.IsExpired() {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, errors.New("session expired"))
	}

	user, err := a.service.user.GetByID(ctx.Context(), session.UserID.String())
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, errors.New("invalid user"))
	}

	ctx.Locals("session", session)
	ctx.Locals("user", user)
	return ctx.Next()
}
