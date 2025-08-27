package http

import (
	"social-media-go-ddd/internal/application/service"
	"social-media-go-ddd/internal/domain/dto"
	"social-media-go-ddd/internal/domain/entity"

	"github.com/gofiber/fiber/v2"
)

type PostHandlerService struct {
	post     *service.PostService
	like     *service.LikeService
	repost   *service.RepostService
	favorite *service.FavoriteService
	session  *service.SessionService
}

func NewPostHandlerService(post *service.PostService, like *service.LikeService, repost *service.RepostService, favorite *service.FavoriteService, session *service.SessionService) *PostHandlerService {
	return &PostHandlerService{
		post:     post,
		like:     like,
		repost:   repost,
		favorite: favorite,
		session:  session,
	}
}

type PostHandlerMiddleware struct {
	auth *AuthMiddleware
}

func NewPostHandlerMiddleware(auth *AuthMiddleware) *PostHandlerMiddleware {
	return &PostHandlerMiddleware{
		auth: auth,
	}
}

type PostHandler struct {
	service    *PostHandlerService
	middleware *PostHandlerMiddleware
}

func NewPostHandler(postService *service.PostService, likeService *service.LikeService, repostService *service.RepostService, favoriteService *service.FavoriteService, sessionService *service.SessionService, authMiddleware *AuthMiddleware) *PostHandler {
	return &PostHandler{
		service:    NewPostHandlerService(postService, likeService, repostService, favoriteService, sessionService),
		middleware: NewPostHandlerMiddleware(authMiddleware),
	}
}

func (h *PostHandler) RegisterRoutes(app *fiber.App) {
	apiPosts := app.Group("/api/v1/public/posts")
	apiPosts.Get("/:id", h.GetPostByID)

	apiPostsProtected := app.Group("/api/v1/posts", h.middleware.auth.Handler)
	apiPostsProtected.Post("/", h.CreatePost)
	apiPostsProtected.Put("/:id", h.UpdatePost)
	apiPostsProtected.Delete("/:id", h.DeletePost)
	apiPostsProtected.Post("/:id/like", h.LikePost)
	apiPostsProtected.Delete("/:id/like", h.UnlikePost)
	apiPostsProtected.Post("/:id/favorite", h.FavoritePost)
	apiPostsProtected.Delete("/:id/favorite", h.UnfavoritePost)
	apiPostsProtected.Post("/:id/repost", h.RepostPost)
	apiPostsProtected.Delete("/:id/repost", h.UnrepostPost)
}

func (p *PostHandler) getCurrentUserId(ctx *fiber.Ctx) *string {
	var currentUserID string

	// Try to read bearer token and get session/user, but it's optional
	// Such that when getting the user, we know if we have followed that person or not yet
	token, err := readBearerToken(ctx)
	if err == nil && token != "" {
		session, err := p.service.session.GetByID(ctx.Context(), token)
		if err == nil && !session.IsExpired() {
			currentUserID = session.UserID.String()
		}
	}
	return &currentUserID
}

func (h *PostHandler) CreatePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	type request struct {
		dto.NewPost
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	body.NewPost.UserID = user.ID

	post, err := h.service.post.Create(ctx.Context(), body.NewPost)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) DeletePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}
	id := ctx.Params("id")

	userId := user.ID.String()
	post, err := h.service.post.GetByID(ctx.Context(), id, &userId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}
	if post.UserID != user.ID {
		return ErrorResponse(ctx, fiber.StatusForbidden, "you are not allowed to update this post")
	}

	if err := h.service.post.Delete(ctx.Context(), dto.DeletePost{ID: id, UserID: user.ID}); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}
	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) UpdatePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	var body dto.UpdatePost
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	body.UserID = user.ID
	body.ID = id

	userId := user.ID.String()
	post, err := h.service.post.GetByID(ctx.Context(), id, &userId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}
	if post.UserID != user.ID {
		return ErrorResponse(ctx, fiber.StatusForbidden, "you are not allowed to update this post")
	}

	updatedPost, err := h.service.post.Update(ctx.Context(), &post.Post, body)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"post": updatedPost,
	})
}

func (h *PostHandler) GetPostByID(ctx *fiber.Ctx) error {
	currentUserId := h.getCurrentUserId(ctx)

	id := ctx.Params("id")
	post, err := h.service.post.GetByID(ctx.Context(), id, currentUserId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) LikePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	userId := user.ID.String()
	post, err := h.service.post.GetByID(ctx.Context(), id, &userId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	if !post.Liked {
		_, err = h.service.like.Create(ctx.Context(), dto.NewLike{
			UserID: user.ID,
			PostID: post.ID,
		})
		if err != nil {
			return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}
	}

	// Invalidate cache for the post and user feed (post is already invalidated by like service)
	h.service.post.InvalidateCacheForUserId(ctx.Context(), userId)

	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) UnlikePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	userId := user.ID.String()
	post, err := h.service.post.GetByID(ctx.Context(), id, &userId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	if post.Liked {
		err = h.service.like.Delete(ctx.Context(), dto.DeleteLike{
			UserID: user.ID,
			PostID: post.ID,
		})
		if err != nil {
			return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}
	}

	h.service.post.InvalidateCacheForUserId(ctx.Context(), userId)

	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) FavoritePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	userId := user.ID.String()
	post, err := h.service.post.GetByID(ctx.Context(), id, &userId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	if !post.Favorited {
		_, err = h.service.favorite.Create(ctx.Context(), dto.NewFavorite{
			UserID: user.ID,
			PostID: post.ID,
		})
		if err != nil {
			return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}
	}

	h.service.post.InvalidateCacheForUserId(ctx.Context(), userId)

	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) UnfavoritePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	userId := user.ID.String()
	post, err := h.service.post.GetByID(ctx.Context(), id, &userId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	if post.Favorited {
		err = h.service.favorite.Delete(ctx.Context(), dto.DeleteFavorite{
			UserID: user.ID,
			PostID: post.ID,
		})
		if err != nil {
			return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}
	}

	h.service.post.InvalidateCacheForUserId(ctx.Context(), userId)

	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) RepostPost(ctx *fiber.Ctx) error {
	type request struct {
		dto.NewRepost
	}

	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	userId := user.ID.String()
	post, err := h.service.post.GetByID(ctx.Context(), id, &userId)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}
	body.UserID = user.ID
	body.PostID = post.ID

	repost, err := h.service.repost.Create(ctx.Context(), body.NewRepost)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"repost": repost,
	})
}

func (h *PostHandler) UnrepostPost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	postID := ctx.Params("id")
	postIDUUID, err := entity.StringToUUID(postID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	err = h.service.repost.Delete(ctx.Context(), dto.DeleteRepost{
		UserID: user.ID,
		PostID: postIDUUID,
	})
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, nil)
}
