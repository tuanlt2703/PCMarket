package optional

import (
	"encoding/json"
	"testing"
)

func TestFieldDistinguishesOmittedNullAndValue(t *testing.T) {
	var request struct {
		A Field[string]  `json:"a"`
		B Field[string]  `json:"b"`
		C Field[string]  `json:"c"`
		D Field[[]int64] `json:"d"`
	}
	if err := json.Unmarshal([]byte(`{"b":null,"c":"x","d":[]}`), &request); err != nil {
		t.Fatal(err)
	}

	if request.A.Present {
		t.Error("omitted field must not be present")
	}
	if !request.B.IsNull() {
		t.Error("explicit null must be present with a nil value")
	}
	if !request.C.Present || *request.C.Value != "x" {
		t.Errorf("value = %+v", request.C)
	}
	if !request.D.Present || request.D.Value == nil || len(*request.D.Value) != 0 {
		t.Errorf("empty array = %+v", request.D)
	}
}

func TestApplyToPointer(t *testing.T) {
	keep, clear, replace := "keep", "clear", "old"
	a, b, c := &keep, &clear, &replace

	Field[string]{}.ApplyToPointer(&a)
	Null[string]().ApplyToPointer(&b)
	Of("new").ApplyToPointer(&c)

	if a == nil || *a != "keep" {
		t.Errorf("omitted field changed the value: %v", a)
	}
	if b != nil {
		t.Errorf("null field must clear the value: %v", *b)
	}
	if c == nil || *c != "new" {
		t.Errorf("present field must replace the value: %v", c)
	}
}

func TestApplyToValue(t *testing.T) {
	ids := []int64{1}
	Null[[]int64]().ApplyToValue(&ids)
	if ids != nil {
		t.Fatal("null must reset a slice")
	}
}

func TestIsBlankAndTrimSpace(t *testing.T) {
	if IsBlank(Field[string]{}) {
		t.Error("omitted is not blank")
	}
	if !IsBlank(Null[string]()) || !IsBlank(Of("  ")) {
		t.Error("null and whitespace are blank")
	}
	if trimmed := TrimSpace(Of(" x ")); *trimmed.Value != "x" {
		t.Errorf("trim = %q", *trimmed.Value)
	}
}
