package catalog

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
)

// assertInputMapping checks that Request.Input() copies every field by name,
// so a field added to a request is not silently dropped before the service.
func assertInputMapping(t *testing.T, request, input any) {
	t.Helper()
	from, to := reflect.ValueOf(request), reflect.ValueOf(input)
	if from.NumField() != to.NumField() {
		t.Fatalf("request/input field count differs: %s vs %s", from.Type(), to.Type())
	}
	for i := 0; i < from.NumField(); i++ {
		name := from.Type().Field(i).Name
		value := to.FieldByName(name)
		if !value.IsValid() || !reflect.DeepEqual(from.Field(i).Interface(), value.Interface()) {
			t.Errorf("field %s lost during request conversion", name)
		}
	}
}

func decode[T any](t *testing.T, body string) T {
	t.Helper()
	var value T
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		t.Fatal(err)
	}
	return value
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
	input := decode[UpdateProductModelRequest](t, `{"description":null}`).Input()
	if input.ModelName.Present || input.ReleaseDate.Present || input.ManufacturerID.Present || input.IsActive.Present {
		t.Fatalf("omitted fields must stay absent: %+v", input)
	}
	if !input.Description.IsNull() {
		t.Fatal("explicit null must clear the description")
	}
}

func TestComponentTypeCodeAndNameAreImmutable(t *testing.T) {
	for _, body := range []string{`{"code":"gpu"}`, `{"name":"GPU"}`, `{"code":null,"is_active":true}`} {
		var request UpdateComponentTypeRequest
		if err := httpx.DecodeJSON([]byte(body), &request); !errors.Is(err, apperror.ErrInvalid) {
			t.Errorf("%s: got %v, want ErrInvalid", body, err)
		}
	}

	var request UpdateComponentTypeRequest
	if err := httpx.DecodeJSON([]byte(`{"description":"x","is_active":false}`), &request); err != nil {
		t.Fatal(err)
	}
	input := request.Input()
	if *input.Description.Value != "x" || *input.IsActive.Value {
		t.Fatalf("update input = %+v", input)
	}
}

func TestRequestValidationMessages(t *testing.T) {
	tests := []struct {
		name    string
		request httpx.Request
		body    string
		want    string
	}{
		{"missing manufacturer", &CreateProductModelRequest{}, `{"component_type_id":1,"model_name":"X"}`,
			"manufacturer_id is required"},
		{"negative type", &CreateProductModelRequest{}, `{"component_type_id":-1,"manufacturer_id":1,"model_name":"X"}`,
			"component_type_id must be greater than 0"},
		{"bad date", &CreateProductModelRequest{}, `{"component_type_id":1,"manufacturer_id":1,"model_name":"X","release_date":"17/10/2023"}`,
			"release_date must be a date in the format YYYY-MM-DD"},
		{"several problems", &CreateManufacturerRequest{}, `{"name":"` + strings.Repeat("é", 101) + `"}`,
			"code is required; name must be at most 100 characters"},
		{"update zero manufacturer", &UpdateProductModelRequest{}, `{"manufacturer_id":0}`,
			"manufacturer_id must be greater than 0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := httpx.DecodeJSON([]byte(tt.body), tt.request)
			if !errors.Is(err, apperror.ErrInvalid) || err.Error() != tt.want {
				t.Fatalf("got %v, want %q", err, tt.want)
			}
		})
	}
}

func TestUpdateServicesRejectNullForRequiredColumns(t *testing.T) {
	ctx := t.Context()
	models := NewProductModelService(nil)
	for name, input := range map[string]UpdateProductModelInput{
		"blank name":        {ModelName: decode[UpdateProductModelRequest](t, `{"model_name":" "}`).ModelName},
		"null manufacturer": {ManufacturerID: decode[UpdateProductModelRequest](t, `{"manufacturer_id":null}`).ManufacturerID},
		"null is_active":    {IsActive: decode[UpdateProductModelRequest](t, `{"is_active":null}`).IsActive},
	} {
		if _, err := models.Update(ctx, 1, input); !errors.Is(err, apperror.ErrInvalid) {
			t.Errorf("model %s: got %v", name, err)
		}
	}

	manufacturers := NewManufacturerService(nil)
	if _, err := manufacturers.Update(ctx, 1, decode[UpdateManufacturerRequest](t, `{"code":null}`).Input()); !errors.Is(err, apperror.ErrInvalid) {
		t.Errorf("manufacturer null code: got %v", err)
	}

	types := NewComponentTypeService(nil)
	nullIsActive := decode[UpdateComponentTypeRequest](t, `{"is_active":null}`).IsActive
	if _, err := types.Update(ctx, 1, UpdateComponentTypeInput{IsActive: nullIsActive}); !errors.Is(err, apperror.ErrInvalid) {
		t.Errorf("component type null is_active: got %v", err)
	}
}
