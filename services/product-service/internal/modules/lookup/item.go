package lookup

import (
	"time"

	"product-service/internal/platform/optional"
)

// Table describes one lookup table, e.g. cpu_sockets.
// Every lookup table has the same columns:
// (id, code, name, description, is_active, created_at, updated_at).
// Tables only come from the registry, never from user input.
type Table struct {
	// Name is the SQL table name, e.g. "cpu_sockets".
	Name string
	// Label is a human-readable singular name used in error messages, e.g. "cpu sockets".
	Label string
}

// Item is one row of a lookup table.
type Item struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInput struct {
	Code        string
	Name        string
	Description *string
	// IsActive defaults to true when nil.
	IsActive *bool
}

// UpdateInput is a partial update: only present fields are changed.
type UpdateInput struct {
	Code        optional.Field[string]
	Name        optional.Field[string]
	Description optional.Field[string]
	IsActive    optional.Field[bool]
}

// Filter narrows an admin list. Empty Query and nil IsActive mean "no filter".
type Filter struct {
	Query    string
	IsActive *bool
}
