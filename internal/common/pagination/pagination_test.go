package pagination

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestFromQuery(t *testing.T) {
	tests := []struct {
		query string
		want  Params
	}{
		{"", Params{Page: 1, PageSize: 10}},
		{"?page=3&pageSize=20", Params{Page: 3, PageSize: 20}},
		{"?page=0&pageSize=-5", Params{Page: 1, PageSize: 10}},
		{"?page=abc&pageSize=xyz", Params{Page: 1, PageSize: 10}},
		{"?pageSize=1000", Params{Page: 1, PageSize: 100}},
	}

	for _, tt := range tests {
		var got Params
		app := fiber.New()
		app.Get("/", func(c fiber.Ctx) error {
			got = FromQuery(c)
			return nil
		})
		if _, err := app.Test(httptest.NewRequest("GET", "/"+tt.query, nil)); err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Errorf("query %q: got %+v, want %+v", tt.query, got, tt.want)
		}
	}
}

func TestMeta(t *testing.T) {
	p := Params{Page: 2, PageSize: 10}
	if p.Offset() != 10 {
		t.Errorf("offset: got %d, want 10", p.Offset())
	}
	if got := p.Meta(21).TotalPages; got != 3 {
		t.Errorf("totalPages for 21 items: got %d, want 3", got)
	}
	if got := p.Meta(0).TotalPages; got != 0 {
		t.Errorf("totalPages for 0 items: got %d, want 0", got)
	}
}
