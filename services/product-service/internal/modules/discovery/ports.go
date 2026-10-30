package discovery

import "context"

// SpecReader loads a model's specification. The Specs module satisfies it via an
// adapter in app wiring, so Discovery never depends on Specs internals.
type SpecReader interface {
	SpecFor(ctx context.Context, modelID int64) (any, error)
}

// Offer is a sellable SKU of a model (price, availability). SKU data is out of
// scope for now; replace NoopOfferReader once the SKU source exists.
type Offer struct {
	SKUID    int64   `json:"sku_id"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	InStock  bool    `json:"in_stock"`
}

type OfferReader interface {
	OffersByModelIDs(ctx context.Context, modelIDs []int64) (map[int64][]Offer, error)
}

type NoopOfferReader struct{}

func (NoopOfferReader) OffersByModelIDs(context.Context, []int64) (map[int64][]Offer, error) {
	return map[int64][]Offer{}, nil
}
