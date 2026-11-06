package main

import (
	"log"

	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/database"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/route"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters (generate one with: openssl rand -base64 32)")
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database successfully")

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{ErrorHandler: response.ErrorHandler})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	route.Setup(app, db, cfg)

	log.Println("Server is running on port " + cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
