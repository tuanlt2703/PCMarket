package catalog

import (
	"time"

	"product-service/internal/platform/optional"
)

type Ref struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProductModel struct {
	ID            int64     `json:"id"`
	ComponentType Ref       `json:"component_type"`
	Manufacturer  Ref       `json:"manufacturer"`
	ModelName     string    `json:"model_name"`
	Description   *string   `json:"description"`
	ReleaseDate   *string   `json:"release_date"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateProductModelInput struct {
	ComponentTypeID int64
	ManufacturerID  int64
	ModelName       string
	Description     *string
	ReleaseDate     *string
	IsActive        *bool
}

// UpdateProductModelInput is a partial update. It intentionally omits
// component_type_id: a model's type decides which specification table it owns,
// so changing it would orphan its specs.
type UpdateProductModelInput struct {
	ManufacturerID optional.Field[int64]
	ModelName      optional.Field[string]
	Description    optional.Field[string]
	ReleaseDate    optional.Field[string]
	IsActive       optional.Field[bool]
}

type ProductModelFilter struct {
	Query           string
	ComponentTypeID *int64
	ManufacturerID  *int64
	IsActive        *bool
}
