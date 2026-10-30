package catalog

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
)

func assertInputMapping(t *testing.T, request, input any) {
	t.Helper()
	from, to := reflect.ValueOf(request), reflect.ValueOf(input)
	if from.NumField() != to.NumField() {
		t.Fatalf("request/input field count differs: %s vs %s", from.Type(), to.Type())
	}
	for i := 0; i < from.NumField(); i++ {
		field := from.Type().Field(i)
		value := to.FieldByName(field.Name)
		if !value.IsValid() || !reflect.DeepEqual(from.Field(i).Interface(), value.Interface()) {
			t.Errorf("field %s lost during DTO conversion", field.Name)
		}
		if domainField, ok := to.Type().FieldByName(field.Name); ok && domainField.Tag.Get("binding") != "" {
			t.Errorf("service input %s has HTTP validation tags", domainField.Name)
		}
	}
}

func decode[T any](t *testing.T, body string) T {
	t.Helper()
	var v T
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestCatalogRequestMappings(t *testing.T) {
	payload := `{"code":"cpu","name":"CPU","description":"test","is_active":false}`
	componentType := decode[CreateComponentTypeRequest](t, payload)
	assertInputMapping(t, componentType, componentType.Input())

	manufacturer := decode[CreateManufacturerRequest](t, payload)
	assertInputMapping(t, manufacturer, manufacturer.Input())
	manufacturerUpdate := decode[UpdateManufacturerRequest](t, payload)
	assertInputMapping(t, manufacturerUpdate, manufacturerUpdate.Input())

	modelPayload := `{"component_type_id":3,"manufacturer_id":4,"model_name":"Core i7","description":"test","release_date":"2023-10-17","is_active":false}`
	create := decode[CreateProductModelRequest](t, modelPayload)
	assertInputMapping(t, create, create.Input())
	if create.Input().IsActive == nil || *create.Input().IsActive {
		t.Fatal("explicit false must be preserved")
	}
	update := decode[UpdateProductModelRequest](t, `{"manufacturer_id":4,"model_name":"Core i7","description":"test","release_date":"2023-10-17","is_active":false}`)
	assertInputMapping(t, update, update.Input())

	omitted := decode[CreateManufacturerRequest](t, `{"code":"intel","name":"Intel"}`)
	if omitted.Input().IsActive != nil {
		t.Fatal("omitted is_active must retain default semantics")
	}
}

func TestUpdateRequestsDistinguishOmittedFromNull(t *testing.T) {
	in := decode[UpdateProductModelRequest](t, `{"description":null}`).Input()
	if in.ModelName.Set || in.ReleaseDate.Set || in.ManufacturerID.Set || in.IsActive.Set {
		t.Fatalf("omitted fields must stay unset: %+v", in)
	}
	if !in.Description.IsNull() {
		t.Fatal("explicit null must clear the description")
	}
}

func TestComponentTypeCodeAndNameAreImmutable(t *testing.T) {
	for _, body := range []string{`{"code":"gpu"}`, `{"name":"GPU"}`, `{"code":null,"is_active":true}`} {
		var req UpdateComponentTypeRequest
		if err := httpx.DecodeJSON([]byte(body), &req); !errors.Is(err, apperror.ErrInvalid) {
			t.Errorf("%s: got %v, want ErrInvalid", body, err)
		}
	}
	var req UpdateComponentTypeRequest
	if err := httpx.DecodeJSON([]byte(`{"description":"x","is_active":false}`), &req); err != nil {
		t.Fatal(err)
	}
	in := req.Input()
	if *in.Description.Value != "x" || *in.IsActive.Value {
		t.Fatalf("update input = %+v", in)
	}
}

func TestCreateProductModelRequiresManufacturer(t *testing.T) {
	var req CreateProductModelRequest
	err := httpx.DecodeJSON([]byte(`{"component_type_id":1,"model_name":"X"}`), &req)
	if !errors.Is(err, apperror.ErrInvalid) || err.Error() != "manufacturer_id is required" {
		t.Fatalf("got %v", err)
	}
}

func TestUpdateServicesRejectNullForRequiredColumns(t *testing.T) {
	ctx := t.Context()
	models := NewProductModelService(nil)
	for name, in := range map[string]UpdateProductModelInput{
		"blank name":        {ModelName: decode[UpdateProductModelRequest](t, `{"model_name":" "}`).ModelName},
		"null manufacturer": {ManufacturerID: decode[UpdateProductModelRequest](t, `{"manufacturer_id":null}`).ManufacturerID},
		"null is_active":    {IsActive: decode[UpdateProductModelRequest](t, `{"is_active":null}`).IsActive},
	} {
		if _, err := models.Update(ctx, 1, in); !errors.Is(err, apperror.ErrInvalid) {
			t.Errorf("model %s: got %v", name, err)
		}
	}
	manufacturers := NewManufacturerService(nil)
	if _, err := manufacturers.Update(ctx, 1, decode[UpdateManufacturerRequest](t, `{"code":null}`).Input()); !errors.Is(err, apperror.ErrInvalid) {
		t.Errorf("manufacturer null code: got %v", err)
	}
	types := NewComponentTypeService(nil)
	if _, err := types.Update(ctx, 1, UpdateComponentTypeInput{IsActive: decode[UpdateComponentTypeRequest](t, `{"is_active":null}`).IsActive}); !errors.Is(err, apperror.ErrInvalid) {
		t.Errorf("component type null is_active: got %v", err)
	}
}
