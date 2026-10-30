package catalog

import (
	"product-service/internal/platform/apperror"
	"product-service/internal/platform/optional"
)

type CreateComponentTypeRequest struct {
	// Code links a type to its specification module, e.g. "cpu".
	Code        string  `json:"code" binding:"required,max=50"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description"`
	// IsActive defaults to true when omitted.
	IsActive *bool `json:"is_active"`
}

func (r CreateComponentTypeRequest) Input() CreateComponentTypeInput {
	return CreateComponentTypeInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// UpdateComponentTypeRequest is a partial update: omitted fields are unchanged.
type UpdateComponentTypeRequest struct {
	Description optional.Field[string] `json:"description"`
	IsActive    optional.Field[bool]   `json:"is_active"`

	// Code and Name are decoded only to reject them with a clear message.
	Code optional.Field[any] `json:"code"`
	Name optional.Field[any] `json:"name"`
}

func (r UpdateComponentTypeRequest) Validate() error {
	if r.Code.Set || r.Name.Set {
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
