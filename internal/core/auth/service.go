package auth

import (
	"errors"

	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/password"
	"github.com/devaeem/go-rest-api/internal/common/token"
	"github.com/devaeem/go-rest-api/internal/core/users"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Service struct {
	users *users.Service
	cfg   *config.Config
}

func NewService(users *users.Service, cfg *config.Config) *Service {
	return &Service{users: users, cfg: cfg}
}

func (s *Service) Register(name, email, plainPassword string) (*users.User, error) {
	return s.users.Create(name, email, plainPassword)
}

func (s *Service) Login(email, plainPassword string) (*users.User, string, error) {
	user, err := s.users.GetByEmail(email)
	if errors.Is(err, users.ErrUserNotFound) {
		return nil, "", ErrInvalidCredentials
	}
	if err != nil {
		return nil, "", err
	}

	if !password.Check(user.Password, plainPassword) {
		return nil, "", ErrInvalidCredentials
	}

	accessToken, err := token.Generate(user.ID.String(), s.cfg.JWTSecret, s.cfg.JWTExpiresIn)
	if err != nil {
		return nil, "", err
	}
	return user, accessToken, nil
}
