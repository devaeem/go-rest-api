package middleware

import (
	"slices"

	"github.com/devaeem/go-rest-api/internal/common/token"
	"github.com/gofiber/fiber/v3"
)

const (
	AuthCookieName = "access_token"
	userIDKey      = "userID"
	roleKey        = "role"
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

func UserID(c fiber.Ctx) string {
	userID, _ := c.Locals(userIDKey).(string)
	return userID
}

func SetRole(c fiber.Ctx, role string) {
	c.Locals(roleKey, role)
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role, _ := c.Locals(roleKey).(string)
		if slices.Contains(roles, role) {
			return c.Next()
		}
		return fiber.NewError(fiber.StatusForbidden, "you do not have permission to do this")
	}
}
