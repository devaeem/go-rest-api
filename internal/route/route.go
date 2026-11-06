package route

import (
	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/middleware"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/core/users"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	app.Get("/", func(c fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, fiber.Map{"message": "Hello, Fiber v3!"})
	})

	auth := middleware.Auth(cfg.JWTSecret)
	userHandler := users.NewHandler(users.NewService(users.NewRepository(db)), cfg.JWTSecret, cfg.JWTExpiresIn)
	api := app.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.Post("/register", userHandler.Register)
	authRoutes.Post("/login", userHandler.Login)

	userRoutes := api.Group("/users", auth)
	userRoutes.Get("/", userHandler.List)
	userRoutes.Get("/:id", userHandler.Get)
	userRoutes.Patch("/:id", userHandler.Update)
	userRoutes.Delete("/:id", userHandler.Delete)
}
