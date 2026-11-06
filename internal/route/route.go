package route

import (
	"time"

	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/middleware"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/core/auth"
	"github.com/devaeem/go-rest-api/internal/core/users"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	app.Get("/", func(c fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, fiber.Map{"message": "Hello, Fiber v3!"})
	})

	userService := users.NewService(users.NewRepository(db))
	userHandler := users.NewHandler(userService)
	authHandler := auth.NewHandler(auth.NewService(userService, cfg), cfg)

	requireToken := middleware.Auth(cfg.JWTSecret)
	requireAdmin := middleware.RequireRole(users.RoleAdmin)
	loginLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "too many login attempts, try again later")
		},
	})

	api := app.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.Post("/login", loginLimiter, authHandler.Login)
	authRoutes.Post("/logout", authHandler.Logout)
	authRoutes.Get("/me", requireToken, authHandler.RequireLogin, authHandler.Me)

	userRoutes := api.Group("/users", requireToken, authHandler.RequireLogin)
	userRoutes.Post("/", requireAdmin, userHandler.Create)
	userRoutes.Get("/", userHandler.List)
	userRoutes.Get("/:id", userHandler.Get)
	userRoutes.Patch("/:id", userHandler.Update)
	userRoutes.Delete("/:id", userHandler.Delete)
}
