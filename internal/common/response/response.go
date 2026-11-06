package response

import "github.com/gofiber/fiber/v3"

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func Success(c fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(BaseResponse{Success: true, Data: data})
}

func Paginated(c fiber.Ctx, data any, meta any) error {
	return c.Status(fiber.StatusOK).JSON(BaseResponse{Success: true, Data: data, Meta: meta})
}

func Error(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(BaseResponse{Success: false, Message: message})
}
