package catalog

import (
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

// CreateProductModelRequest is the JSON body of POST /admin/models.
type CreateProductModelRequest struct {
	ComponentTypeID int64   `json:"component_type_id"`
	ManufacturerID  int64   `json:"manufacturer_id"`
	ModelName       string  `json:"model_name"`
	Description     *string `json:"description"`
	ReleaseDate     *string `json:"release_date"` // YYYY-MM-DD
	IsActive        *bool   `json:"is_active"`
}

func (r CreateProductModelRequest) Validate() error {
	var v httpx.Validator
	v.RequiredID("component_type_id", r.ComponentTypeID)
	v.RequiredID("manufacturer_id", r.ManufacturerID)
	v.RequiredString("model_name", r.ModelName, 200)
	v.Date("release_date", r.ReleaseDate)
	return v.Err()
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

// UpdateProductModelRequest is the JSON body of PUT /admin/models/:id.
// It is a partial update: omitted fields are unchanged and an explicit null
// clears a nullable field. component_type_id cannot be changed.
type UpdateProductModelRequest struct {
	ManufacturerID optional.Field[int64]  `json:"manufacturer_id"`
	ModelName      optional.Field[string] `json:"model_name"`
	Description    optional.Field[string] `json:"description"`
	ReleaseDate    optional.Field[string] `json:"release_date"` // YYYY-MM-DD
	IsActive       optional.Field[bool]   `json:"is_active"`
}

func (r UpdateProductModelRequest) Validate() error {
	var v httpx.Validator
	httpx.CheckGreaterThan(&v, "manufacturer_id", r.ManufacturerID.Value, 0)
	v.MaxLength("model_name", r.ModelName.Value, 200)
	v.Date("release_date", r.ReleaseDate.Value)
	return v.Err()
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
