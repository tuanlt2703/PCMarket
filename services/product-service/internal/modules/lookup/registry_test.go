package lookup

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var (
	createTableRe = regexp.MustCompile(`(?s)CREATE TABLE (\w+) \((.*?)\n\);`)
	columnRe      = regexp.MustCompile(`(?m)^\s+(\w+)\s`)
)

// catalogOwned tables have the lookup shape but are exposed by the Catalog module.
var catalogOwned = map[string]bool{"component_types": true, "manufacturers": true}

// TestRegistryMatchesMigrations fails when a lookup-shaped table is added to or
// removed from the migrations without updating the registry.
func TestRegistryMatchesMigrations(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "..", "migrations", "*.up.sql"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no migrations found: %v", err)
	}

	want := "code,created_at,description,id,is_active,name,updated_at"
	var fromMigrations []string
	for _, f := range files {
		sql, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range createTableRe.FindAllStringSubmatch(strings.ReplaceAll(string(sql), "\r\n", "\n"), -1) {
			var cols []string
			for _, c := range columnRe.FindAllStringSubmatch(m[2], -1) {
				cols = append(cols, c[1])
			}
			sort.Strings(cols)
			if strings.Join(cols, ",") == want && !catalogOwned[m[1]] {
				fromMigrations = append(fromMigrations, m[1])
			}
		}
	}

	registered := append([]string(nil), tables...)
	sort.Strings(fromMigrations)
	sort.Strings(registered)
	if strings.Join(fromMigrations, ",") != strings.Join(registered, ",") {
		t.Fatalf("registry out of sync with migrations\nmigrations: %v\nregistry:   %v", fromMigrations, registered)
	}
}

func TestRegistrySlugs(t *testing.T) {
	r := NewRegistry()
	tbl, ok := r.Get("cpu-sockets")
	if !ok || tbl.Name != "cpu_sockets" {
		t.Fatalf("cpu-sockets resolved to %+v, %v", tbl, ok)
	}
	if _, ok := r.Get("cpu_sockets"); ok {
		t.Fatal("snake_case slug must not resolve")
	}
	if _, ok := r.Get("models"); ok {
		t.Fatal("non-lookup table must not resolve")
	}
	if len(r.Slugs()) != len(tables) {
		t.Fatalf("duplicate table names in registry")
	}
}
