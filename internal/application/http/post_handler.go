package http

import (
	"social-media-go-ddd/internal/application/service"
	"social-media-go-ddd/internal/domain/dto"

	"github.com/gofiber/fiber/v2"
)

type PostHandlerService struct {
	post     *service.PostService
	like     *service.LikeService
	repost   *service.RepostService
	favorite *service.FavoriteService
}

func NewPostHandlerService(post *service.PostService, like *service.LikeService, repost *service.RepostService, favorite *service.FavoriteService) *PostHandlerService {
	return &PostHandlerService{
		post:     post,
		like:     like,
		repost:   repost,
		favorite: favorite,
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

func NewPostHandler(postService *service.PostService, likeService *service.LikeService, repostService *service.RepostService, favoriteService *service.FavoriteService, authMiddleware *AuthMiddleware) *PostHandler {
	return &PostHandler{
		service:    NewPostHandlerService(postService, likeService, repostService, favoriteService),
		middleware: NewPostHandlerMiddleware(authMiddleware),
	}
}

func (h *PostHandler) RegisterRoutes(app *fiber.App) {
	apiPosts := app.Group("/api/v1/posts")
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

	apiRepostsProtected := app.Group("/api/v1/reposts", h.middleware.auth.Handler)
	apiRepostsProtected.Delete("/:id", h.UnrepostPost)
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

	post, err := h.service.post.GetByID(ctx.Context(), id, user.ID.String())
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

	post, err := h.service.post.GetByID(ctx.Context(), id, user.ID.String())
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
	id := ctx.Params("id")
	post, err := h.service.post.GetByID(ctx.Context(), id, "")
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
	post, err := h.service.post.GetByID(ctx.Context(), id, user.ID.String())
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

	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) UnlikePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	post, err := h.service.post.GetByID(ctx.Context(), id, user.ID.String())
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

	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) FavoritePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	post, err := h.service.post.GetByID(ctx.Context(), id, user.ID.String())
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

	return SuccessResponse(ctx, nil)
}

func (h *PostHandler) UnfavoritePost(ctx *fiber.Ctx) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	post, err := h.service.post.GetByID(ctx.Context(), id, user.ID.String())
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

	post, err := h.service.post.GetByID(ctx.Context(), id, user.ID.String())
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
	_, err := GetUserFromCtx(ctx)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}

	id := ctx.Params("id")
	repost, err := h.service.repost.GetByID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	err = h.service.repost.Delete(ctx.Context(), dto.DeleteRepost{
		ID: repost.ID.String(),
	})
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, nil)
}
