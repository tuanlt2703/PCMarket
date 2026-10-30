package database

import (
	"reflect"
	"testing"

	"product-service/internal/platform/optional"
)

func TestUpdatesOnlyWritesPresentFields(t *testing.T) {
	u := NewUpdates(int64(7))
	Set(u, "code", optional.Field[string]{})
	Set(u, "name", optional.Of("Intel"))
	Set(u, "description", optional.Null[string]())
	if u.Empty() {
		t.Fatal("updates must not be empty")
	}
	want := "UPDATE manufacturers SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id"
	if got := u.SQL("manufacturers", "id"); got != want {
		t.Fatalf("sql =\n%s\nwant\n%s", got, want)
	}
	args := u.Args()
	if len(args) != 3 || args[0] != int64(7) || *args[1].(*string) != "Intel" || args[2].(*string) != nil {
		t.Fatalf("args = %#v", args)
	}
	// SQL must be repeatable without leaking updated_at into the set list.
	if u.SQL("manufacturers", "") != "UPDATE manufacturers SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1" {
		t.Fatal("second SQL call changed the statement")
	}
}

func TestUpdatesEmpty(t *testing.T) {
	u := NewUpdates(1)
	Set(u, "name", optional.Field[string]{})
	if !u.Empty() || !reflect.DeepEqual(u.Args(), []any{1}) {
		t.Fatalf("expected empty update, got %#v", u)
	}
}
