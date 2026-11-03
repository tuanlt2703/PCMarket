package catalog

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

// selectProductModel joins the type and manufacturer so a model can be
// returned with their code and name. Column order must match scanProductModel.
const selectProductModel = `
SELECT m.id,
       ct.id, ct.code, ct.name,
       mf.id, mf.code, mf.name,
       m.model_name, m.description, m.release_date::text, m.is_active, m.created_at, m.updated_at
FROM models m
JOIN component_types ct ON ct.id = m.component_type_id
JOIN manufacturers mf ON mf.id = m.manufacturer_id`

// ProductModelRepository reads and writes the models table.
type ProductModelRepository struct {
	db *database.DB
}

func NewProductModelRepository(db *database.DB) *ProductModelRepository {
	return &ProductModelRepository{db: db}
}

func scanProductModel(row pgx.Row) (ProductModel, error) {
	var m ProductModel
	err := row.Scan(&m.ID,
		&m.ComponentType.ID, &m.ComponentType.Code, &m.ComponentType.Name,
		&m.Manufacturer.ID, &m.Manufacturer.Code, &m.Manufacturer.Name,
		&m.ModelName, &m.Description, &m.ReleaseDate, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (r *ProductModelRepository) List(ctx context.Context, filter ProductModelFilter, page httpx.Page) ([]ProductModel, int64, error) {
	// An empty search text ($1 = '') or a NULL id/flag ($2..$4) disables that filter.
	const where = `
WHERE ($1 = '' OR m.model_name ILIKE '%' || $1 || '%')
  AND ($2::bigint IS NULL OR m.component_type_id = $2)
  AND ($3::bigint IS NULL OR m.manufacturer_id = $3)
  AND ($4::boolean IS NULL OR m.is_active = $4)`
	args := []any{filter.Query, filter.ComponentTypeID, filter.ManufacturerID, filter.IsActive}

	var total int64
	err := r.db.Conn(ctx).QueryRow(ctx, "SELECT COUNT(*) FROM models m"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listSQL := selectProductModel + where + " ORDER BY m.model_name, m.id LIMIT $5 OFFSET $6"
	args = append(args, page.Size, page.Offset())
	rows, err := r.db.Conn(ctx).Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []ProductModel
	for rows.Next() {
		m, err := scanProductModel(rows)
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

func (r *ProductModelRepository) Get(ctx context.Context, id int64) (ProductModel, error) {
	m, err := scanProductModel(r.db.Conn(ctx).QueryRow(ctx, selectProductModel+" WHERE m.id = $1", id))
	if err != nil {
		return ProductModel{}, database.MapError(err, productModelNotFound(id))
	}
	return m, nil
}

// Create inserts a model and returns its new id.
func (r *ProductModelRepository) Create(ctx context.Context, input CreateProductModelInput) (int64, error) {
	const sql = `
INSERT INTO models (component_type_id, manufacturer_id, model_name, description, release_date, is_active)
VALUES ($1, $2, $3, $4, $5::date, COALESCE($6, TRUE))
RETURNING id`
	var id int64
	err := r.db.Conn(ctx).QueryRow(ctx, sql,
		input.ComponentTypeID, input.ManufacturerID, input.ModelName,
		input.Description, input.ReleaseDate, input.IsActive).Scan(&id)
	if err != nil {
		return 0, database.MapError(err, "")
	}
	return id, nil
}

func (r *ProductModelRepository) Update(ctx context.Context, id int64, input UpdateProductModelInput) error {
	update := database.NewUpdateBuilder(id)
	database.SetIfPresent(update, "manufacturer_id", input.ManufacturerID)
	database.SetIfPresent(update, "model_name", input.ModelName)
	database.SetIfPresent(update, "description", input.Description)
	database.SetIfPresent(update, "release_date", input.ReleaseDate)
	database.SetIfPresent(update, "is_active", input.IsActive)
	if update.IsEmpty() {
		// Nothing to write, but still report 404 for an unknown id.
		_, err := r.Get(ctx, id)
		return err
	}

	result, err := r.db.Conn(ctx).Exec(ctx, update.SQL("models", ""), update.Args()...)
	if err != nil {
		return database.MapError(err, "")
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound(productModelNotFound(id))
	}
	return nil
}

// Delete removes a model. Its specification rows are removed by ON DELETE CASCADE.
func (r *ProductModelRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Conn(ctx).Exec(ctx, "DELETE FROM models WHERE id = $1", id)
	if err != nil {
		return database.MapError(err, "")
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound(productModelNotFound(id))
	}
	return nil
}

// ComponentTypeCode returns the code of the model's component type, e.g. "cpu".
func (r *ProductModelRepository) ComponentTypeCode(ctx context.Context, modelID int64) (string, error) {
	const sql = `
SELECT ct.code
FROM models m
JOIN component_types ct ON ct.id = m.component_type_id
WHERE m.id = $1`
	var code string
	if err := r.db.Conn(ctx).QueryRow(ctx, sql, modelID).Scan(&code); err != nil {
		return "", database.MapError(err, productModelNotFound(modelID))
	}
	return code, nil
}

func productModelNotFound(id int64) string {
	return fmt.Sprintf("model %d not found", id)
}
