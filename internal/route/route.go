package route

import (
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/core/users"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Setup รวม route ของทุกฟีเจอร์ไว้ที่เดียว ฟีเจอร์ใหม่ให้มาเพิ่มที่นี่
func Setup(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, fiber.Map{"message": "Hello, Fiber v3!"})
	})

	userHandler := users.NewHandler(users.NewService(users.NewRepository(db)))
	api := app.Group("/api/v1")

	userRoutes := api.Group("/users")
	userRoutes.Post("/", userHandler.Create)
	userRoutes.Get("/", userHandler.List)
	userRoutes.Get("/:id", userHandler.Get)
	userRoutes.Patch("/:id", userHandler.Update)
	userRoutes.Delete("/:id", userHandler.Delete)
}
