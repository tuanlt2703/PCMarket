package catalog

import (
	"context"
	"strings"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

type ManufacturerService struct {
	repo *ManufacturerRepository
}

func NewManufacturerService(repo *ManufacturerRepository) *ManufacturerService {
	return &ManufacturerService{repo: repo}
}

func (s *ManufacturerService) List(ctx context.Context, f ManufacturerFilter, p httpx.Page) ([]Manufacturer, int64, error) {
	f.Query = strings.TrimSpace(f.Query)
	return s.repo.List(ctx, f, p)
}

func (s *ManufacturerService) Get(ctx context.Context, id int64) (Manufacturer, error) {
	return s.repo.Get(ctx, id)
}

func (s *ManufacturerService) Create(ctx context.Context, in CreateManufacturerInput) (Manufacturer, error) {
	in.Code = strings.TrimSpace(in.Code)
	in.Name = strings.TrimSpace(in.Name)
	if in.Code == "" || in.Name == "" {
		return Manufacturer{}, apperror.Invalid("code and name must not be blank")
	}
	return s.repo.Create(ctx, in)
}

// Update changes only the fields present in the input.
func (s *ManufacturerService) Update(ctx context.Context, id int64, in UpdateManufacturerInput) (Manufacturer, error) {
	if optional.Blank(in.Code) || optional.Blank(in.Name) {
		return Manufacturer{}, apperror.Invalid("code and name must not be blank")
	}
	if in.IsActive.IsNull() {
		return Manufacturer{}, apperror.Invalid("is_active must not be null")
	}
	in.Code = optional.TrimSpace(in.Code)
	in.Name = optional.TrimSpace(in.Name)
	return s.repo.Update(ctx, id, in)
}

// Delete fails with a conflict while models still reference the manufacturer; deactivate it instead.
func (s *ManufacturerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
