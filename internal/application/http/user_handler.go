package http

import (
	"social-media-go-ddd/internal/application/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/users")
	api.Post("/", h.CreateUser)
	api.Get("/:id", h.GetUserByID)
}

func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.Create(ctx.Context(), body.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(user)
}

func (h *UserHandler) GetUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := h.service.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(user)
}
