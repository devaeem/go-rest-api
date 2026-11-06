package middleware

import (
	"strings"

	"github.com/devaeem/go-rest-api/internal/common/pkg"
	"github.com/gofiber/fiber/v3"
)

const userIDKey = "userID"

// Auth ตรวจ header "Authorization: Bearer <token>" ถ้าผ่านจะเก็บ userID ไว้ให้ handler ใช้
func Auth(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, ok := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
		if !ok || token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}

		userID, err := pkg.ParseToken(token, secret)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals(userIDKey, userID)
		return c.Next()
	}
}

// UserID ดึง userID ที่ Auth เก็บไว้ ใช้ใน handler ที่อยู่หลัง Auth
func UserID(c fiber.Ctx) string {
	userID, _ := c.Locals(userIDKey).(string)
	return userID
}
