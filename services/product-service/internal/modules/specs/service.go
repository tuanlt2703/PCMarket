package specs

import (
	"context"

	"product-service/internal/platform/apperror"
)

// ModelTypeResolver is implemented by the Catalog module.
type ModelTypeResolver interface {
	ComponentTypeCode(ctx context.Context, modelID int64) (string, error)
}

type Transactor interface {
	InTx(context.Context, func(context.Context) error) error
}

type Service struct {
	db       Transactor
	models   ModelTypeResolver
	registry *Registry
}

func NewService(db Transactor, models ModelTypeResolver, registry *Registry) *Service {
	return &Service{db: db, models: models, registry: registry}
}

func (s *Service) Get(ctx context.Context, modelID int64) (Spec, error) {
	code, spec, err := s.resolve(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	data, err := spec.Get(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	return Spec{ComponentType: code, Data: data}, nil
}

// GetPublic returns the storefront view, without inactive lookup values.
func (s *Service) GetPublic(ctx context.Context, modelID int64) (Spec, error) {
	code, spec, err := s.resolve(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	data, err := spec.GetPublic(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	return Spec{ComponentType: code, Data: data}, nil
}

func (s *Service) ComponentTypeCode(ctx context.Context, modelID int64) (string, error) {
	return s.models.ComponentTypeCode(ctx, modelID)
}

func (s *Service) Upsert(ctx context.Context, modelID int64, input Input) (Spec, error) {
	var out Spec
	err := s.db.InTx(ctx, func(ctx context.Context) error {
		code, spec, err := s.resolve(ctx, modelID)
		if err != nil {
			return err
		}
		if input == nil || input.ComponentTypeCode() != code {
			return apperror.Invalid("specification input does not match the model's component type")
		}
		data, err := spec.Upsert(ctx, modelID, input)
		if err != nil {
			return err
		}
		out = Spec{ComponentType: code, Data: data}
		return nil
	})
	return out, err
}

func (s *Service) Delete(ctx context.Context, modelID int64) error {
	_, spec, err := s.resolve(ctx, modelID)
	if err != nil {
		return err
	}
	return spec.Delete(ctx, modelID)
}

func (s *Service) resolve(ctx context.Context, modelID int64) (string, ComponentSpec, error) {
	code, err := s.models.ComponentTypeCode(ctx, modelID)
	if err != nil {
		return "", nil, err
	}
	registration, err := s.registry.Get(code)
	return code, registration.Service, err
}
