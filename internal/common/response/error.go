package response

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(c, fiberErr.Code, fiberErr.Message)
	}

	log.Printf("%s %s: %v", c.Method(), c.Path(), err)
	return Error(c, fiber.StatusInternalServerError, "internal server error")
}
