package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

const componentTypeColumns = "id, code, name, description, is_active, created_at, updated_at"

type ComponentTypeRepository struct {
	db *database.DB
}

func NewComponentTypeRepository(db *database.DB) *ComponentTypeRepository {
	return &ComponentTypeRepository{db: db}
}

func scanComponentType(row pgx.Row) (ComponentType, error) {
	var ct ComponentType
	err := row.Scan(&ct.ID, &ct.Code, &ct.Name, &ct.Description, &ct.IsActive, &ct.CreatedAt, &ct.UpdatedAt)
	return ct, err
}

func (r *ComponentTypeRepository) List(ctx context.Context, f ComponentTypeFilter, p httpx.Page) ([]ComponentType, int64, error) {
	const where = `
WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
  AND ($2::boolean IS NULL OR is_active = $2)`
	args := []any{f.Query, f.IsActive}

	var total int64
	if err := r.db.Q(ctx).QueryRow(ctx, "SELECT COUNT(*) FROM component_types"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Q(ctx).Query(ctx,
		"SELECT "+componentTypeColumns+" FROM component_types"+where+" ORDER BY name, id LIMIT $3 OFFSET $4",
		append(args, p.PageSize, p.Offset())...)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (ComponentType, error) { return scanComponentType(row) })
	return items, total, err
}

func (r *ComponentTypeRepository) Get(ctx context.Context, id int64) (ComponentType, error) {
	ct, err := scanComponentType(r.db.Q(ctx).QueryRow(ctx,
		"SELECT "+componentTypeColumns+" FROM component_types WHERE id = $1", id))
	return ct, database.MapError(err, componentTypeNotFound(id))
}

func (r *ComponentTypeRepository) Create(ctx context.Context, in CreateComponentTypeInput) (ComponentType, error) {
	ct, err := scanComponentType(r.db.Q(ctx).QueryRow(ctx, `
INSERT INTO component_types (code, name, description, is_active)
VALUES ($1, $2, $3, COALESCE($4, TRUE))
RETURNING `+componentTypeColumns,
		in.Code, in.Name, in.Description, in.IsActive))
	return ct, database.MapError(err, "")
}

func (r *ComponentTypeRepository) Update(ctx context.Context, id int64, in UpdateComponentTypeInput) (ComponentType, error) {
	u := database.NewUpdates(id)
	database.Set(u, "description", in.Description)
	database.Set(u, "is_active", in.IsActive)
	if u.Empty() {
		return r.Get(ctx, id)
	}
	ct, err := scanComponentType(r.db.Q(ctx).QueryRow(ctx,
		u.SQL("component_types", componentTypeColumns), u.Args()...))
	return ct, database.MapError(err, componentTypeNotFound(id))
}

func (r *ComponentTypeRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Q(ctx).Exec(ctx, "DELETE FROM component_types WHERE id = $1", id)
	if err != nil {
		return database.MapError(err, "")
	}
	if tag.RowsAffected() == 0 {
		return database.MapError(pgx.ErrNoRows, componentTypeNotFound(id))
	}
	return nil
}

func componentTypeNotFound(id int64) string { return fmt.Sprintf("component type %d not found", id) }
