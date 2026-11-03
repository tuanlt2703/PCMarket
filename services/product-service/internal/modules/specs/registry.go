package specs

import (
	"fmt"

	"product-service/internal/platform/apperror"
)

// Registry maps a component type code ("cpu") to its Component.
type Registry struct {
	componentsByCode map[string]Component
}

// NewRegistry panics on a nil or duplicate component: both are programming
// errors that must stop the service at startup.
func NewRegistry(components ...Component) *Registry {
	registry := &Registry{componentsByCode: make(map[string]Component, len(components))}
	for _, component := range components {
		if component == nil {
			panic("specs: nil component")
		}
		code := component.ComponentTypeCode()
		if _, exists := registry.componentsByCode[code]; exists {
			panic(fmt.Sprintf("specs: duplicate component type %q", code))
		}
		registry.componentsByCode[code] = component
	}
	return registry
}

// Get returns the component for code, or an Invalid error when specifications
// for that component type are not implemented yet.
func (r *Registry) Get(code string) (Component, error) {
	component, found := r.componentsByCode[code]
	if !found {
		return nil, apperror.Invalid(fmt.Sprintf("specifications for component type %q are not supported", code))
	}
	return component, nil
}
