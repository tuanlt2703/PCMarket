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

func (s *ComponentTypeService) List(ctx context.Context, f ComponentTypeFilter, p httpx.Page) ([]ComponentType, int64, error) {
	f.Query = strings.TrimSpace(f.Query)
	return s.repo.List(ctx, f, p)
}

func (s *ComponentTypeService) Get(ctx context.Context, id int64) (ComponentType, error) {
	return s.repo.Get(ctx, id)
}

func (s *ComponentTypeService) Create(ctx context.Context, in CreateComponentTypeInput) (ComponentType, error) {
	in.Code = strings.TrimSpace(in.Code)
	in.Name = strings.TrimSpace(in.Name)
	if in.Code == "" || in.Name == "" {
		return ComponentType{}, apperror.Invalid("code and name must not be blank")
	}
	return s.repo.Create(ctx, in)
}

// Update changes only the fields present in the input; code and name are fixed.
func (s *ComponentTypeService) Update(ctx context.Context, id int64, in UpdateComponentTypeInput) (ComponentType, error) {
	if in.IsActive.IsNull() {
		return ComponentType{}, apperror.Invalid("is_active must not be null")
	}
	return s.repo.Update(ctx, id, in)
}

// Delete fails with a conflict while models still use the type; deactivate it instead.
func (s *ComponentTypeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
