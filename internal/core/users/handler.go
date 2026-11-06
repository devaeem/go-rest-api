package users

import (
	"errors"
	"net/mail"

	"github.com/devaeem/go-rest-api/internal/common/pagination"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
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
	id, err := parseID(c)
	if err != nil {
		return err
	}

	var req UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name != nil && *req.Name == "" {
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
	id, err := parseID(c)
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

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// toHTTPError แปลง error จาก service เป็น status code ที่เหมาะสม
func toHTTPError(err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, ErrEmailTaken):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	default:
		return err // response.ErrorHandler จะ log และตอบ 500 ให้
	}
}
