package product

import (
	"product-service/modules/catalog/componenttype"
	"product-service/modules/catalog/manufacturer"
	"time"
)

type Model struct {
	ID            int64                       `json:"id"`
	ComponentType componenttype.ComponentType `json:"component_type"`
	Manufacturer  manufacturer.Manufacturer   `json:"manufacturer"`
	ModelName     string                      `json:"model_name"`
	Description   *string                     `json:"description"`
	ReleaseDate   time.Time                   `json:"release_date"`
	IsActive      bool                        `json:"is_active"`
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
}

type CreateModel struct {
	ComponentTypeID int64     `json:"component_type_id"`
	ManufacturerID  int64     `json:"manufacturer_id"`
	ModelName       string    `json:"model_name"`
	Description     *string   `json:"description"`
	ReleaseDate     time.Time `json:"release_date"`
	IsActive        *bool     `json:"is_active"`
}

type UpdateModel struct {
	ModelName   string    `json:"model_name"`
	Description *string   `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	IsActive    *bool     `json:"is_active"`
}
