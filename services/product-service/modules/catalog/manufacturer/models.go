package manufacturer

import "time"

type Manufacturer struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateManufacturer struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsActive    bool    `json:"is_active"`
}

type UpdateManufacturer struct {
	Description *string `json:"description"`
	IsActive    bool    `json:"is_active"`
}
