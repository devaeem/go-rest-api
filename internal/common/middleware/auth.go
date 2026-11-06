package middleware

import (
	"github.com/devaeem/go-rest-api/internal/common/token"
	"github.com/gofiber/fiber/v3"
)

const (
	AuthCookieName = "access_token"
	userIDKey      = "userID"
)

func Auth(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies(AuthCookieName)
		if accessToken == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "not logged in")
		}

		userID, err := token.Parse(accessToken, secret)
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
