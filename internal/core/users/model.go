package users

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	// uuidv7() มีใน postgres 18 เรียงตามเวลา ทำให้ index เร็วกว่า uuid แบบสุ่ม
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // เก็บเป็น hash เท่านั้น และไม่ส่งออกใน JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // soft delete: ลบแล้วแค่ซ่อน ไม่ได้ลบจริง
}
