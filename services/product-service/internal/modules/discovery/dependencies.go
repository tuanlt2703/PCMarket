package discovery

import "context"

// This file declares what Discovery needs from the rest of the system.
// They are interfaces (like abstract base classes) so Discovery does not
// import the Specs module or a future SKU module directly. The concrete
// implementations are plugged in by internal/app.

// SpecReader loads the public specification of a model.
type SpecReader interface {
	SpecFor(ctx context.Context, modelID int64) (any, error)
}

// Offer is a sellable SKU of a model (price, availability).
type Offer struct {
	SKUID    int64   `json:"sku_id"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	InStock  bool    `json:"in_stock"`
}

// OfferReader loads the offers of several models in one call.
type OfferReader interface {
	OffersByModelIDs(ctx context.Context, modelIDs []int64) (map[int64][]Offer, error)
}

// NoOffers is a placeholder OfferReader that always returns no offers.
// SKU data does not exist yet; replace it once there is a SKU source.
type NoOffers struct{}

func (NoOffers) OffersByModelIDs(ctx context.Context, modelIDs []int64) (map[int64][]Offer, error) {
	return map[int64][]Offer{}, nil
}
