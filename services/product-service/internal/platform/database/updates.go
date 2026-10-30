package database

import (
	"fmt"
	"strings"

	"product-service/internal/platform/optional"
)

// Updates builds a partial UPDATE from optional fields: only fields present in
// the request are written, and an explicit null writes NULL. $1 is the row key.
type Updates struct {
	sets []string
	args []any
}

func NewUpdates(key any) *Updates { return &Updates{args: []any{key}} }

// Set adds "column = value" when the field is present. column must be a
// compile-time constant, never user input.
func Set[T any](u *Updates, column string, f optional.Field[T]) {
	if !f.Set {
		return
	}
	u.args = append(u.args, f.Value)
	u.sets = append(u.sets, fmt.Sprintf("%s = $%d", column, len(u.args)))
}

func (u *Updates) Empty() bool { return len(u.sets) == 0 }

func (u *Updates) Args() []any { return u.args }

// SQL returns the UPDATE statement; it also refreshes updated_at. table must
// be a sanitized identifier.
func (u *Updates) SQL(table, returning string) string {
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE id = $1",
		table, strings.Join(append(u.sets[:len(u.sets):len(u.sets)], "updated_at = CURRENT_TIMESTAMP"), ", "))
	if returning != "" {
		sql += " RETURNING " + returning
	}
	return sql
}
