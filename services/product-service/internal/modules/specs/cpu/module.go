package cpu

import (
	"context"

	"product-service/internal/modules/specs"
	"product-service/internal/platform/apperror"
	"product-service/internal/platform/database"
)

const ComponentTypeCode = "cpu"

func New(db *database.DB) specs.Registration {
	return specs.Registration{
		Service: component{service: NewService(NewRepository(db))},
		Decode:  decodeRequest,
	}
}

// component adapts the typed CPU service to the multi-type specification registry.
type component struct {
	service *Service
}

func (component) ComponentTypeCode() string { return ComponentTypeCode }

func (c component) Get(ctx context.Context, modelID int64) (any, error) {
	return c.service.Get(ctx, modelID)
}

func (c component) GetPublic(ctx context.Context, modelID int64) (any, error) {
	return c.service.GetPublic(ctx, modelID)
}

func (c component) Upsert(ctx context.Context, modelID int64, input specs.Input) (any, error) {
	patch, ok := input.(Patch)
	if !ok {
		return nil, apperror.Invalid("invalid input for cpu specification")
	}
	return c.service.Upsert(ctx, modelID, patch)
}

func (c component) Delete(ctx context.Context, modelID int64) error {
	return c.service.Delete(ctx, modelID)
}
