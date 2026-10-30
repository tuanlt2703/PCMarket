package catalog

import "product-service/internal/platform/optional"

type CreateProductModelRequest struct {
	ComponentTypeID int64   `json:"component_type_id" binding:"required,gt=0"`
	ManufacturerID  int64   `json:"manufacturer_id" binding:"required,gt=0"`
	ModelName       string  `json:"model_name" binding:"required,max=200"`
	Description     *string `json:"description"`
	ReleaseDate     *string `json:"release_date" binding:"omitempty,datetime=2006-01-02"`
	IsActive        *bool   `json:"is_active"`
}

func (r CreateProductModelRequest) Input() CreateProductModelInput {
	return CreateProductModelInput{
		ComponentTypeID: r.ComponentTypeID,
		ManufacturerID:  r.ManufacturerID,
		ModelName:       r.ModelName,
		Description:     r.Description,
		ReleaseDate:     r.ReleaseDate,
		IsActive:        r.IsActive,
	}
}

// UpdateProductModelRequest is a partial update: omitted fields are unchanged
// and an explicit null clears a nullable field.
type UpdateProductModelRequest struct {
	ManufacturerID optional.Field[int64]  `json:"manufacturer_id" binding:"omitempty,gt=0"`
	ModelName      optional.Field[string] `json:"model_name" binding:"omitempty,max=200"`
	Description    optional.Field[string] `json:"description"`
	ReleaseDate    optional.Field[string] `json:"release_date" binding:"omitempty,datetime=2006-01-02"`
	IsActive       optional.Field[bool]   `json:"is_active"`
}

func (r UpdateProductModelRequest) Input() UpdateProductModelInput {
	return UpdateProductModelInput{
		ManufacturerID: r.ManufacturerID,
		ModelName:      r.ModelName,
		Description:    r.Description,
		ReleaseDate:    r.ReleaseDate,
		IsActive:       r.IsActive,
	}
}
