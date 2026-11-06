package users

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:255;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"size:20;not null;default:user" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
