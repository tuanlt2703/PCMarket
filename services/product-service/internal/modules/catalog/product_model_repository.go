package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

const modelSelect = `
SELECT m.id,
       ct.id, ct.code, ct.name,
       mf.id, mf.code, mf.name,
       m.model_name, m.description, m.release_date::text, m.is_active, m.created_at, m.updated_at
FROM models m
JOIN component_types ct ON ct.id = m.component_type_id
JOIN manufacturers mf ON mf.id = m.manufacturer_id`

type ProductModelRepository struct {
	db *database.DB
}

func NewProductModelRepository(db *database.DB) *ProductModelRepository {
	return &ProductModelRepository{db: db}
}

func scanModel(row pgx.Row) (ProductModel, error) {
	var m ProductModel
	err := row.Scan(&m.ID,
		&m.ComponentType.ID, &m.ComponentType.Code, &m.ComponentType.Name,
		&m.Manufacturer.ID, &m.Manufacturer.Code, &m.Manufacturer.Name,
		&m.ModelName, &m.Description, &m.ReleaseDate, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (r *ProductModelRepository) List(ctx context.Context, f ProductModelFilter, p httpx.Page) ([]ProductModel, int64, error) {
	where := `
WHERE ($1 = '' OR m.model_name ILIKE '%' || $1 || '%')
  AND ($2::bigint IS NULL OR m.component_type_id = $2)
  AND ($3::bigint IS NULL OR m.manufacturer_id = $3)
  AND ($4::boolean IS NULL OR m.is_active = $4)`
	args := []any{f.Query, f.ComponentTypeID, f.ManufacturerID, f.IsActive}

	var total int64
	if err := r.db.Q(ctx).QueryRow(ctx, "SELECT COUNT(*) FROM models m"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Q(ctx).Query(ctx, modelSelect+where+" ORDER BY m.model_name, m.id LIMIT $5 OFFSET $6",
		append(args, p.PageSize, p.Offset())...)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (ProductModel, error) { return scanModel(row) })
	return items, total, err
}

func (r *ProductModelRepository) Get(ctx context.Context, id int64) (ProductModel, error) {
	m, err := scanModel(r.db.Q(ctx).QueryRow(ctx, modelSelect+" WHERE m.id = $1", id))
	return m, database.MapError(err, modelNotFound(id))
}

func (r *ProductModelRepository) Create(ctx context.Context, in CreateProductModelInput) (int64, error) {
	var id int64
	err := r.db.Q(ctx).QueryRow(ctx, `
INSERT INTO models (component_type_id, manufacturer_id, model_name, description, release_date, is_active)
VALUES ($1, $2, $3, $4, $5::date, COALESCE($6, TRUE))
RETURNING id`,
		in.ComponentTypeID, in.ManufacturerID, in.ModelName, in.Description, in.ReleaseDate, in.IsActive).Scan(&id)
	return id, database.MapError(err, "")
}

func (r *ProductModelRepository) Update(ctx context.Context, id int64, in UpdateProductModelInput) error {
	u := database.NewUpdates(id)
	database.Set(u, "manufacturer_id", in.ManufacturerID)
	database.Set(u, "model_name", in.ModelName)
	database.Set(u, "description", in.Description)
	database.Set(u, "release_date", in.ReleaseDate)
	database.Set(u, "is_active", in.IsActive)
	if u.Empty() {
		_, err := r.Get(ctx, id)
		return err
	}
	tag, err := r.db.Q(ctx).Exec(ctx, u.SQL("models", ""), u.Args()...)
	if err != nil {
		return database.MapError(err, "")
	}
	if tag.RowsAffected() == 0 {
		return database.MapError(pgx.ErrNoRows, modelNotFound(id))
	}
	return nil
}

// Delete removes a model; its specification rows are removed by ON DELETE CASCADE.
func (r *ProductModelRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Q(ctx).Exec(ctx, "DELETE FROM models WHERE id = $1", id)
	if err != nil {
		return database.MapError(err, "")
	}
	if tag.RowsAffected() == 0 {
		return database.MapError(pgx.ErrNoRows, modelNotFound(id))
	}
	return nil
}

func (r *ProductModelRepository) ComponentTypeCode(ctx context.Context, modelID int64) (string, error) {
	var code string
	err := r.db.Q(ctx).QueryRow(ctx, `
SELECT ct.code FROM models m JOIN component_types ct ON ct.id = m.component_type_id
WHERE m.id = $1`, modelID).Scan(&code)
	return code, database.MapError(err, modelNotFound(modelID))
}

func modelNotFound(id int64) string { return fmt.Sprintf("model %d not found", id) }
