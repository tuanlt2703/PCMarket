package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

// manufacturerColumns must stay in the same order as scanManufacturer.
const manufacturerColumns = "id, code, name, description, is_active, created_at, updated_at"

// ManufacturerRepository reads and writes the manufacturers table.
type ManufacturerRepository struct {
	db *database.DB
}

func NewManufacturerRepository(db *database.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

func scanManufacturer(row pgx.Row) (Manufacturer, error) {
	var m Manufacturer
	err := row.Scan(&m.ID, &m.Code, &m.Name, &m.Description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (r *ManufacturerRepository) List(ctx context.Context, filter ManufacturerFilter, page httpx.Page) ([]Manufacturer, int64, error) {
	// An empty search text ($1 = '') or a NULL is_active ($2) disables that filter.
	const where = `
WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
  AND ($2::boolean IS NULL OR is_active = $2)`

	var total int64
	countSQL := "SELECT COUNT(*) FROM manufacturers" + where
	err := r.db.Conn(ctx).QueryRow(ctx, countSQL, filter.Query, filter.IsActive).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listSQL := "SELECT " + manufacturerColumns + " FROM manufacturers" + where +
		" ORDER BY name, id LIMIT $3 OFFSET $4"
	rows, err := r.db.Conn(ctx).Query(ctx, listSQL, filter.Query, filter.IsActive, page.Size, page.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []Manufacturer
	for rows.Next() {
		m, err := scanManufacturer(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ManufacturerRepository) Get(ctx context.Context, id int64) (Manufacturer, error) {
	sql := "SELECT " + manufacturerColumns + " FROM manufacturers WHERE id = $1"
	m, err := scanManufacturer(r.db.Conn(ctx).QueryRow(ctx, sql, id))
	if err != nil {
		return Manufacturer{}, database.MapError(err, manufacturerNotFound(id))
	}
	return m, nil
}

func (r *ManufacturerRepository) Create(ctx context.Context, input CreateManufacturerInput) (Manufacturer, error) {
	sql := `
INSERT INTO manufacturers (code, name, description, is_active)
VALUES ($1, $2, $3, COALESCE($4, TRUE))
RETURNING ` + manufacturerColumns
	row := r.db.Conn(ctx).QueryRow(ctx, sql, input.Code, input.Name, input.Description, input.IsActive)
	m, err := scanManufacturer(row)
	if err != nil {
		return Manufacturer{}, database.MapError(err, "")
	}
	return m, nil
}

func (r *ManufacturerRepository) Update(ctx context.Context, id int64, input UpdateManufacturerInput) (Manufacturer, error) {
	update := database.NewUpdateBuilder(id)
	database.SetIfPresent(update, "code", input.Code)
	database.SetIfPresent(update, "name", input.Name)
	database.SetIfPresent(update, "description", input.Description)
	database.SetIfPresent(update, "is_active", input.IsActive)
	if update.IsEmpty() {
		return r.Get(ctx, id)
	}

	row := r.db.Conn(ctx).QueryRow(ctx, update.SQL("manufacturers", manufacturerColumns), update.Args()...)
	m, err := scanManufacturer(row)
	if err != nil {
		return Manufacturer{}, database.MapError(err, manufacturerNotFound(id))
	}
	return m, nil
}

func (r *ManufacturerRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Conn(ctx).Exec(ctx, "DELETE FROM manufacturers WHERE id = $1", id)
	if err != nil {
		return database.MapError(err, "")
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound(manufacturerNotFound(id))
	}
	return nil
}

func manufacturerNotFound(id int64) string {
	return fmt.Sprintf("manufacturer %d not found", id)
}
