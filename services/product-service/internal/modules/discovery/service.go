package discovery

import (
	"context"
	"errors"
	"strings"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
)

type Service struct {
	repo   *QueryRepository
	specs  SpecReader
	offers OfferReader
}

func NewService(repo *QueryRepository, specs SpecReader, offers OfferReader) *Service {
	return &Service{repo: repo, specs: specs, offers: offers}
}

func (s *Service) Search(ctx context.Context, f SearchFilter, p httpx.Page) ([]ProductSummary, int64, error) {
	f.Query = strings.TrimSpace(f.Query)
	if f.Sort != "" && f.Sort != SortName && f.Sort != SortNewest {
		return nil, 0, apperror.Invalid("sort must be one of: name, newest")
	}

	items, total, err := s.repo.Search(ctx, f, p)
	if err != nil {
		return nil, 0, err
	}

	ids := make([]int64, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	offers, err := s.offers.OffersByModelIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].Offers = nonNil(offers[items[i].ID])
	}
	return items, total, nil
}

func (s *Service) Get(ctx context.Context, id int64) (ProductDetail, error) {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return p, err
	}

	spec, err := s.specs.SpecFor(ctx, id)
	switch {
	case err == nil:
		p.Specs = spec
	case errors.Is(err, apperror.ErrNotFound), errors.Is(err, apperror.ErrInvalid):
		// No spec yet, or component type without spec support: show the product without specs.
	default:
		return p, err
	}

	offers, err := s.offers.OffersByModelIDs(ctx, []int64{id})
	if err != nil {
		return p, err
	}
	p.Offers = nonNil(offers[id])
	return p, nil
}

func nonNil(o []Offer) []Offer {
	if o == nil {
		return []Offer{}
	}
	return o
}
