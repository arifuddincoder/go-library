package query

import (
	"strings"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Params struct {
	Page   int
	Limit  int
	Search string
	SortBy string
	Order  string
}

// Echo context theke query param parse kore
func Parse(c *echo.Context) Params {
	page, _ := echo.QueryParam[int](c, "page")
	limit, _ := echo.QueryParam[int](c, "limit")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	order := strings.ToLower(c.QueryParam("order"))
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return Params{
		Page:   page,
		Limit:  limit,
		Search: c.QueryParam("search"),
		SortBy: c.QueryParam("sort_by"),
		Order:  order,
	}
}

// Pagination scope
func Paginate(p Params) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	}
}

// Search scope — kon kon column e search hobe seita pass korbe
func Search(keyword string, columns ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(columns) == 0 {
			return db
		}
		q := db
		like := "%" + keyword + "%"
		for i, col := range columns {
			if i == 0 {
				q = q.Where(col+" ILIKE ?", like)
			} else {
				q = q.Or(col+" ILIKE ?", like)
			}
		}
		return db.Where(q)
	}
}

// Sort scope
func Sort(p Params, allowed map[string]bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.SortBy != "" && allowed[p.SortBy] {
			return db.Order(p.SortBy + " " + p.Order)
		}
		return db.Order("created_at " + p.Order)
	}
}
