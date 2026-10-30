package lookup

import (
	"context"
	"strings"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

type ItemRepository interface {
	List(context.Context, Table, Filter, *httpx.Page) ([]Item, int64, error)
	Get(context.Context, Table, int64) (Item, error)
	Create(context.Context, Table, Input) (Item, error)
	Update(context.Context, Table, int64, UpdateInput) (Item, error)
	Delete(context.Context, Table, int64) error
}

type Service struct {
	repo     ItemRepository
	registry *Registry
}

func NewService(repo ItemRepository, registry *Registry) *Service {
	return &Service{repo: repo, registry: registry}
}

func (s *Service) ListTypes() []string { return s.registry.Slugs() }

func (s *Service) table(slug string) (Table, error) {
	t, ok := s.registry.Get(slug)
	if !ok {
		return Table{}, apperror.NotFound("unknown lookup type " + slug)
	}
	return t, nil
}

func (s *Service) List(ctx context.Context, slug string, f Filter, p httpx.Page) ([]Item, int64, error) {
	t, err := s.table(slug)
	if err != nil {
		return nil, 0, err
	}
	f.Query = strings.TrimSpace(f.Query)
	return s.repo.List(ctx, t, f, &p)
}

// ListActive returns all active rows without pagination; reference tables are small.
func (s *Service) ListActive(ctx context.Context, slug string) ([]Item, error) {
	t, err := s.table(slug)
	if err != nil {
		return nil, err
	}
	active := true
	items, _, err := s.repo.List(ctx, t, Filter{IsActive: &active}, nil)
	return items, err
}

func (s *Service) Get(ctx context.Context, slug string, id int64) (Item, error) {
	t, err := s.table(slug)
	if err != nil {
		return Item{}, err
	}
	return s.repo.Get(ctx, t, id)
}

func (s *Service) Create(ctx context.Context, slug string, in Input) (Item, error) {
	t, err := s.table(slug)
	if err != nil {
		return Item{}, err
	}
	if err := normalize(&in); err != nil {
		return Item{}, err
	}
	return s.repo.Create(ctx, t, in)
}

// Update changes only the fields present in the input.
func (s *Service) Update(ctx context.Context, slug string, id int64, in UpdateInput) (Item, error) {
	t, err := s.table(slug)
	if err != nil {
		return Item{}, err
	}
	if optional.Blank(in.Code) || optional.Blank(in.Name) {
		return Item{}, apperror.Invalid("code and name must not be blank")
	}
	if in.IsActive.IsNull() {
		return Item{}, apperror.Invalid("is_active must not be null")
	}
	in.Code = optional.TrimSpace(in.Code)
	in.Name = optional.TrimSpace(in.Name)
	return s.repo.Update(ctx, t, id, in)
}

// Delete removes a row. Rows still referenced by other data return a conflict;
// deactivate them via is_active instead.
func (s *Service) Delete(ctx context.Context, slug string, id int64) error {
	t, err := s.table(slug)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, t, id)
}

func normalize(in *Input) error {
	in.Code = strings.TrimSpace(in.Code)
	in.Name = strings.TrimSpace(in.Name)
	if in.Code == "" || in.Name == "" {
		return apperror.Invalid("code and name must not be blank")
	}
	return nil
}
