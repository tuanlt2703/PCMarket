package catalog

import (
	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

// CreateComponentTypeRequest is the JSON body of POST /admin/component-types.
type CreateComponentTypeRequest struct {
	// Code links a type to its specification module, e.g. "cpu".
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	// IsActive defaults to true when omitted.
	IsActive *bool `json:"is_active"`
}

func (r CreateComponentTypeRequest) Validate() error {
	var v httpx.Validator
	v.RequiredString("code", r.Code, 50)
	v.RequiredString("name", r.Name, 100)
	return v.Err()
}

func (r CreateComponentTypeRequest) Input() CreateComponentTypeInput {
	return CreateComponentTypeInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// UpdateComponentTypeRequest is the JSON body of PUT /admin/component-types/:id.
// It is a partial update: omitted fields are unchanged.
type UpdateComponentTypeRequest struct {
	Description optional.Field[string] `json:"description"`
	IsActive    optional.Field[bool]   `json:"is_active"`

	// Code and Name cannot be changed. They are decoded only so the client
	// gets a clear 422 message instead of a generic "unknown field" 400.
	Code optional.Field[any] `json:"code"`
	Name optional.Field[any] `json:"name"`
}

func (r UpdateComponentTypeRequest) Validate() error {
	if r.Code.Present || r.Name.Present {
		return apperror.Invalid("code and name cannot be changed after creation")
	}
	return nil
}

func (r UpdateComponentTypeRequest) Input() UpdateComponentTypeInput {
	return UpdateComponentTypeInput{
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}
