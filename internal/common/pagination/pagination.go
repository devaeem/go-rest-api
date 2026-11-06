package pagination

import "github.com/gofiber/fiber/v3"

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

type Params struct {
	Page     int
	PageSize int
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// FromQuery อ่าน ?page=&pageSize= ค่าที่ผิดหรือเกินขอบเขตจะถูกแทนด้วยค่า default
func FromQuery(c fiber.Ctx) Params {
	page := fiber.Query(c, "page", defaultPage)
	pageSize := fiber.Query(c, "pageSize", defaultPageSize)

	if page < 1 {
		page = defaultPage
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return Params{Page: page, PageSize: pageSize}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Params) Meta(total int64) Meta {
	totalPages := int((total + int64(p.PageSize) - 1) / int64(p.PageSize))
	return Meta{Page: p.Page, PageSize: p.PageSize, Total: total, TotalPages: totalPages}
}
