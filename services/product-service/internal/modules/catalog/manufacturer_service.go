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

func (s *ManufacturerService) List(ctx context.Context, filter ManufacturerFilter, page httpx.Page) ([]Manufacturer, int64, error) {
	filter.Query = strings.TrimSpace(filter.Query)
	return s.repo.List(ctx, filter, page)
}

func (s *ManufacturerService) Get(ctx context.Context, id int64) (Manufacturer, error) {
	return s.repo.Get(ctx, id)
}

func (s *ManufacturerService) Create(ctx context.Context, input CreateManufacturerInput) (Manufacturer, error) {
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	if input.Code == "" || input.Name == "" {
		return Manufacturer{}, apperror.Invalid("code and name must not be blank")
	}
	return s.repo.Create(ctx, input)
}

// Update changes only the fields present in the input.
func (s *ManufacturerService) Update(ctx context.Context, id int64, input UpdateManufacturerInput) (Manufacturer, error) {
	if optional.IsBlank(input.Code) || optional.IsBlank(input.Name) {
		return Manufacturer{}, apperror.Invalid("code and name must not be blank")
	}
	if input.IsActive.IsNull() {
		return Manufacturer{}, apperror.Invalid("is_active must not be null")
	}
	input.Code = optional.TrimSpace(input.Code)
	input.Name = optional.TrimSpace(input.Name)
	return s.repo.Update(ctx, id, input)
}

// Delete fails with a conflict while models still reference the manufacturer; deactivate it instead.
func (s *ManufacturerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
