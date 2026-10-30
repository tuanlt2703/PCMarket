package catalog

import (
	"time"

	"product-service/internal/platform/optional"
)

type Manufacturer struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateManufacturerInput struct {
	Code        string
	Name        string
	Description *string
	// IsActive defaults to true when omitted.
	IsActive *bool
}

// UpdateManufacturerInput is a partial update: only present fields are changed.
type UpdateManufacturerInput struct {
	Code        optional.Field[string]
	Name        optional.Field[string]
	Description optional.Field[string]
	IsActive    optional.Field[bool]
}

type ManufacturerFilter struct {
	Query    string
	IsActive *bool
}
