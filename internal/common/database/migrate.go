package database

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	return nil
}
