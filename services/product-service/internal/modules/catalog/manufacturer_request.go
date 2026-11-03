package catalog

import (
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

// CreateManufacturerRequest is the JSON body of POST /admin/manufacturers.
type CreateManufacturerRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	// IsActive defaults to true when omitted.
	IsActive *bool `json:"is_active"`
}

func (r CreateManufacturerRequest) Validate() error {
	var v httpx.Validator
	v.RequiredString("code", r.Code, 50)
	v.RequiredString("name", r.Name, 100)
	return v.Err()
}

func (r CreateManufacturerRequest) Input() CreateManufacturerInput {
	return CreateManufacturerInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// UpdateManufacturerRequest is the JSON body of PUT /admin/manufacturers/:id.
// It is a partial update: omitted fields are unchanged and an explicit null
// clears a nullable field.
type UpdateManufacturerRequest struct {
	Code        optional.Field[string] `json:"code"`
	Name        optional.Field[string] `json:"name"`
	Description optional.Field[string] `json:"description"`
	IsActive    optional.Field[bool]   `json:"is_active"`
}

func (r UpdateManufacturerRequest) Validate() error {
	var v httpx.Validator
	v.MaxLength("code", r.Code.Value, 50)
	v.MaxLength("name", r.Name.Value, 100)
	return v.Err()
}

func (r UpdateManufacturerRequest) Input() UpdateManufacturerInput {
	return UpdateManufacturerInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}
