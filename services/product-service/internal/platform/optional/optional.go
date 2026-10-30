// Package optional provides a field type for partial updates that tells an
// omitted JSON field apart from an explicit null.
package optional

import (
	"encoding/json"
	"strings"
)

// Field is a request field with three states:
//
//	omitted         -> Set == false: leave the stored value unchanged
//	"field": null   -> Set == true, Value == nil: clear the stored value
//	"field": value  -> Set == true, Value != nil: store the value
type Field[T any] struct {
	Set   bool
	Value *T
}

func Of[T any](v T) Field[T] { return Field[T]{Set: true, Value: &v} }

func Null[T any]() Field[T] { return Field[T]{Set: true} }

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	f.Set = true
	if string(data) == "null" {
		f.Value = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	f.Value = &v
	return nil
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	if f.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*f.Value)
}

// IsNull reports whether the request explicitly sent null.
func (f Field[T]) IsNull() bool { return f.Set && f.Value == nil }

// ApplyTo copies a present value (or null) into a nullable destination.
func (f Field[T]) ApplyTo(dst **T) {
	if f.Set {
		*dst = f.Value
	}
}

// ApplyValue copies a present value into dst; null resets dst to its zero value.
func (f Field[T]) ApplyValue(dst *T) {
	if !f.Set {
		return
	}
	if f.Value == nil {
		var zero T
		*dst = zero
		return
	}
	*dst = *f.Value
}

// ValidationValue exposes the pointer to struct validators so binding tags such
// as "omitempty,gt=0" apply to the value and skip omitted or null fields.
func (f Field[T]) ValidationValue() any { return f.Value }

// TrimSpace trims a present string value.
func TrimSpace(f Field[string]) Field[string] {
	if f.Value != nil {
		s := strings.TrimSpace(*f.Value)
		f.Value = &s
	}
	return f
}

// Blank reports whether a present string field is null or empty.
func Blank(f Field[string]) bool {
	return f.Set && (f.Value == nil || strings.TrimSpace(*f.Value) == "")
}
