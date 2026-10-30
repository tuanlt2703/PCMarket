package cpu

import (
	"errors"
	"testing"

	"product-service/internal/platform/apperror"
)

func TestDecodeRequestAndValidate(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr error
	}{
		{"valid", `{"core_count":8,"thread_count":16,"base_clock_mhz":3600,"max_boost_clock_mhz":5000,"memory_type_ids":[1,2]}`, nil},
		{"empty object", `{}`, nil},
		{"explicit nulls", `{"core_count":null,"memory_type_ids":null}`, nil},
		{"unknown field", `{"cores":8}`, apperror.ErrMalformed},
		{"malformed json", `{`, apperror.ErrMalformed},
		{"null", `null`, apperror.ErrMalformed},
		{"array", `[]`, apperror.ErrMalformed},
		{"multiple objects", `{} {}`, apperror.ErrMalformed},
		{"trailing garbage", `{} trailing`, apperror.ErrMalformed},
		{"wrong type", `{"core_count":"8"}`, apperror.ErrMalformed},
		{"int32 overflow", `{"core_count":3000000000}`, apperror.ErrMalformed},
		{"non-positive core count", `{"core_count":0}`, apperror.ErrInvalid},
		{"non-positive link id", `{"memory_type_ids":[0]}`, apperror.ErrInvalid},
		{"threads below cores", `{"core_count":8,"thread_count":4}`, apperror.ErrInvalid},
		{"boost below base", `{"base_clock_mhz":4000,"max_boost_clock_mhz":3000}`, apperror.ErrInvalid},
		{"igpu chipset without igpu", `{"integrated_gpu_chipset_id":1}`, apperror.ErrInvalid},
		{"igpu chipset with igpu", `{"has_integrated_graphics":true,"integrated_gpu_chipset_id":1}`, nil},
		{"memory channel overflow", `{"memory_channel_count":40000}`, apperror.ErrInvalid},
		{"clock overflows NUMERIC(8,2)", `{"base_clock_mhz":10000000}`, apperror.ErrInvalid},
		{"clock at NUMERIC(8,2) max", `{"base_clock_mhz":999999.99}`, nil},
		{"cache overflows NUMERIC(10,2)", `{"l3_cache_mb":100000000}`, apperror.ErrInvalid},
		{"temperature overflows NUMERIC(5,2)", `{"max_operating_temperature_c":1000}`, apperror.ErrInvalid},
		{"negative temperature", `{"max_operating_temperature_c":-40}`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := decodeRequest([]byte(tt.body))
			if err == nil {
				err = validate(input.(Patch).Apply(Attributes{}))
			}
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatchApplyKeepsOmittedAndClearsNull(t *testing.T) {
	cores, socket := int32(8), int64(3)
	current := Attributes{CoreCount: &cores, SocketID: &socket, MemoryTypeIDs: []int64{1, 2}}
	input, err := decodeRequest([]byte(`{"socket_id":null,"memory_type_ids":[5]}`))
	if err != nil {
		t.Fatal(err)
	}
	merged := input.(Patch).Apply(current)
	if merged.CoreCount == nil || *merged.CoreCount != 8 {
		t.Fatal("omitted core_count must be kept")
	}
	if merged.SocketID != nil {
		t.Fatal("explicit null must clear socket_id")
	}
	if len(merged.MemoryTypeIDs) != 1 || merged.MemoryTypeIDs[0] != 5 {
		t.Fatalf("memory_type_ids = %v", merged.MemoryTypeIDs)
	}
	if current.SocketID == nil {
		t.Fatal("Apply must not mutate its argument")
	}
}
