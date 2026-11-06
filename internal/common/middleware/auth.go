package middleware

import (
	"strings"

	"github.com/devaeem/go-rest-api/internal/common/token"
	"github.com/gofiber/fiber/v3"
)

const userIDKey = "userID"

func Auth(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		bearer, ok := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
		if !ok || bearer == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}

		userID, err := token.Parse(bearer, secret)
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
