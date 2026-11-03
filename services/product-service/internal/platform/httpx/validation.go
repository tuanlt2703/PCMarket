package httpx

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"product-service/internal/platform/apperror"
)

// Validator collects every validation problem of a request so the client can
// fix all of them at once. Typical use inside a request's Validate method:
//
//	var v httpx.Validator
//	v.RequiredString("name", r.Name, 100)
//	httpx.CheckGreaterThan(&v, "core_count", r.CoreCount.Value, 0)
//	return v.Err()
//
// All checks that take a pointer skip nil, because an omitted or null
// optional field has no value to check.
type Validator struct {
	problems []string
}

// Add records a problem message, e.g. "name is required".
func (v *Validator) Add(format string, args ...any) {
	v.problems = append(v.problems, fmt.Sprintf(format, args...))
}

// Err returns nil when there is no problem, otherwise an apperror.Invalid
// error with all messages joined by "; ".
func (v *Validator) Err() error {
	if len(v.problems) == 0 {
		return nil
	}
	return apperror.Invalid(strings.Join(v.problems, "; "))
}

// RequiredString checks a mandatory string field.
func (v *Validator) RequiredString(field, value string, maxLength int) {
	if value == "" {
		v.Add("%s is required", field)
		return
	}
	v.MaxLength(field, &value, maxLength)
}

// MaxLength checks the length in characters (not bytes), like PostgreSQL VARCHAR(n).
func (v *Validator) MaxLength(field string, value *string, maxLength int) {
	if value != nil && utf8.RuneCountInString(*value) > maxLength {
		v.Add("%s must be at most %d characters", field, maxLength)
	}
}

// RequiredID checks a mandatory foreign key. 0 means the client did not send it.
func (v *Validator) RequiredID(field string, value int64) {
	switch {
	case value == 0:
		v.Add("%s is required", field)
	case value < 0:
		v.Add("%s must be greater than 0", field)
	}
}

// PositiveIDs checks that every id in a list is greater than 0.
func (v *Validator) PositiveIDs(field string, ids *[]int64) {
	if ids == nil {
		return
	}
	for i, id := range *ids {
		if id <= 0 {
			v.Add("%s[%d] must be greater than 0", field, i)
		}
	}
}

// Date checks the format YYYY-MM-DD.
func (v *Validator) Date(field string, value *string) {
	if value == nil {
		return
	}
	if _, err := time.Parse(time.DateOnly, *value); err != nil {
		v.Add("%s must be a date in the format YYYY-MM-DD", field)
	}
}

// Number is the set of numeric types accepted by the Check* functions below.
type Number interface {
	~int32 | ~int64 | ~float64
}

// The numeric checks are functions, not Validator methods, because Go methods
// cannot have their own type parameters.

// CheckGreaterThan checks value > limit.
func CheckGreaterThan[T Number](v *Validator, field string, value *T, limit T) {
	if value != nil && *value <= limit {
		v.Add("%s must be greater than %s", field, formatNumber(limit))
	}
}

// CheckRange checks low <= value <= high.
func CheckRange[T Number](v *Validator, field string, value *T, low, high T) {
	if value == nil {
		return
	}
	if *value < low {
		v.Add("%s must be greater than or equal to %s", field, formatNumber(low))
	} else if *value > high {
		v.Add("%s must be less than or equal to %s", field, formatNumber(high))
	}
}

// CheckPositiveMax checks 0 < value <= high.
func CheckPositiveMax[T Number](v *Validator, field string, value *T, high T) {
	if value == nil {
		return
	}
	if *value <= 0 {
		v.Add("%s must be greater than 0", field)
	} else if *value > high {
		v.Add("%s must be less than or equal to %s", field, formatNumber(high))
	}
}

// formatNumber prints 99999999.99 as written instead of 9.999999999e+07.
func formatNumber[T Number](n T) string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}
