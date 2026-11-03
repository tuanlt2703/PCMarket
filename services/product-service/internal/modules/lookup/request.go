package lookup

import (
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

// CreateRequest is the JSON body of POST /admin/lookups/:type.
type CreateRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	// IsActive defaults to true when omitted.
	IsActive *bool `json:"is_active"`
}

func (r CreateRequest) Validate() error {
	var v httpx.Validator
	v.RequiredString("code", r.Code, 50)
	v.RequiredString("name", r.Name, 100)
	return v.Err()
}

func (r CreateRequest) Input() CreateInput {
	return CreateInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// UpdateRequest is the JSON body of PUT /admin/lookups/:type/:id.
// It is a partial update: omitted fields are unchanged and an explicit null
// clears a nullable field.
type UpdateRequest struct {
	Code        optional.Field[string] `json:"code"`
	Name        optional.Field[string] `json:"name"`
	Description optional.Field[string] `json:"description"`
	IsActive    optional.Field[bool]   `json:"is_active"`
}

func (r UpdateRequest) Validate() error {
	var v httpx.Validator
	v.MaxLength("code", r.Code.Value, 50)
	v.MaxLength("name", r.Name.Value, 100)
	return v.Err()
}

func (r UpdateRequest) Input() UpdateInput {
	return UpdateInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// PublicItem is the storefront view of an Item, without audit fields.
type PublicItem struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func toPublicItems(items []Item) []PublicItem {
	result := make([]PublicItem, len(items))
	for i, item := range items {
		result[i] = PublicItem{ID: item.ID, Code: item.Code, Name: item.Name, Description: item.Description}
	}
	return result
}
