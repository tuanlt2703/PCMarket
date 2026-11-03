package lookup

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

// itemColumns must stay in the same order as scanItem.
const itemColumns = "id, code, name, description, is_active, created_at, updated_at"

// Repository reads and writes any lookup table. The table is passed to each
// method because all lookup tables share the same columns.
type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// quotedTableName returns the table name quoted for safe use inside SQL text.
func quotedTableName(table Table) string {
	return pgx.Identifier{table.Name}.Sanitize()
}

func scanItem(row pgx.Row) (Item, error) {
	var item Item
	err := row.Scan(&item.ID, &item.Code, &item.Name, &item.Description, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func readItems(rows pgx.Rows) ([]Item, error) {
	defer rows.Close()
	var items []Item
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// List returns one page of rows matching filter, plus the total number of matches.
func (r *Repository) List(ctx context.Context, table Table, filter Filter, page httpx.Page) ([]Item, int64, error) {
	// An empty search text ($1 = '') or a NULL is_active ($2) disables that filter.
	where := `
WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
  AND ($2::boolean IS NULL OR is_active = $2)`

	var total int64
	countSQL := "SELECT COUNT(*) FROM " + quotedTableName(table) + where
	err := r.db.Conn(ctx).QueryRow(ctx, countSQL, filter.Query, filter.IsActive).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listSQL := "SELECT " + itemColumns + " FROM " + quotedTableName(table) + where +
		" ORDER BY name, id LIMIT $3 OFFSET $4"
	rows, err := r.db.Conn(ctx).Query(ctx, listSQL, filter.Query, filter.IsActive, page.Size, page.Offset())
	if err != nil {
		return nil, 0, err
	}
	items, err := readItems(rows)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListActive returns every active row without paging; lookup tables are small.
func (r *Repository) ListActive(ctx context.Context, table Table) ([]Item, error) {
	sql := "SELECT " + itemColumns + " FROM " + quotedTableName(table) + " WHERE is_active ORDER BY name, id"
	rows, err := r.db.Conn(ctx).Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return readItems(rows)
}

func (r *Repository) Get(ctx context.Context, table Table, id int64) (Item, error) {
	sql := "SELECT " + itemColumns + " FROM " + quotedTableName(table) + " WHERE id = $1"
	item, err := scanItem(r.db.Conn(ctx).QueryRow(ctx, sql, id))
	if err != nil {
		return Item{}, database.MapError(err, notFoundMessage(table, id))
	}
	return item, nil
}

func (r *Repository) Create(ctx context.Context, table Table, input CreateInput) (Item, error) {
	sql := "INSERT INTO " + quotedTableName(table) + ` (code, name, description, is_active)
VALUES ($1, $2, $3, COALESCE($4, TRUE))
RETURNING ` + itemColumns
	row := r.db.Conn(ctx).QueryRow(ctx, sql, input.Code, input.Name, input.Description, input.IsActive)
	item, err := scanItem(row)
	if err != nil {
		return Item{}, database.MapError(err, "")
	}
	return item, nil
}

func (r *Repository) Update(ctx context.Context, table Table, id int64, input UpdateInput) (Item, error) {
	update := database.NewUpdateBuilder(id)
	database.SetIfPresent(update, "code", input.Code)
	database.SetIfPresent(update, "name", input.Name)
	database.SetIfPresent(update, "description", input.Description)
	database.SetIfPresent(update, "is_active", input.IsActive)
	if update.IsEmpty() {
		return r.Get(ctx, table, id)
	}

	row := r.db.Conn(ctx).QueryRow(ctx, update.SQL(quotedTableName(table), itemColumns), update.Args()...)
	item, err := scanItem(row)
	if err != nil {
		return Item{}, database.MapError(err, notFoundMessage(table, id))
	}
	return item, nil
}

func (r *Repository) Delete(ctx context.Context, table Table, id int64) error {
	sql := "DELETE FROM " + quotedTableName(table) + " WHERE id = $1"
	result, err := r.db.Conn(ctx).Exec(ctx, sql, id)
	if err != nil {
		return database.MapError(err, "")
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound(notFoundMessage(table, id))
	}
	return nil
}

func notFoundMessage(table Table, id int64) string {
	return fmt.Sprintf("%s %d not found", table.Label, id)
}
