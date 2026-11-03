package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

// componentTypeColumns must stay in the same order as scanComponentType.
const componentTypeColumns = "id, code, name, description, is_active, created_at, updated_at"

// ComponentTypeRepository reads and writes the component_types table.
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

func (r *ComponentTypeRepository) List(ctx context.Context, filter ComponentTypeFilter, page httpx.Page) ([]ComponentType, int64, error) {
	// An empty search text ($1 = '') or a NULL is_active ($2) disables that filter.
	const where = `
WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
  AND ($2::boolean IS NULL OR is_active = $2)`

	var total int64
	countSQL := "SELECT COUNT(*) FROM component_types" + where
	err := r.db.Conn(ctx).QueryRow(ctx, countSQL, filter.Query, filter.IsActive).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listSQL := "SELECT " + componentTypeColumns + " FROM component_types" + where +
		" ORDER BY name, id LIMIT $3 OFFSET $4"
	rows, err := r.db.Conn(ctx).Query(ctx, listSQL, filter.Query, filter.IsActive, page.Size, page.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []ComponentType
	for rows.Next() {
		ct, err := scanComponentType(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, ct)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ComponentTypeRepository) Get(ctx context.Context, id int64) (ComponentType, error) {
	sql := "SELECT " + componentTypeColumns + " FROM component_types WHERE id = $1"
	ct, err := scanComponentType(r.db.Conn(ctx).QueryRow(ctx, sql, id))
	if err != nil {
		return ComponentType{}, database.MapError(err, componentTypeNotFound(id))
	}
	return ct, nil
}

func (r *ComponentTypeRepository) Create(ctx context.Context, input CreateComponentTypeInput) (ComponentType, error) {
	sql := `
INSERT INTO component_types (code, name, description, is_active)
VALUES ($1, $2, $3, COALESCE($4, TRUE))
RETURNING ` + componentTypeColumns
	row := r.db.Conn(ctx).QueryRow(ctx, sql, input.Code, input.Name, input.Description, input.IsActive)
	ct, err := scanComponentType(row)
	if err != nil {
		return ComponentType{}, database.MapError(err, "")
	}
	return ct, nil
}

func (r *ComponentTypeRepository) Update(ctx context.Context, id int64, input UpdateComponentTypeInput) (ComponentType, error) {
	update := database.NewUpdateBuilder(id)
	database.SetIfPresent(update, "description", input.Description)
	database.SetIfPresent(update, "is_active", input.IsActive)
	if update.IsEmpty() {
		return r.Get(ctx, id)
	}

	row := r.db.Conn(ctx).QueryRow(ctx, update.SQL("component_types", componentTypeColumns), update.Args()...)
	ct, err := scanComponentType(row)
	if err != nil {
		return ComponentType{}, database.MapError(err, componentTypeNotFound(id))
	}
	return ct, nil
}

func (r *ComponentTypeRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Conn(ctx).Exec(ctx, "DELETE FROM component_types WHERE id = $1", id)
	if err != nil {
		return database.MapError(err, "")
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound(componentTypeNotFound(id))
	}
	return nil
}

func componentTypeNotFound(id int64) string {
	return fmt.Sprintf("component type %d not found", id)
}
