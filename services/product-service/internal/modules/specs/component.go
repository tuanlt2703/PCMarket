package specs

import "context"

// Input is a decoded specification write request, e.g. cpu.Patch.
// ComponentTypeCode tells which component type the input was decoded for, so
// the service can refuse a CPU input for a GPU model.
type Input interface {
	ComponentTypeCode() string
}

// Component is implemented once per component type (cpu, gpu, ...) and
// registered in internal/app. It works like an abstract base class: the specs
// module only talks to this interface and never to cpu/gpu code directly.
//
// The return values are `any` because each component has its own spec type
// (cpu.Spec, gpu.Spec, ...). They are only serialized to JSON.
type Component interface {
	// ComponentTypeCode must match component_types.code, e.g. "cpu".
	ComponentTypeCode() string

	// DecodeRequest parses the JSON body of PUT /models/:id/specs.
	// It is called by the handler; services only receive typed inputs.
	DecodeRequest(body []byte) (Input, error)

	// Get returns the admin view, including references to inactive lookup values.
	Get(ctx context.Context, modelID int64) (any, error)

	// GetPublic returns the storefront view. It must not expose inactive lookup values.
	GetPublic(ctx context.Context, modelID int64) (any, error)

	// Upsert applies a partial write: fields absent from input keep their stored value.
	// It is always called inside a transaction.
	Upsert(ctx context.Context, modelID int64, input Input) (any, error)

	Delete(ctx context.Context, modelID int64) error
}
