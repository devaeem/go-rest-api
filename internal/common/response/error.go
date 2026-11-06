package response

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

// ErrorHandler ทำให้ทุก error ตอบเป็น BaseResponse รูปแบบเดียวกัน
func ErrorHandler(c fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(c, fiberErr.Code, fiberErr.Message)
	}

	// error ที่ไม่ได้ตั้งใจส่งให้ client: log ไว้ แล้วตอบแค่ 500 ไม่บอกรายละเอียด
	log.Printf("%s %s: %v", c.Method(), c.Path(), err)
	return Error(c, fiber.StatusInternalServerError, "internal server error")
}
