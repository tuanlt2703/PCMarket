package catalog

import "product-service/internal/platform/optional"

type CreateManufacturerRequest struct {
	Code        string  `json:"code" binding:"required,max=50"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description"`
	// IsActive defaults to true when omitted.
	IsActive *bool `json:"is_active"`
}

func (r CreateManufacturerRequest) Input() CreateManufacturerInput {
	return CreateManufacturerInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// UpdateManufacturerRequest is a partial update: omitted fields are unchanged
// and an explicit null clears a nullable field.
type UpdateManufacturerRequest struct {
	Code        optional.Field[string] `json:"code" binding:"omitempty,max=50"`
	Name        optional.Field[string] `json:"name" binding:"omitempty,max=100"`
	Description optional.Field[string] `json:"description"`
	IsActive    optional.Field[bool]   `json:"is_active"`
}

func (r UpdateManufacturerRequest) Input() UpdateManufacturerInput {
	return UpdateManufacturerInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}
