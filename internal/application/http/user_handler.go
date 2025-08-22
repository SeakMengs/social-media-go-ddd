package http

import (
	"errors"
	"social-media-go-ddd/internal/application/service"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandlerService struct {
	user    *service.UserService
	session *service.SessionService
	post    *service.PostService
	follow  *service.FollowService
}

func NewUserHandlerService(user *service.UserService, session *service.SessionService, post *service.PostService, follow *service.FollowService) *UserHandlerService {
	return &UserHandlerService{
		user:    user,
		session: session,
		post:    post,
		follow:  follow,
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

func NewUserHandler(userService *service.UserService, sessionService *service.SessionService, postService *service.PostService, followService *service.FollowService, authMiddleware *AuthMiddleware) *UserHandler {
	return &UserHandler{
		service:    NewUserHandlerService(userService, sessionService, postService, followService),
		middleware: NewUserHandlerMiddleware(authMiddleware),
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	apiUsersProtected := app.Group("/api/v1/users", h.middleware.auth.Handler)
	apiUsersProtected.Get("/me", h.Me)
	apiUsersProtected.Get("/me/posts", h.GetMyPosts)
	apiUsersProtected.Post("/:id/follow", h.FollowUser)
	apiUsersProtected.Delete("/:id/follow", h.UnfollowUser)

	// Public user routes
	apiUsers := app.Group("/api/v1/users")
	apiUsers.Get("/:id", h.GetUserByID)
	apiUsers.Get("/:id/posts", h.GetUserPosts)

	// No middleware route
	apiAuth := app.Group("/api/v1/auth")
	apiAuth.Post("/register", h.CreateUser)
	apiAuth.Post("/login", h.Login)
	apiAuth.Delete("/logout", h.Logout)
}

func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	type request struct {
		dto.NewUser
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	user, err := h.service.user.Create(ctx.Context(), body.NewUser)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) GetUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := h.service.user.GetByID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) Me(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	type request struct {
		dto.UserLogin
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	user, err := h.service.user.GetByName(ctx.Context(), body.Username)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	if !user.Password.Match(body.Password) {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, "invalid credentials")
	}

	session, err := h.service.session.Create(ctx.Context(), dto.NewSession{
		UserID:   user.ID,
		ExpireAt: time.Now().Add(24 * time.Hour),
	})
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

	if err := h.service.session.Delete(ctx.Context(), dto.DeleteSession{
		ID: session.ID.String(),
	}); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, nil)
}

func (h *UserHandler) GetMyPosts(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	posts, err := h.service.post.GetByUserID(ctx.Context(), user.ID.String())
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	if posts == nil {
		posts = []*aggregate.Post{}
	}

	return SuccessResponse(ctx, fiber.Map{
		"posts": posts,
	})
}

func (h *UserHandler) GetUserPosts(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	posts, err := h.service.post.GetByUserID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	if posts == nil {
		posts = []*aggregate.Post{}
	}

	return SuccessResponse(ctx, fiber.Map{
		"posts": posts,
	})
}

func (h *UserHandler) FollowUser(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	targetID := ctx.Params("id")
	targetUser, err := h.service.user.GetByID(ctx.Context(), targetID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	_, err = h.service.follow.Create(ctx.Context(), dto.NewFollow{
		FollowerID: user.ID,
		FolloweeID: targetUser.ID,
	})
	if err != nil {
		if errors.Is(err, entity.ErrFollowSelfFollow) {
			return ErrorResponse(ctx, fiber.StatusBadRequest, err)
		}

		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, nil)
}

func (h *UserHandler) UnfollowUser(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	targetID := ctx.Params("id")
	targetUser, err := h.service.user.GetByID(ctx.Context(), targetID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	err = h.service.follow.Delete(ctx.Context(), dto.DeleteFollow{
		FollowerID: user.ID,
		FolloweeID: targetUser.ID,
	})
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, nil)
}
