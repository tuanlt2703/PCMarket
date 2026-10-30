package lookup

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"product-service/internal/platform/database"
	"product-service/internal/platform/httpx"
)

const columns = "id, code, name, description, is_active, created_at, updated_at"

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func ident(t Table) string { return pgx.Identifier{t.Name}.Sanitize() }

func scanItem(row pgx.Row) (Item, error) {
	var it Item
	err := row.Scan(&it.ID, &it.Code, &it.Name, &it.Description, &it.IsActive, &it.CreatedAt, &it.UpdatedAt)
	return it, err
}

func (r *Repository) List(ctx context.Context, t Table, f Filter, p *httpx.Page) ([]Item, int64, error) {
	where := `WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
	            AND ($2::boolean IS NULL OR is_active = $2)`
	args := []any{f.Query, f.IsActive}

	var total int64
	if p != nil {
		countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", ident(t), where)
		if err := r.db.Q(ctx).QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	sql := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY name, id", columns, ident(t), where)
	if p != nil {
		sql += " LIMIT $3 OFFSET $4"
		args = append(args, p.PageSize, p.Offset())
	}

	rows, err := r.db.Q(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Item, error) { return scanItem(row) })
	if err != nil {
		return nil, 0, err
	}
	if p == nil {
		total = int64(len(items))
	}
	return items, total, nil
}

func (r *Repository) Get(ctx context.Context, t Table, id int64) (Item, error) {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", columns, ident(t))
	it, err := scanItem(r.db.Q(ctx).QueryRow(ctx, sql, id))
	return it, database.MapError(err, notFound(t, id))
}

func (r *Repository) Create(ctx context.Context, t Table, in Input) (Item, error) {
	sql := fmt.Sprintf(`INSERT INTO %s (code, name, description, is_active)
		VALUES ($1, $2, $3, COALESCE($4, TRUE))
		RETURNING %s`, ident(t), columns)
	it, err := scanItem(r.db.Q(ctx).QueryRow(ctx, sql, in.Code, in.Name, in.Description, in.IsActive))
	return it, database.MapError(err, "")
}

func (r *Repository) Update(ctx context.Context, t Table, id int64, in UpdateInput) (Item, error) {
	u := database.NewUpdates(id)
	database.Set(u, "code", in.Code)
	database.Set(u, "name", in.Name)
	database.Set(u, "description", in.Description)
	database.Set(u, "is_active", in.IsActive)
	if u.Empty() {
		return r.Get(ctx, t, id)
	}
	it, err := scanItem(r.db.Q(ctx).QueryRow(ctx, u.SQL(ident(t), columns), u.Args()...))
	return it, database.MapError(err, notFound(t, id))
}

func (r *Repository) Delete(ctx context.Context, t Table, id int64) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE id = $1", ident(t))
	tag, err := r.db.Q(ctx).Exec(ctx, sql, id)
	if err != nil {
		return database.MapError(err, "")
	}
	if tag.RowsAffected() == 0 {
		return database.MapError(pgx.ErrNoRows, notFound(t, id))
	}
	return nil
}

func notFound(t Table, id int64) string {
	return fmt.Sprintf("%s %d not found", t.Label, id)
}
