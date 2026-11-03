// Package optional provides Field[T], a request field type for partial updates.
//
// A plain Go pointer only has two states (nil or not nil), so it cannot tell
// "the client did not send this field" apart from "the client sent null".
// Field[T] adds a Present flag to cover the third state.
package optional

import (
	"encoding/json"
	"strings"
)

// Field is a request field with three states:
//
//	field omitted    -> Present == false               : keep the stored value
//	"field": null    -> Present == true, Value == nil  : clear the stored value
//	"field": value   -> Present == true, Value != nil  : store the new value
type Field[T any] struct {
	Present bool
	Value   *T
}

// Of returns a present field holding v.
func Of[T any](v T) Field[T] {
	return Field[T]{Present: true, Value: &v}
}

// Null returns a present field holding an explicit null.
func Null[T any]() Field[T] {
	return Field[T]{Present: true, Value: nil}
}

// UnmarshalJSON is called by encoding/json only when the key exists in the
// JSON object, so reaching this method means the field is present.
func (f *Field[T]) UnmarshalJSON(data []byte) error {
	f.Present = true
	if string(data) == "null" {
		f.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	f.Value = &value
	return nil
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	if f.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*f.Value)
}

// IsNull reports whether the client explicitly sent null.
func (f Field[T]) IsNull() bool {
	return f.Present && f.Value == nil
}

// ApplyToPointer copies the field into a nullable destination when present.
// A null field sets *dst to nil.
func (f Field[T]) ApplyToPointer(dst **T) {
	if !f.Present {
		return
	}
	*dst = f.Value
}

// ApplyToValue copies the field into a non-pointer destination when present.
// A null field resets *dst to the zero value of T.
func (f Field[T]) ApplyToValue(dst *T) {
	if !f.Present {
		return
	}
	if f.Value == nil {
		var zero T
		*dst = zero
		return
	}
	*dst = *f.Value
}

// TrimSpace returns a copy of f with surrounding spaces removed from its value.
func TrimSpace(f Field[string]) Field[string] {
	if f.Value == nil {
		return f
	}
	trimmed := strings.TrimSpace(*f.Value)
	f.Value = &trimmed
	return f
}

// IsBlank reports whether a present string field is null, empty or only spaces.
// An omitted field is not blank because it will not be written.
func IsBlank(f Field[string]) bool {
	if !f.Present {
		return false
	}
	return f.Value == nil || strings.TrimSpace(*f.Value) == ""
}
