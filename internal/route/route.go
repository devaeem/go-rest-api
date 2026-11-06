package route

import (
	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/middleware"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/core/auth"
	"github.com/devaeem/go-rest-api/internal/core/users"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	app.Get("/", func(c fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, fiber.Map{"message": "Hello, Fiber v3!"})
	})

	requireAuth := middleware.Auth(cfg.JWTSecret)
	userService := users.NewService(users.NewRepository(db))
	userHandler := users.NewHandler(userService)
	authHandler := auth.NewHandler(auth.NewService(userService, cfg), cfg)
	api := app.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.Post("/register", authHandler.Register)
	authRoutes.Post("/login", authHandler.Login)
	authRoutes.Post("/logout", authHandler.Logout)

	userRoutes := api.Group("/users", requireAuth)
	userRoutes.Get("/", userHandler.List)
	userRoutes.Get("/:id", userHandler.Get)
	userRoutes.Patch("/:id", userHandler.Update)
	userRoutes.Delete("/:id", userHandler.Delete)
}
