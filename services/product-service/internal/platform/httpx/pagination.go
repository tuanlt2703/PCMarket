package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Page is the requested page; Number starts at 1.
type Page struct {
	Number int
	Size   int
}

// Offset is the number of rows to skip, for SQL "LIMIT size OFFSET offset".
func (p Page) Offset() int {
	return (p.Number - 1) * p.Size
}

// PageResult is the JSON body of every list endpoint.
type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func NewPageResult[T any](items []T, page Page, total int64) PageResult[T] {
	if items == nil {
		items = []T{} // send [] instead of null
	}
	return PageResult[T]{Items: items, Page: page.Number, PageSize: page.Size, Total: total}
}

// ParsePage reads ?page=&page_size=. Missing or invalid values fall back to
// page 1 and the default size; the size is capped at maxPageSize.
func ParsePage(c *gin.Context) Page {
	number, err := strconv.Atoi(c.Query("page"))
	if err != nil || number < 1 {
		number = 1
	}

	size, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || size < 1 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}

	return Page{Number: number, Size: size}
}
