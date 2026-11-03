package catalog

import (
	"context"
	"strings"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
)

type ComponentTypeService struct {
	repo *ComponentTypeRepository
}

func NewComponentTypeService(repo *ComponentTypeRepository) *ComponentTypeService {
	return &ComponentTypeService{repo: repo}
}

func (s *ComponentTypeService) List(ctx context.Context, filter ComponentTypeFilter, page httpx.Page) ([]ComponentType, int64, error) {
	filter.Query = strings.TrimSpace(filter.Query)
	return s.repo.List(ctx, filter, page)
}

func (s *ComponentTypeService) Get(ctx context.Context, id int64) (ComponentType, error) {
	return s.repo.Get(ctx, id)
}

func (s *ComponentTypeService) Create(ctx context.Context, input CreateComponentTypeInput) (ComponentType, error) {
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	if input.Code == "" || input.Name == "" {
		return ComponentType{}, apperror.Invalid("code and name must not be blank")
	}
	return s.repo.Create(ctx, input)
}

// Update changes only the fields present in the input; code and name are fixed.
func (s *ComponentTypeService) Update(ctx context.Context, id int64, input UpdateComponentTypeInput) (ComponentType, error) {
	if input.IsActive.IsNull() {
		return ComponentType{}, apperror.Invalid("is_active must not be null")
	}
	return s.repo.Update(ctx, id, input)
}

// Delete fails with a conflict while models still use the type; deactivate it instead.
func (s *ComponentTypeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
