package catalog

import (
	"time"

	"product-service/internal/platform/optional"
)

type ComponentType struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateComponentTypeInput struct {
	// Code links a type to its specification module, e.g. "cpu".
	Code        string
	Name        string
	Description *string
	// IsActive defaults to true when omitted.
	IsActive *bool
}

// UpdateComponentTypeInput has no Code or Name: both are fixed after creation
// because specification modules and clients identify the type by them.
type UpdateComponentTypeInput struct {
	Description optional.Field[string]
	IsActive    optional.Field[bool]
}

type ComponentTypeFilter struct {
	Query    string
	IsActive *bool
}
