package specs

import (
	"context"
	"fmt"

	"product-service/internal/platform/apperror"
)

type Input interface {
	ComponentTypeCode() string
}

type ComponentSpec interface {
	ComponentTypeCode() string
	Get(ctx context.Context, modelID int64) (any, error)
	// GetPublic is the storefront view; it must not expose inactive lookup values.
	GetPublic(ctx context.Context, modelID int64) (any, error)
	// Upsert applies a partial write: fields absent from input keep their stored value.
	Upsert(ctx context.Context, modelID int64, input Input) (any, error)
	Delete(ctx context.Context, modelID int64) error
}

// Decode runs at the HTTP boundary; services only receive typed inputs.
type Registration struct {
	Service ComponentSpec
	Decode  func(body []byte) (Input, error)
}

type Registry struct {
	byCode map[string]Registration
}

func NewRegistry(specs ...Registration) *Registry {
	r := &Registry{byCode: make(map[string]Registration, len(specs))}
	for _, s := range specs {
		if s.Service == nil || s.Decode == nil {
			panic("specs: service and request decoder are required")
		}
		code := s.Service.ComponentTypeCode()
		if _, dup := r.byCode[code]; dup {
			panic(fmt.Sprintf("specs: duplicate component type %q", code))
		}
		r.byCode[code] = s
	}
	return r
}

func (r *Registry) Get(code string) (Registration, error) {
	s, ok := r.byCode[code]
	if !ok {
		return Registration{}, apperror.Invalid(fmt.Sprintf("specifications for component type %q are not supported", code))
	}
	return s, nil
}
