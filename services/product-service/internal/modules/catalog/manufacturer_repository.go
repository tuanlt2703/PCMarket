package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

const manufacturerColumns = "id, code, name, description, is_active, created_at, updated_at"

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

func (r *ManufacturerRepository) List(ctx context.Context, f ManufacturerFilter, p httpx.Page) ([]Manufacturer, int64, error) {
	const where = `
WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
  AND ($2::boolean IS NULL OR is_active = $2)`
	args := []any{f.Query, f.IsActive}

	var total int64
	if err := r.db.Q(ctx).QueryRow(ctx, "SELECT COUNT(*) FROM manufacturers"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Q(ctx).Query(ctx,
		"SELECT "+manufacturerColumns+" FROM manufacturers"+where+" ORDER BY name, id LIMIT $3 OFFSET $4",
		append(args, p.PageSize, p.Offset())...)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Manufacturer, error) { return scanManufacturer(row) })
	return items, total, err
}

func (r *ManufacturerRepository) Get(ctx context.Context, id int64) (Manufacturer, error) {
	m, err := scanManufacturer(r.db.Q(ctx).QueryRow(ctx,
		"SELECT "+manufacturerColumns+" FROM manufacturers WHERE id = $1", id))
	return m, database.MapError(err, manufacturerNotFound(id))
}

func (r *ManufacturerRepository) Create(ctx context.Context, in CreateManufacturerInput) (Manufacturer, error) {
	m, err := scanManufacturer(r.db.Q(ctx).QueryRow(ctx, `
INSERT INTO manufacturers (code, name, description, is_active)
VALUES ($1, $2, $3, COALESCE($4, TRUE))
RETURNING `+manufacturerColumns,
		in.Code, in.Name, in.Description, in.IsActive))
	return m, database.MapError(err, "")
}

func (r *ManufacturerRepository) Update(ctx context.Context, id int64, in UpdateManufacturerInput) (Manufacturer, error) {
	u := database.NewUpdates(id)
	database.Set(u, "code", in.Code)
	database.Set(u, "name", in.Name)
	database.Set(u, "description", in.Description)
	database.Set(u, "is_active", in.IsActive)
	if u.Empty() {
		return r.Get(ctx, id)
	}
	m, err := scanManufacturer(r.db.Q(ctx).QueryRow(ctx,
		u.SQL("manufacturers", manufacturerColumns), u.Args()...))
	return m, database.MapError(err, manufacturerNotFound(id))
}

func (r *ManufacturerRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Q(ctx).Exec(ctx, "DELETE FROM manufacturers WHERE id = $1", id)
	if err != nil {
		return database.MapError(err, "")
	}
	if tag.RowsAffected() == 0 {
		return database.MapError(pgx.ErrNoRows, manufacturerNotFound(id))
	}
	return nil
}

func manufacturerNotFound(id int64) string { return fmt.Sprintf("manufacturer %d not found", id) }
