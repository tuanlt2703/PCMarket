package lookup

import (
	"time"

	"product-service/internal/platform/optional"
)

// Table is a lookup table. All lookup tables share the shape
// (id, code, name, description, is_active, created_at, updated_at).
// Names come only from the registry, never from user input.
type Table struct {
	// Name is the SQL table name.
	Name string
	// Label is a human-readable singular name used in error messages.
	Label string
}

type Item struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Input struct {
	Code        string
	Name        string
	Description *string
	// IsActive defaults to true when omitted.
	IsActive *bool
}

// UpdateInput is a partial update: only present fields are changed.
type UpdateInput struct {
	Code        optional.Field[string]
	Name        optional.Field[string]
	Description optional.Field[string]
	IsActive    optional.Field[bool]
}

type Filter struct {
	Query    string
	IsActive *bool
}
