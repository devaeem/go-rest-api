package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/middleware"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/common/validate"
	"github.com/devaeem/go-rest-api/internal/core/users"
	"github.com/gofiber/fiber/v3"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Handler struct {
	service *Service
	cfg     *config.Config
}

func NewHandler(service *Service, cfg *config.Config) *Handler {
	return &Handler{service: service, cfg: cfg}
}

func (h *Handler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name, email and password are required")
	}
	if !validate.Email(req.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email")
	}
	if len(req.Password) < 8 || len(req.Password) > 72 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be 8-72 characters")
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password)
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

	user, accessToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return toHTTPError(err)
	}

	h.setAuthCookie(c, accessToken, time.Now().Add(h.cfg.JWTExpiresIn))
	return response.Success(c, fiber.StatusOK, user)
}

func (h *Handler) Logout(c fiber.Ctx) error {
	h.setAuthCookie(c, "", time.Unix(0, 0))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) setAuthCookie(c fiber.Ctx, value string, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     middleware.AuthCookieName,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func toHTTPError(err error) error {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, users.ErrEmailTaken):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	default:
		return err
	}
}
