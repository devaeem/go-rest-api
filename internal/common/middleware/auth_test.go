package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devaeem/go-rest-api/internal/common/token"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

const secret = "test-secret"

func TestAuth(t *testing.T) {
	app := fiber.New()
	app.Get("/me", Auth(secret), func(c fiber.Ctx) error {
		return c.SendString(UserID(c))
	})

	valid, _ := token.Generate("user-123", secret, time.Hour)
	wrongSecret, _ := token.Generate("user-123", "other-secret", time.Hour)
	expired, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "user-123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
	}).SignedString([]byte(secret))
	algNone, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject: "user-123",
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)

	tests := []struct {
		name       string
		cookie     string
		wantStatus int
	}{
		{"valid token", valid, fiber.StatusOK},
		{"no cookie", "", fiber.StatusUnauthorized},
		{"garbage", "not-a-jwt", fiber.StatusUnauthorized},
		{"wrong secret", wrongSecret, fiber.StatusUnauthorized},
		{"expired", expired, fiber.StatusUnauthorized},
		{"alg none", algNone, fiber.StatusUnauthorized},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/me", nil)
		if tt.cookie != "" {
			req.AddCookie(&http.Cookie{Name: AuthCookieName, Value: tt.cookie})
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
	if _, err := token.Generate("user-123", "", time.Hour); err != token.ErrEmptySecret {
		t.Errorf("GenerateToken with empty secret: got %v, want ErrEmptySecret", err)
	}
}
