package database

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"product-service/internal/platform/apperror"
)

func TestMapError(t *testing.T) {
	other := errors.New("boom")
	tests := []struct {
		name string
		in   error
		want error
	}{
		{"nil", nil, nil},
		{"no rows", fmt.Errorf("wrap: %w", pgx.ErrNoRows), apperror.ErrNotFound},
		{"unique", &pgconn.PgError{Code: "23505"}, apperror.ErrConflict},
		{"fk on insert", &pgconn.PgError{Code: "23503", Message: `insert or update on table "models" violates foreign key constraint`}, apperror.ErrInvalid},
		{"fk on delete", &pgconn.PgError{Code: "23503", Message: `update or delete on table "cpu_sockets" violates foreign key constraint`}, apperror.ErrConflict},
		{"check", &pgconn.PgError{Code: "23514"}, apperror.ErrInvalid},
		{"not null", &pgconn.PgError{Code: "23502", ColumnName: "name"}, apperror.ErrInvalid},
		{"string too long", &pgconn.PgError{Code: "22001"}, apperror.ErrInvalid},
		{"numeric overflow", &pgconn.PgError{Code: "22003"}, apperror.ErrInvalid},
		{"bad date", &pgconn.PgError{Code: "22007"}, apperror.ErrInvalid},
		{"date overflow", &pgconn.PgError{Code: "22008"}, apperror.ErrInvalid},
		{"bad text", &pgconn.PgError{Code: "22P02"}, apperror.ErrInvalid},
		{"unknown", other, other},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapError(tt.in, "missing")
			if tt.want == nil {
				if got != nil {
					t.Fatalf("got %v, want nil", got)
				}
				return
			}
			if !errors.Is(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
