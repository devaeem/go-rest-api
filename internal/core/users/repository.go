package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindAll(offset, limit int) ([]User, int64, error) {
	var total int64
	if err := r.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []User
	err := r.db.Order("created_at, id").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *Repository) FindByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Update(user *User, changes map[string]any) error {
	return r.db.Model(user).Updates(changes).Error
}

func (r *Repository) Delete(user *User) error {
	return r.db.Delete(user).Error
}
