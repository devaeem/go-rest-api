package route

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/devaeem/go-rest-api/internal/common/config"
	"github.com/devaeem/go-rest-api/internal/common/database"
	"github.com/devaeem/go-rest-api/internal/common/middleware"
	"github.com/devaeem/go-rest-api/internal/common/response"
	"github.com/devaeem/go-rest-api/internal/core/users"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

const testPassword = "secret123"

type testApp struct {
	t   *testing.T
	app *fiber.App
	db  *gorm.DB
}

type apiResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// newTestApp needs TEST_DATABASE_URL pointing at a database whose name ends in _test, because every test wipes the users table.
func newTestApp(t *testing.T) *testApp {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	if !strings.Contains(url, "_test") {
		t.Fatal("TEST_DATABASE_URL must point at a *_test database")
	}

	db, err := database.Connect(url)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("TRUNCATE users").Error; err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		JWTSecret:    "test-secret-that-is-at-least-32-chars",
		JWTExpiresIn: time.Hour,
	}
	app := fiber.New(fiber.Config{ErrorHandler: response.ErrorHandler})
	Setup(app, db, cfg)
	return &testApp{t: t, app: app, db: db}
}

func (a *testApp) createUser(email, role string) users.User {
	a.t.Helper()
	user, err := users.NewService(users.NewRepository(a.db)).Create("Test", email, testPassword, role)
	if err != nil {
		a.t.Fatal(err)
	}
	return *user
}

func (a *testApp) login(email string) *http.Cookie {
	a.t.Helper()
	resp, _ := a.request("POST", "/api/v1/auth/login", `{"email":"`+email+`","password":"`+testPassword+`"}`, nil)
	if resp.StatusCode != http.StatusOK {
		a.t.Fatalf("login %s: got %d", email, resp.StatusCode)
	}
	for _, c := range resp.Cookies() {
		if c.Name == middleware.AuthCookieName {
			return c
		}
	}
	a.t.Fatal("login did not set the auth cookie")
	return nil
}

func (a *testApp) request(method, path, body string, cookie *http.Cookie) (*http.Response, apiResponse) {
	a.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}

	resp, err := a.app.Test(req)
	if err != nil {
		a.t.Fatal(err)
	}

	var parsed apiResponse
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	return resp, parsed
}

func (a *testApp) expectStatus(name string, resp *http.Response, want int) {
	a.t.Helper()
	if resp.StatusCode != want {
		a.t.Errorf("%s: got %d, want %d", name, resp.StatusCode, want)
	}
}

func TestLogin(t *testing.T) {
	a := newTestApp(t)
	a.createUser("ann@example.com", users.RoleUser)

	resp, _ := a.request("POST", "/api/v1/auth/login", `{"email":"ann@example.com","password":"wrong-password"}`, nil)
	a.expectStatus("wrong password", resp, http.StatusUnauthorized)

	resp, _ = a.request("POST", "/api/v1/auth/login", `{"email":"nobody@example.com","password":"secret123"}`, nil)
	a.expectStatus("unknown email", resp, http.StatusUnauthorized)

	cookie := a.login("ANN@example.com")
	if !cookie.HttpOnly {
		t.Error("auth cookie must be HttpOnly")
	}

	resp, body := a.request("GET", "/api/v1/auth/me", "", cookie)
	a.expectStatus("me", resp, http.StatusOK)
	if !strings.Contains(string(body.Data), "ann@example.com") {
		t.Errorf("me returned %s", body.Data)
	}

	resp, _ = a.request("GET", "/api/v1/auth/me", "", nil)
	a.expectStatus("me without cookie", resp, http.StatusUnauthorized)
}

func TestLoginRateLimit(t *testing.T) {
	a := newTestApp(t)

	var last *http.Response
	for range 6 {
		last, _ = a.request("POST", "/api/v1/auth/login", `{"email":"x@example.com","password":"wrong-password"}`, nil)
	}
	a.expectStatus("6th login attempt in a minute", last, http.StatusTooManyRequests)
}

