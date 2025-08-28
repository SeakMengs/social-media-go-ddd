package http

import (
	"database/sql"
	"errors"
	"social-media-go-ddd/internal/application/service"
	"social-media-go-ddd/internal/domain/aggregate"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type UserHandlerService struct {
	user     *service.UserService
	session  *service.SessionService
	post     *service.PostService
	repost   *service.RepostService
	follow   *service.FollowService
	favorite *service.FavoriteService
}

func NewUserHandlerService(user *service.UserService, session *service.SessionService, post *service.PostService, repost *service.RepostService, follow *service.FollowService, favorite *service.FavoriteService) *UserHandlerService {
	return &UserHandlerService{
		user:     user,
		session:  session,
		post:     post,
		repost:   repost,
		follow:   follow,
		favorite: favorite,
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

func NewUserHandler(userService *service.UserService, sessionService *service.SessionService, postService *service.PostService, repostService *service.RepostService, followService *service.FollowService, favoriteService *service.FavoriteService, authMiddleware *AuthMiddleware) *UserHandler {
	return &UserHandler{
		service:    NewUserHandlerService(userService, sessionService, postService, repostService, followService, favoriteService),
		middleware: NewUserHandlerMiddleware(authMiddleware),
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	apiUsersProtected := app.Group("/api/v1/users", h.middleware.auth.Handler)
	apiUsersProtected.Get("/me", h.Me)
	apiUsersProtected.Get("/me/posts", h.GetMyPosts)
	apiUsersProtected.Get("/me/feed", h.GetMyFeed)
	apiUsersProtected.Get("/me/posts/favorites", h.GetMyFavoritePosts)
	apiUsersProtected.Get("/me/reposts", h.GetMyReposts)
	apiUsersProtected.Post("/:id/follow", h.FollowUser)
	apiUsersProtected.Delete("/:id/follow", h.UnfollowUser)

	// Public user routes
	apiUsers := app.Group("/api/v1/public/users")
	apiUsers.Get("/name/:name", h.GetUsersByName)
	apiUsers.Get("/:id", h.GetUserByID)
	apiUsers.Get("/:id/posts", h.GetUserPosts)
	apiUsers.Get("/:id/reposts", h.GetUserReposts)

	// No middleware route
	apiAuth := app.Group("/api/v1/auth")
	apiAuth.Post("/register", h.CreateUser)
	apiAuth.Post("/login", h.Login)
	apiAuth.Delete("/logout", h.Logout)
}

func (h *UserHandler) getCurrentUserId(ctx *fiber.Ctx) *string {
	var currentUserID *string = nil

	// Try to read bearer token and get session/user, but it's optional
	// Such that when getting the user, we know if we have followed that person or not yet
	token, err := readBearerToken(ctx)
	if err == nil && token != "" {
		session, err := h.service.session.GetByID(ctx.Context(), token)
		if err == nil && !session.IsExpired() {
			uid := session.UserID.String()
			currentUserID = &uid
		}
	}
	return currentUserID
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
	currentUserID := h.getCurrentUserId(ctx)

	id := ctx.Params("id")
	user, err := h.service.user.GetByID(ctx.Context(), id, currentUserID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) GetUsersByName(ctx *fiber.Ctx) error {
	currentUserID := h.getCurrentUserId(ctx)

	name := ctx.Params("name")
	users, err := h.service.user.GetManyByName(ctx.Context(), name, currentUserID)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) || errors.Is(sql.ErrNoRows, err) {
			return ErrorResponse(ctx, fiber.StatusNotFound, "no users found")
		}
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"users": users,
	})
}

func (h *UserHandler) Me(ctx *fiber.Ctx) error {
	session, err := GetSessionFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	now := time.Now()
	if session.ExpireAt.Before(now) {
		_ = h.service.session.Delete(ctx.Context(), dto.DeleteSession{ID: session.ID.String()})
		return ErrorResponse(ctx, fiber.StatusUnauthorized, errors.New("session expired"))
	}

	// Refresh session if less than 1 day left
	refreshThreshold := now.Add(24 * time.Hour)
	if session.ExpireAt.Before(refreshThreshold) {
		newSession, err := h.service.session.UpdateExpireAt(ctx.Context(), session, dto.UpdateSessionExpireAt{
			ExpireAt: entity.DefaultSessionExpireAt(),
			ID:       session.ID.String(),
		})
		if err != nil {
			return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}

		return SuccessResponse(ctx, fiber.Map{
			"session": newSession,
			"user":    user,
		})
	}

	return SuccessResponse(ctx, fiber.Map{
		"session": session,
		"user":    user,
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

	user, err := h.service.user.GetByName(ctx.Context(), body.Username, nil)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) || errors.Is(sql.ErrNoRows, err) {
			return ErrorResponse(ctx, fiber.StatusNotFound, "user not found")
		}
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	if !user.Password.Match(body.Password) {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, "invalid credentials")
	}

	session, err := h.service.session.Create(ctx.Context(), dto.NewSession{
		UserID:   user.ID,
		ExpireAt: entity.DefaultSessionExpireAt(),
	})
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"session": session,
		"user":    user,
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
func (h *UserHandler) GetMyFavoritePosts(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	posts, err := h.service.favorite.GetByUserID(ctx.Context(), user.ID.String())
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

func (h *UserHandler) GetMyReposts(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	reposts, err := h.service.repost.GetByUserID(ctx.Context(), user.ID.String())
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	if reposts == nil {
		reposts = []*aggregate.Post{}
	}

	return SuccessResponse(ctx, fiber.Map{
		"reposts": reposts,
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

func (h *UserHandler) GetUserReposts(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	reposts, err := h.service.repost.GetByUserID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	if reposts == nil {
		reposts = []*aggregate.Post{}
	}

	return SuccessResponse(ctx, fiber.Map{
		"reposts": reposts,
	})
}

func (h *UserHandler) FollowUser(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	targetID := ctx.Params("id")
	userId := user.ID.String()
	targetUser, err := h.service.user.GetByID(ctx.Context(), targetID, &userId)
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
	userId := user.ID.String()
	targetUser, err := h.service.user.GetByID(ctx.Context(), targetID, &userId)
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

func (h *UserHandler) GetMyFeed(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	p, pSize := getPaginationParams(ctx)
	limit, offset := paginationToLimitOffset(p, pSize)

	feed, total, err := h.service.post.GetFeed(ctx.Context(), user.ID.String(), limit, offset)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	if feed == nil {
		feed = []*aggregate.Post{}
	}

	return SuccessResponse(ctx, fiber.Map{
		"feed": feed,
		"pagination": fiber.Map{
			"page":     p,
			"pageSize": pSize,
			"total":    total,
		},
	})
}
