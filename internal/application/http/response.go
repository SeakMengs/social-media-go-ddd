package http

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  int         `json:"status"`
	Error   any         `json:"error,omitempty"`
}

func NewResponse(success bool, message string, data interface{}, status int, err any) *Response {
	return &Response{
		Success: success,
		Message: message,
		Data:    data,
		Status:  status,
		Error:   err,
	}
}

func SuccessResponse(ctx *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Request success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	resp := NewResponse(true, msg, data, fiber.StatusOK, nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func ErrorResponse(ctx *fiber.Ctx, status int, err any, message ...string) error {
	msg := "Request failed"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	if e, ok := err.(error); ok {
		err = e.Error()
	}

	resp := NewResponse(false, msg, nil, status, err)
	return ctx.Status(status).JSON(resp)
}
