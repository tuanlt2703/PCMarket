package discovery

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

// Repository runs read-only storefront queries. It reads the catalog tables
// directly instead of calling the Catalog module, so the queries can be
// tuned for the storefront.
type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

const fromProducts = `
FROM models m
JOIN component_types ct ON ct.id = m.component_type_id
JOIN manufacturers mf ON mf.id = m.manufacturer_id`

// selectProduct column order must match scanProduct.
const selectProduct = `
SELECT m.id, m.model_name, ct.code, ct.name, mf.code, mf.name, m.release_date::text, m.description` + fromProducts

// Customers only see active models whose type and manufacturer are also active.
const whereVisible = `
WHERE m.is_active AND ct.is_active AND mf.is_active`

func scanProduct(row pgx.Row) (ProductDetail, error) {
	var p ProductDetail
	err := row.Scan(&p.ID, &p.Name, &p.ComponentType.Code, &p.ComponentType.Name,
		&p.Manufacturer.Code, &p.Manufacturer.Name, &p.ReleaseDate, &p.Description)
	return p, err
}

// Search returns one page of visible products plus the total number of matches.
func (r *Repository) Search(ctx context.Context, filter SearchFilter, page httpx.Page) ([]ProductSummary, int64, error) {
	// An empty value ($1, $2, $3 = '') disables that filter.
	where := whereVisible + `
  AND ($1 = '' OR m.model_name ILIKE '%' || $1 || '%')
  AND ($2 = '' OR ct.code = $2)
  AND ($3 = '' OR mf.code = $3)`

	var total int64
	countSQL := "SELECT COUNT(*)" + fromProducts + where
	err := r.db.Conn(ctx).QueryRow(ctx, countSQL, filter.Query, filter.ComponentType, filter.Manufacturer).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := " ORDER BY m.model_name, m.id"
	if filter.Sort == SortNewest {
		orderBy = " ORDER BY m.release_date DESC NULLS LAST, m.id DESC"
	}

	listSQL := selectProduct + where + orderBy + " LIMIT $4 OFFSET $5"
	rows, err := r.db.Conn(ctx).Query(ctx, listSQL,
		filter.Query, filter.ComponentType, filter.Manufacturer, page.Size, page.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []ProductSummary
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, product.ProductSummary)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Get returns one visible product, without specs and offers.
func (r *Repository) Get(ctx context.Context, id int64) (ProductDetail, error) {
	sql := selectProduct + whereVisible + " AND m.id = $1"
	product, err := scanProduct(r.db.Conn(ctx).QueryRow(ctx, sql, id))
	if err != nil {
		return ProductDetail{}, database.MapError(err, fmt.Sprintf("product %d not found", id))
	}
	return product, nil
}
