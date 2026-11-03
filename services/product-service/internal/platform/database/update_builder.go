package database

import (
	"fmt"
	"strings"

	"product-service/internal/platform/optional"
)

// UpdateBuilder builds a partial "UPDATE table SET ... WHERE id = $1" statement.
// Only fields the client sent are written; an explicit null writes NULL.
//
// Usage:
//
//	update := database.NewUpdateBuilder(id)
//	database.SetIfPresent(update, "name", input.Name)
//	if !update.IsEmpty() {
//		db.Conn(ctx).Exec(ctx, update.SQL("manufacturers", ""), update.Args()...)
//	}
type UpdateBuilder struct {
	assignments []string // e.g. "name = $2"
	args        []any    // args[0] is the row id ($1)
}

func NewUpdateBuilder(id any) *UpdateBuilder {
	return &UpdateBuilder{args: []any{id}}
}

// SetIfPresent adds "column = $n" when field is present.
// column must be a constant written in code, never a value from the request.
//
// This is a function and not a method because Go methods cannot have their
// own type parameters.
func SetIfPresent[T any](update *UpdateBuilder, column string, field optional.Field[T]) {
	if !field.Present {
		return
	}
	update.args = append(update.args, field.Value)
	placeholder := len(update.args)
	update.assignments = append(update.assignments, fmt.Sprintf("%s = $%d", column, placeholder))
}

// IsEmpty reports whether no field was added, i.e. there is nothing to update.
func (b *UpdateBuilder) IsEmpty() bool {
	return len(b.assignments) == 0
}

func (b *UpdateBuilder) Args() []any {
	return b.args
}

// SQL returns the UPDATE statement and always refreshes updated_at.
// table must be a constant or a sanitized identifier. returning may be empty.
func (b *UpdateBuilder) SQL(table, returning string) string {
	assignments := make([]string, 0, len(b.assignments)+1)
	assignments = append(assignments, b.assignments...)
	assignments = append(assignments, "updated_at = CURRENT_TIMESTAMP")

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE id = $1", table, strings.Join(assignments, ", "))
	if returning != "" {
		sql += " RETURNING " + returning
	}
	return sql
}
