package discovery

import (
	"context"
	"errors"
	"strings"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
)

// Service builds the storefront product views: catalog data from Repository,
// plus specs and offers from the readers plugged in by internal/app.
type Service struct {
	repo   *Repository
	specs  SpecReader
	offers OfferReader
}

func NewService(repo *Repository, specs SpecReader, offers OfferReader) *Service {
	return &Service{repo: repo, specs: specs, offers: offers}
}

func (s *Service) Search(ctx context.Context, filter SearchFilter, page httpx.Page) ([]ProductSummary, int64, error) {
	filter.Query = strings.TrimSpace(filter.Query)
	if filter.Sort != "" && filter.Sort != SortName && filter.Sort != SortNewest {
		return nil, 0, apperror.Invalid("sort must be one of: name, newest")
	}

	items, total, err := s.repo.Search(ctx, filter, page)
	if err != nil {
		return nil, 0, err
	}

	// Load the offers of the whole page in one call.
	modelIDs := make([]int64, len(items))
	for i, item := range items {
		modelIDs[i] = item.ID
	}
	offersByModel, err := s.offers.OffersByModelIDs(ctx, modelIDs)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].Offers = emptyIfNil(offersByModel[items[i].ID])
	}
	return items, total, nil
}

func (s *Service) Get(ctx context.Context, id int64) (ProductDetail, error) {
	product, err := s.repo.Get(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}

	spec, err := s.specs.SpecFor(ctx, id)
	switch {
	case err == nil:
		product.Specs = spec
	case errors.Is(err, apperror.ErrNotFound), errors.Is(err, apperror.ErrInvalid):
		// The model has no spec yet, or its component type has no spec
		// support yet: show the product with "specs": null.
	default:
		return ProductDetail{}, err
	}

	offersByModel, err := s.offers.OffersByModelIDs(ctx, []int64{id})
	if err != nil {
		return ProductDetail{}, err
	}
	product.Offers = emptyIfNil(offersByModel[id])
	return product, nil
}

// emptyIfNil makes the JSON output [] instead of null.
func emptyIfNil(offers []Offer) []Offer {
	if offers == nil {
		return []Offer{}
	}
	return offers
}
