package main

import (
	"log"
	"strings"

	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	app := fiber.New()
	config := config.LoadConfig()
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(config.CORSAllowedOrigins, ","), // .env เป็น string คั่นด้วย , เลยต้องแยกเป็น []string
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	app.Get("/", func(c fiber.Ctx) error {
		log.Println("this is a test log")
		return c.JSON(fiber.Map{
			"message": "Hello, Fiber v3!",
		})
	})
	log.Println("Server is running on port " + config.Port)
	log.Fatal(app.Listen(":" + config.Port))

}
