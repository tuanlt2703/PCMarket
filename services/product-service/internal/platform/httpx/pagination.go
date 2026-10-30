package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type Page struct {
	Page     int
	PageSize int
}

func (p Page) Offset() int { return (p.Page - 1) * p.PageSize }

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func NewPageResult[T any](items []T, p Page, total int64) PageResult[T] {
	if items == nil {
		items = []T{}
	}
	return PageResult[T]{Items: items, Page: p.Page, PageSize: p.PageSize, Total: total}
}

// ParsePage reads ?page=&page_size= with safe defaults and bounds.
func ParsePage(c *gin.Context) Page {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || size < 1 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return Page{Page: page, PageSize: size}
}

// QueryBool parses an optional boolean query parameter, writing a 400 response on failure.
func QueryBool(c *gin.Context, name string) (*bool, bool) {
	raw, ok := c.GetQuery(name)
	if !ok || raw == "" {
		return nil, true
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		BadRequest(c, "invalid "+name)
		return nil, false
	}
	return &v, true
}

// QueryInt64 parses an optional positive int64 query parameter, writing a 400 response on failure.
func QueryInt64(c *gin.Context, name string) (*int64, bool) {
	raw, ok := c.GetQuery(name)
	if !ok || raw == "" {
		return nil, true
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v <= 0 {
		BadRequest(c, "invalid "+name)
		return nil, false
	}
	return &v, true
}
