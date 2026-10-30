package database

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"product-service/internal/platform/apperror"
)

const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgCheckViolation      = "23514"
	pgNotNullViolation    = "23502"
	pgStringTooLong       = "22001"
	pgNumericOutOfRange   = "22003"
	pgInvalidDatetime     = "22007"
	pgDatetimeOverflow    = "22008"
	pgInvalidText         = "22P02"
)

// MapError converts driver errors into application errors. notFoundMsg is used
// when the query returned no rows.
func MapError(err error, notFoundMsg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound(notFoundMsg)
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case pgUniqueViolation:
		return apperror.Conflict("resource already exists (" + pgErr.ConstraintName + ")")
	case pgForeignKeyViolation:
		// Postgres reports deletes of referenced rows as "update or delete on table ...".
		if strings.HasPrefix(pgErr.Message, "update or delete") {
			return apperror.Conflict("resource is still referenced by other data (" + pgErr.ConstraintName + ")")
		}
		return apperror.Invalid("referenced resource does not exist (" + pgErr.ConstraintName + ")")
	case pgCheckViolation:
		return apperror.Invalid("value violates constraint " + pgErr.ConstraintName)
	case pgNotNullViolation:
		return apperror.Invalid(pgErr.ColumnName + " must not be null")
	case pgStringTooLong:
		return apperror.Invalid("value is too long")
	case pgNumericOutOfRange:
		return apperror.Invalid("numeric value is out of range")
	case pgInvalidDatetime, pgDatetimeOverflow, pgInvalidText:
		return apperror.Invalid("value has an invalid format")
	}
	return err
}
