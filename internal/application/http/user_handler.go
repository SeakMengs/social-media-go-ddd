package http

import (
	"errors"
	"social-media-go-ddd/internal/application/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandlerService struct {
	user    *service.UserService
	session *service.SessionService
}

func NewUserHandlerService(user *service.UserService, session *service.SessionService) *UserHandlerService {
	return &UserHandlerService{
		user:    user,
		session: session,
	}
}

type UserHandlerMiddleware struct {
	auth *AuthMiddleware
}

func NewUserHandlerMiddleware(auth *AuthMiddleware) *UserHandlerMiddleware {
	return &UserHandlerMiddleware{
		auth: auth,
	}
}

type UserHandler struct {
	service    *UserHandlerService
	middleware *UserHandlerMiddleware
}

func NewUserHandler(userService *service.UserService, sessionService *service.SessionService, authMiddleware *AuthMiddleware) *UserHandler {
	return &UserHandler{
		service:    NewUserHandlerService(userService, sessionService),
		middleware: NewUserHandlerMiddleware(authMiddleware),
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	apiUsers := app.Group("/api/v1/users", h.middleware.auth.Handler)
	apiUsers.Get("/me", h.Me)
	apiUsers.Get("/:id", h.GetUserByID)

	apiAuth := app.Group("/api/v1/auth")
	apiAuth.Post("/register", h.CreateUser)
	apiAuth.Post("/login", h.Login)
	apiAuth.Delete("/logout", h.Logout)
}

func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	type request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	user, err := h.service.user.Create(ctx.Context(), body.Name, body.Password)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, user)
}

func (h *UserHandler) GetUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := h.service.user.GetByID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	return SuccessResponse(ctx, user)
}

func (h *UserHandler) Me(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	return SuccessResponse(ctx, user)
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	type request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	user, err := h.service.user.GetByName(ctx.Context(), body.Name)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	if !user.Password.Match(body.Password) {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, "invalid credentials")
	}

	session, err := h.service.session.Create(ctx.Context(), user.ID, time.Now().Add(1*time.Second))
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"session": fiber.Map{
			"id":        session.ID.String(),
			"expire_at": session.ExpireAt,
		},
		"user": user,
	})
}

func (h *UserHandler) Logout(ctx *fiber.Ctx) error {
	token, err := readBearerToken(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	session, err := h.service.session.GetByID(ctx.Context(), token)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, errors.New("invalid session"))
	}

	if err := h.service.session.Delete(ctx.Context(), session.ID.String()); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, nil)
}
