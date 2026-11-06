package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect เปิดการเชื่อมต่อ PostgreSQL ผ่าน GORM (gorm.Open จะ ping ให้เองว่าต่อได้จริง)
// url ตัวอย่าง: postgres://user:pass@localhost:5432/dbname?sslmode=disable
func Connect(url string) (*gorm.DB, error) {
	if url == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return db, nil
}
