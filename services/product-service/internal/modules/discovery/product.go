package discovery

type Ref struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProductSummary struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	ComponentType Ref     `json:"component_type"`
	Manufacturer  Ref     `json:"manufacturer"`
	ReleaseDate   *string `json:"release_date"`
	Offers        []Offer `json:"offers"`
}

type ProductDetail struct {
	ProductSummary
	Description *string `json:"description"`
	// Specs is null when the model has no specification yet.
	Specs any `json:"specs"`
}

type SearchFilter struct {
	Query         string
	ComponentType string
	Manufacturer  string
	Sort          string
}

const (
	SortName   = "name"
	SortNewest = "newest"
)
