package middleware

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devaeem/go-rest-api/internal/common/pkg"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

const secret = "test-secret"

func TestAuth(t *testing.T) {
	app := fiber.New()
	app.Get("/me", Auth(secret), func(c fiber.Ctx) error {
		return c.SendString(UserID(c))
	})

	valid, _ := pkg.GenerateToken("user-123", secret, time.Hour)
	wrongSecret, _ := pkg.GenerateToken("user-123", "other-secret", time.Hour)
	expired, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "user-123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
	}).SignedString([]byte(secret))
	algNone, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject: "user-123",
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{"valid token", "Bearer " + valid, fiber.StatusOK},
		{"no header", "", fiber.StatusUnauthorized},
		{"no Bearer prefix", valid, fiber.StatusUnauthorized},
		{"wrong secret", "Bearer " + wrongSecret, fiber.StatusUnauthorized},
		{"expired", "Bearer " + expired, fiber.StatusUnauthorized},
		{"alg none", "Bearer " + algNone, fiber.StatusUnauthorized},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/me", nil)
		if tt.header != "" {
			req.Header.Set("Authorization", tt.header)
		}
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != tt.wantStatus {
			t.Errorf("%s: got %d, want %d", tt.name, resp.StatusCode, tt.wantStatus)
		}
		if tt.wantStatus == fiber.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			if string(body) != "user-123" {
				t.Errorf("%s: userID = %q, want user-123", tt.name, body)
			}
		}
	}
}

func TestEmptySecret(t *testing.T) {
	if _, err := pkg.GenerateToken("user-123", "", time.Hour); err != pkg.ErrEmptySecret {
		t.Errorf("GenerateToken with empty secret: got %v, want ErrEmptySecret", err)
	}
}
