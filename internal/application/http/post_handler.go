package http

import (
	"social-media-go-ddd/internal/application/service"
	"social-media-go-ddd/internal/domain/dto"

	"github.com/gofiber/fiber/v2"
)

type PostHandlerService struct {
	post *service.PostService
}

func NewPostHandlerService(post *service.PostService) *PostHandlerService {
	return &PostHandlerService{
		post: post,
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

func NewPostHandler(postService *service.PostService, authMiddleware *AuthMiddleware) *PostHandler {
	return &PostHandler{
		service:    NewPostHandlerService(postService),
		middleware: NewPostHandlerMiddleware(authMiddleware),
	}
}

func (h *PostHandler) RegisterRoutes(app *fiber.App) {
	apiPosts := app.Group("/api/v1/posts")
	apiPosts.Get("/:id", h.GetPostByID)

	apiPostsProtected := app.Group("/api/v1/posts", h.middleware.auth.Handler)
	apiPostsProtected.Post("/", h.CreatePost)

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
	body.NewPost.UserID = user.ID
	if err := ctx.BodyParser(&body); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	post, err := h.service.post.Create(ctx.Context(), body.NewPost)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) GetPostByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	post, err := h.service.post.GetByID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, err)
	}

	return SuccessResponse(ctx, fiber.Map{
		"post": post,
	})
}
