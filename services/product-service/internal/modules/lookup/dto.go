package lookup

import "product-service/internal/platform/optional"

type CreateLookupRequest struct {
	Code        string  `json:"code" binding:"required,max=50"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description"`
	// IsActive defaults to true when omitted.
	IsActive *bool `json:"is_active"`
}

func (r CreateLookupRequest) Input() Input {
	return Input{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// UpdateLookupRequest is a partial update: omitted fields are unchanged and an
// explicit null clears a nullable field.
type UpdateLookupRequest struct {
	Code        optional.Field[string] `json:"code" binding:"omitempty,max=50"`
	Name        optional.Field[string] `json:"name" binding:"omitempty,max=100"`
	Description optional.Field[string] `json:"description"`
	IsActive    optional.Field[bool]   `json:"is_active"`
}

func (r UpdateLookupRequest) Input() UpdateInput {
	return UpdateInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// PublicItem hides audit fields from storefront clients.
type PublicItem struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func toPublic(items []Item) []PublicItem {
	out := make([]PublicItem, len(items))
	for i, it := range items {
		out[i] = PublicItem{ID: it.ID, Code: it.Code, Name: it.Name, Description: it.Description}
	}
	return out
}
