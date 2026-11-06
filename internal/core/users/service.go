package users

import (
	"errors"
	"strings"

	"github.com/devaeem/go-rest-api/internal/common/pagination"
	"github.com/devaeem/go-rest-api/internal/common/password"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email already in use")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, email, plainPassword string) (*User, error) {
	hash, err := password.Hash(plainPassword)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     strings.TrimSpace(name),
		Email:    normalizeEmail(email),
		Password: hash,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, translateError(err)
	}
	return user, nil
}

func (s *Service) List(p pagination.Params) ([]User, pagination.Meta, error) {
	users, total, err := s.repo.FindAll(p.Offset(), p.PageSize)
	if err != nil {
		return nil, pagination.Meta{}, err
	}
	return users, p.Meta(total), nil
}

func (s *Service) Get(id uuid.UUID) (*User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, translateError(err)
	}
	return user, nil
}

func (s *Service) Update(id uuid.UUID, name, email *string) (*User, error) {
	user, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		user.Name = strings.TrimSpace(*name)
	}
	if email != nil {
		user.Email = normalizeEmail(*email)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, translateError(err)
	}
	return user, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	user, err := s.Get(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(user)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func translateError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrUserNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrEmailTaken
	default:
		return err
	}
}
