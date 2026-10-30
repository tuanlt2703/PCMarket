package discovery

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

// QueryRepository holds read-optimised queries for the storefront. It reads
// tables directly instead of calling other modules' services.
type QueryRepository struct {
	db *database.DB
}

func NewQueryRepository(db *database.DB) *QueryRepository {
	return &QueryRepository{db: db}
}

const productSelect = `
SELECT m.id, m.model_name, ct.code, ct.name, mf.code, mf.name, m.release_date::text, m.description
FROM models m
JOIN component_types ct ON ct.id = m.component_type_id
JOIN manufacturers mf ON mf.id = m.manufacturer_id`

// Only active models of active types and manufacturers are visible to customers.
const visible = `
WHERE m.is_active AND ct.is_active AND mf.is_active`

func scanProduct(row pgx.Row) (ProductDetail, error) {
	var p ProductDetail
	err := row.Scan(&p.ID, &p.Name, &p.ComponentType.Code, &p.ComponentType.Name,
		&p.Manufacturer.Code, &p.Manufacturer.Name, &p.ReleaseDate, &p.Description)
	return p, err
}

func (r *QueryRepository) Search(ctx context.Context, f SearchFilter, p httpx.Page) ([]ProductSummary, int64, error) {
	where := visible + `
  AND ($1 = '' OR m.model_name ILIKE '%' || $1 || '%')
  AND ($2 = '' OR ct.code = $2)
  AND ($3 = '' OR mf.code = $3)`
	args := []any{f.Query, f.ComponentType, f.Manufacturer}

	var total int64
	countSQL := `SELECT COUNT(*) FROM models m
JOIN component_types ct ON ct.id = m.component_type_id
JOIN manufacturers mf ON mf.id = m.manufacturer_id` + where
	if err := r.db.Q(ctx).QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	order := " ORDER BY m.model_name, m.id"
	if f.Sort == SortNewest {
		order = " ORDER BY m.release_date DESC NULLS LAST, m.id DESC"
	}

	rows, err := r.db.Q(ctx).Query(ctx, productSelect+where+order+" LIMIT $4 OFFSET $5",
		append(args, p.PageSize, p.Offset())...)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (ProductSummary, error) {
		d, err := scanProduct(row)
		return d.ProductSummary, err
	})
	return items, total, err
}

func (r *QueryRepository) Get(ctx context.Context, id int64) (ProductDetail, error) {
	p, err := scanProduct(r.db.Q(ctx).QueryRow(ctx, productSelect+visible+" AND m.id = $1", id))
	return p, database.MapError(err, fmt.Sprintf("product %d not found", id))
}
