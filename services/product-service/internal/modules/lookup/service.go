package lookup

import (
	"context"
	"strings"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
	"product-service/internal/platform/optional"
)

// ItemRepository is the storage used by Service. *Repository implements it
// with PostgreSQL; tests use an in-memory stub.
type ItemRepository interface {
	List(ctx context.Context, table Table, filter Filter, page httpx.Page) ([]Item, int64, error)
	ListActive(ctx context.Context, table Table) ([]Item, error)
	Get(ctx context.Context, table Table, id int64) (Item, error)
	Create(ctx context.Context, table Table, input CreateInput) (Item, error)
	Update(ctx context.Context, table Table, id int64, input UpdateInput) (Item, error)
	Delete(ctx context.Context, table Table, id int64) error
}

// Service handles every lookup type. Each method takes the URL slug
// (e.g. "cpu-sockets"), finds the matching table in the registry and then
// calls the repository.
type Service struct {
	repo     ItemRepository
	registry *Registry
}

func NewService(repo ItemRepository, registry *Registry) *Service {
	return &Service{repo: repo, registry: registry}
}

// ListTypes returns the slugs of every lookup type, sorted.
func (s *Service) ListTypes() []string {
	return s.registry.Slugs()
}

func (s *Service) findTable(slug string) (Table, error) {
	table, found := s.registry.Get(slug)
	if !found {
		return Table{}, apperror.NotFound("unknown lookup type " + slug)
	}
	return table, nil
}

func (s *Service) List(ctx context.Context, slug string, filter Filter, page httpx.Page) ([]Item, int64, error) {
	table, err := s.findTable(slug)
	if err != nil {
		return nil, 0, err
	}
	filter.Query = strings.TrimSpace(filter.Query)
	return s.repo.List(ctx, table, filter, page)
}

// ListActive returns all active rows of one lookup type, for the storefront.
func (s *Service) ListActive(ctx context.Context, slug string) ([]Item, error) {
	table, err := s.findTable(slug)
	if err != nil {
		return nil, err
	}
	return s.repo.ListActive(ctx, table)
}

func (s *Service) Get(ctx context.Context, slug string, id int64) (Item, error) {
	table, err := s.findTable(slug)
	if err != nil {
		return Item{}, err
	}
	return s.repo.Get(ctx, table, id)
}

func (s *Service) Create(ctx context.Context, slug string, input CreateInput) (Item, error) {
	table, err := s.findTable(slug)
	if err != nil {
		return Item{}, err
	}
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	if input.Code == "" || input.Name == "" {
		return Item{}, apperror.Invalid("code and name must not be blank")
	}
	return s.repo.Create(ctx, table, input)
}

// Update changes only the fields present in the input.
func (s *Service) Update(ctx context.Context, slug string, id int64, input UpdateInput) (Item, error) {
	table, err := s.findTable(slug)
	if err != nil {
		return Item{}, err
	}
	if optional.IsBlank(input.Code) || optional.IsBlank(input.Name) {
		return Item{}, apperror.Invalid("code and name must not be blank")
	}
	if input.IsActive.IsNull() {
		return Item{}, apperror.Invalid("is_active must not be null")
	}
	input.Code = optional.TrimSpace(input.Code)
	input.Name = optional.TrimSpace(input.Name)
	return s.repo.Update(ctx, table, id, input)
}

// Delete removes a row. A row still used by other data returns a conflict;
// set is_active = false instead.
func (s *Service) Delete(ctx context.Context, slug string, id int64) error {
	table, err := s.findTable(slug)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, table, id)
}
