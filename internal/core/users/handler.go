package users

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/devaeem/go-rest-api/internal/common/middleware"
	"github.com/devaeem/go-rest-api/internal/common/pagination"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/common/token"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type Handler struct {
	service      *Service
	jwtSecret    string
	jwtExpiresIn time.Duration
}

func NewHandler(service *Service, jwtSecret string, jwtExpiresIn time.Duration) *Handler {
	return &Handler{service: service, jwtSecret: jwtSecret, jwtExpiresIn: jwtExpiresIn}
}

func (h *Handler) Register(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name, email and password are required")
	}
	if !isValidEmail(req.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email")
	}
	if len(req.Password) < 8 || len(req.Password) > 72 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be 8-72 characters")
	}

	user, err := h.service.Create(req.Name, req.Email, req.Password)
	if err != nil {
		return toHTTPError(err)
	}
	return response.Success(c, fiber.StatusCreated, user)
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return toHTTPError(err)
	}

	accessToken, err := token.Generate(user.ID.String(), h.jwtSecret, h.jwtExpiresIn)
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, LoginResponse{Token: accessToken, User: user})
}

func (h *Handler) List(c fiber.Ctx) error {
	users, meta, err := h.service.List(pagination.FromQuery(c))
	if err != nil {
		return toHTTPError(err)
	}
	return response.Paginated(c, users, meta)
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	user, err := h.service.Get(id)
	if err != nil {
		return toHTTPError(err)
	}
	return response.Success(c, fiber.StatusOK, user)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := parseOwnID(c)
	if err != nil {
		return err
	}

	var req UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
	}
	if req.Email != nil && !isValidEmail(*req.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email")
	}

	user, err := h.service.Update(id, req.Name, req.Email)
	if err != nil {
		return toHTTPError(err)
	}
	return response.Success(c, fiber.StatusOK, user)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := parseOwnID(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(id); err != nil {
		return toHTTPError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseID(c fiber.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}
	return id, nil
}

func parseOwnID(c fiber.Ctx) (uuid.UUID, error) {
	id, err := parseID(c)
	if err != nil {
		return uuid.Nil, err
	}
	if id.String() != middleware.UserID(c) {
		return uuid.Nil, fiber.NewError(fiber.StatusForbidden, "you can only modify your own account")
	}
	return id, nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func toHTTPError(err error) error {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrUserNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, ErrEmailTaken):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	default:
		return err // response.ErrorHandler จะ log และตอบ 500 ให้
	}
}
