package optional

import (
	"encoding/json"
	"testing"
)

func TestFieldDistinguishesOmittedNullAndValue(t *testing.T) {
	var req struct {
		A Field[string]  `json:"a"`
		B Field[string]  `json:"b"`
		C Field[string]  `json:"c"`
		D Field[[]int64] `json:"d"`
	}
	if err := json.Unmarshal([]byte(`{"b":null,"c":"x","d":[]}`), &req); err != nil {
		t.Fatal(err)
	}
	if req.A.Set {
		t.Error("omitted field must not be set")
	}
	if !req.B.IsNull() {
		t.Error("explicit null must be set with nil value")
	}
	if !req.C.Set || *req.C.Value != "x" {
		t.Errorf("value = %+v", req.C)
	}
	if !req.D.Set || req.D.Value == nil || len(*req.D.Value) != 0 {
		t.Errorf("empty array = %+v", req.D)
	}
}

func TestApply(t *testing.T) {
	keep, clear, set := "keep", "clear", "old"
	a, b, c := &keep, &clear, &set
	Field[string]{}.ApplyTo(&a)
	Null[string]().ApplyTo(&b)
	Of("new").ApplyTo(&c)
	if a == nil || *a != "keep" || b != nil || *c != "new" {
		t.Fatalf("a=%v b=%v c=%v", a, b, *c)
	}
	ids := []int64{1}
	Null[[]int64]().ApplyValue(&ids)
	if ids != nil {
		t.Fatal("null must reset a slice")
	}
}

func TestBlankAndTrim(t *testing.T) {
	if Blank(Field[string]{}) {
		t.Error("omitted is not blank")
	}
	if !Blank(Null[string]()) || !Blank(Of("  ")) {
		t.Error("null and whitespace are blank")
	}
	if v := TrimSpace(Of(" x ")); *v.Value != "x" {
		t.Errorf("trim = %q", *v.Value)
	}
}
