package catalog

import (
	"context"
	"strings"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

type ProductModelService struct {
	repo *ProductModelRepository
}

func NewProductModelService(repo *ProductModelRepository) *ProductModelService {
	return &ProductModelService{repo: repo}
}

func (s *ProductModelService) List(ctx context.Context, f ProductModelFilter, p httpx.Page) ([]ProductModel, int64, error) {
	f.Query = strings.TrimSpace(f.Query)
	return s.repo.List(ctx, f, p)
}

func (s *ProductModelService) Get(ctx context.Context, id int64) (ProductModel, error) {
	return s.repo.Get(ctx, id)
}

func (s *ProductModelService) Create(ctx context.Context, in CreateProductModelInput) (ProductModel, error) {
	in.ModelName = strings.TrimSpace(in.ModelName)
	if in.ModelName == "" {
		return ProductModel{}, apperror.Invalid("model_name must not be blank")
	}
	id, err := s.repo.Create(ctx, in)
	if err != nil {
		return ProductModel{}, err
	}
	return s.repo.Get(ctx, id)
}

func (s *ProductModelService) Update(ctx context.Context, id int64, in UpdateProductModelInput) (ProductModel, error) {
	if optional.Blank(in.ModelName) {
		return ProductModel{}, apperror.Invalid("model_name must not be blank")
	}
	if in.ManufacturerID.IsNull() {
		return ProductModel{}, apperror.Invalid("manufacturer_id must not be null")
	}
	if in.IsActive.IsNull() {
		return ProductModel{}, apperror.Invalid("is_active must not be null")
	}
	in.ModelName = optional.TrimSpace(in.ModelName)
	if err := s.repo.Update(ctx, id, in); err != nil {
		return ProductModel{}, err
	}
	return s.repo.Get(ctx, id)
}

func (s *ProductModelService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// ComponentTypeCode returns the component type code (e.g. "cpu") of a model.
// It is used by the Specs module to pick the matching specification handler.
func (s *ProductModelService) ComponentTypeCode(ctx context.Context, modelID int64) (string, error) {
	return s.repo.ComponentTypeCode(ctx, modelID)
}
