package cpu

import (
	"context"

	"product-service/internal/modules/specs"
	"product-service/internal/platform/apperror"
	"product-service/internal/platform/database"
)

// ComponentTypeCode is the component_types.code handled by this package.
const ComponentTypeCode = "cpu"

// Component plugs the CPU specification into the specs module.
// It implements specs.Component by forwarding to the typed CPU Service.
type Component struct {
	service *Service
}

// Compile-time check that *Component implements specs.Component.
var _ specs.Component = (*Component)(nil)

func NewComponent(db *database.DB) *Component {
	return &Component{service: NewService(NewRepository(db))}
}

func (c *Component) ComponentTypeCode() string {
	return ComponentTypeCode
}

func (c *Component) DecodeRequest(body []byte) (specs.Input, error) {
	patch, err := DecodeUpsertRequest(body)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

func (c *Component) Get(ctx context.Context, modelID int64) (any, error) {
	return c.service.Get(ctx, modelID)
}

func (c *Component) GetPublic(ctx context.Context, modelID int64) (any, error) {
	return c.service.GetPublic(ctx, modelID)
}

func (c *Component) Upsert(ctx context.Context, modelID int64, input specs.Input) (any, error) {
	// Type assertion, similar to dynamic_cast in C++.
	patch, ok := input.(Patch)
	if !ok {
		return nil, apperror.Invalid("invalid input for cpu specification")
	}
	return c.service.Upsert(ctx, modelID, patch)
}

func (c *Component) Delete(ctx context.Context, modelID int64) error {
	return c.service.Delete(ctx, modelID)
}