func TestOnlyAdminCanCreateUsers(t *testing.T) {
	a := newTestApp(t)
	a.createUser("admin@example.com", users.RoleAdmin)
	a.createUser("ann@example.com", users.RoleUser)
	newUser := `{"name":"Bob","email":"bob@example.com","password":"secret123"}`

	resp, _ := a.request("POST", "/api/v1/users", newUser, a.login("ann@example.com"))
	a.expectStatus("normal user creates user", resp, http.StatusForbidden)

	admin := a.login("admin@example.com")
	resp, body := a.request("POST", "/api/v1/users", newUser, admin)
	a.expectStatus("admin creates user", resp, http.StatusCreated)
	if !strings.Contains(string(body.Data), `"role":"user"`) {
		t.Errorf("new user should default to role user, got %s", body.Data)
	}

	resp, _ = a.request("POST", "/api/v1/users", newUser, admin)
	a.expectStatus("duplicate email", resp, http.StatusConflict)

	resp, _ = a.request("POST", "/api/v1/users", `{"name":"C","email":"c@example.com","password":"secret123","role":"root"}`, admin)
	a.expectStatus("unknown role", resp, http.StatusBadRequest)
}

func TestDeletedUserTokenIsRejected(t *testing.T) {
	a := newTestApp(t)
	ann := a.createUser("ann@example.com", users.RoleUser)
	cookie := a.login("ann@example.com")

	resp, _ := a.request("DELETE", "/api/v1/users/"+ann.ID.String(), "", cookie)
	a.expectStatus("delete own account", resp, http.StatusNoContent)

	resp, _ = a.request("GET", "/api/v1/users", "", cookie)
	a.expectStatus("list with deleted user's token", resp, http.StatusUnauthorized)
}

func TestDeletedEmailCanBeReused(t *testing.T) {
	a := newTestApp(t)
	ann := a.createUser("ann@example.com", users.RoleUser)

	if err := a.db.Delete(&ann).Error; err != nil {
		t.Fatal(err)
	}
	a.createUser("ann@example.com", users.RoleUser)
}

func TestUpdateUser(t *testing.T) {
	a := newTestApp(t)
	ann := a.createUser("ann@example.com", users.RoleUser)
	bob := a.createUser("bob@example.com", users.RoleUser)
	cookie := a.login("ann@example.com")
	annPath := "/api/v1/users/" + ann.ID.String()

	resp, body := a.request("PATCH", annPath, `{"name":"  Ann B  "}`, cookie)
	a.expectStatus("update own name", resp, http.StatusOK)
	if !strings.Contains(string(body.Data), `"name":"Ann B"`) {
		t.Errorf("name should be trimmed, got %s", body.Data)
	}

	resp, _ = a.request("PATCH", annPath, `{"name":"   "}`, cookie)
	a.expectStatus("blank name", resp, http.StatusBadRequest)

	resp, _ = a.request("PATCH", annPath, `{"email":"bob@example.com"}`, cookie)
	a.expectStatus("email taken", resp, http.StatusConflict)

	resp, _ = a.request("PATCH", "/api/v1/users/"+bob.ID.String(), `{"name":"hacked"}`, cookie)
	a.expectStatus("update someone else", resp, http.StatusForbidden)

	var stored users.User
	a.db.First(&stored, "id = ?", ann.ID)
	if stored.Password == "" || stored.Role != users.RoleUser {
		t.Errorf("update must not touch other columns: %+v", stored)
	}
}

func TestListUsersPagination(t *testing.T) {
	a := newTestApp(t)
	for _, email := range []string{"u1@example.com", "u2@example.com", "u3@example.com"} {
		a.createUser(email, users.RoleUser)
	}
	cookie := a.login("u1@example.com")

	resp, body := a.request("GET", "/api/v1/users?page=2&pageSize=2", "", cookie)
	a.expectStatus("list page 2", resp, http.StatusOK)

	var list []users.User
	if err := json.Unmarshal(body.Data, &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Email != "u3@example.com" {
		t.Errorf("page 2 of 3 users with pageSize 2 should hold only u3, got %+v", list)
	}
	if strings.Contains(string(body.Data), "password") {
		t.Error("password hash must never be returned")
	}
}
