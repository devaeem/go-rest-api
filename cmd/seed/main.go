package main

import (
	"errors"
	"log"
	"os"

	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/database"
	"github.com/devaeem/go-rest-api/internal/core/users"
)

func main() {
	cfg := config.LoadConfig()

	name := os.Getenv("SEED_ADMIN_NAME")
	email := os.Getenv("SEED_ADMIN_EMAIL")
	password := os.Getenv("SEED_ADMIN_PASSWORD")
	if name == "" || email == "" || len(password) < 8 {
		log.Fatal("set SEED_ADMIN_NAME, SEED_ADMIN_EMAIL and SEED_ADMIN_PASSWORD (min 8 characters)")
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	service := users.NewService(users.NewRepository(db))
	user, err := service.Create(name, email, password, users.RoleAdmin)
	if errors.Is(err, users.ErrEmailTaken) {
		log.Printf("user %s already exists, nothing to do", email)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created user %s (%s)", user.Email, user.ID)
}
