package database

import (
	"reflect"
	"testing"

	"product-service/internal/platform/optional"
)

func TestUpdateBuilderOnlyWritesPresentFields(t *testing.T) {
	update := NewUpdateBuilder(int64(7))
	SetIfPresent(update, "code", optional.Field[string]{})
	SetIfPresent(update, "name", optional.Of("Intel"))
	SetIfPresent(update, "description", optional.Null[string]())
	if update.IsEmpty() {
		t.Fatal("update must not be empty")
	}

	want := "UPDATE manufacturers SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id"
	if got := update.SQL("manufacturers", "id"); got != want {
		t.Fatalf("sql =\n%s\nwant\n%s", got, want)
	}

	args := update.Args()
	if len(args) != 3 || args[0] != int64(7) || *args[1].(*string) != "Intel" || args[2].(*string) != nil {
		t.Fatalf("args = %#v", args)
	}

	// Calling SQL twice must give the same statement.
	want = "UPDATE manufacturers SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1"
	if got := update.SQL("manufacturers", ""); got != want {
		t.Fatalf("second SQL call changed the statement: %s", got)
	}
}

func TestUpdateBuilderEmpty(t *testing.T) {
	update := NewUpdateBuilder(1)
	SetIfPresent(update, "name", optional.Field[string]{})
	if !update.IsEmpty() || !reflect.DeepEqual(update.Args(), []any{1}) {
		t.Fatalf("expected empty update, got %#v", update)
	}
}
