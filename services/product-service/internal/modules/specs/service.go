package specs

import (
	"context"

	"product-service/internal/platform/apperror"
)

// ModelTypeResolver finds the component type code of a model.
// It is implemented by catalog.ProductModelService.
type ModelTypeResolver interface {
	ComponentTypeCode(ctx context.Context, modelID int64) (string, error)
}

// Transactor runs fn inside a database transaction.
// It is implemented by *database.DB; tests use a stub.
type Transactor interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// Service finds the right Component for a model and forwards the call to it.
type Service struct {
	db       Transactor
	models   ModelTypeResolver
	registry *Registry
}

func NewService(db Transactor, models ModelTypeResolver, registry *Registry) *Service {
	return &Service{db: db, models: models, registry: registry}
}

// Get returns the admin view of a model's specification.
func (s *Service) Get(ctx context.Context, modelID int64) (Spec, error) {
	code, component, err := s.findComponent(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	data, err := component.Get(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	return Spec{ComponentType: code, Data: data}, nil
}

// GetPublic returns the storefront view, without inactive lookup values.
func (s *Service) GetPublic(ctx context.Context, modelID int64) (Spec, error) {
	code, component, err := s.findComponent(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	data, err := component.GetPublic(ctx, modelID)
	if err != nil {
		return Spec{}, err
	}
	return Spec{ComponentType: code, Data: data}, nil
}

// ComponentFor returns the Component that handles the model's component type.
// The handler uses it to decode the request body with the right decoder.
func (s *Service) ComponentFor(ctx context.Context, modelID int64) (Component, error) {
	_, component, err := s.findComponent(ctx, modelID)
	return component, err
}

// Upsert creates or partially updates a model's specification in one transaction.
func (s *Service) Upsert(ctx context.Context, modelID int64, input Input) (Spec, error) {
	var result Spec
	err := s.db.RunInTransaction(ctx, func(ctx context.Context) error {
		code, component, err := s.findComponent(ctx, modelID)
		if err != nil {
			return err
		}
		if input == nil || input.ComponentTypeCode() != code {
			return apperror.Invalid("specification input does not match the model's component type")
		}
		data, err := component.Upsert(ctx, modelID, input)
		if err != nil {
			return err
		}
		result = Spec{ComponentType: code, Data: data}
		return nil
	})
	if err != nil {
		return Spec{}, err
	}
	return result, nil
}

func (s *Service) Delete(ctx context.Context, modelID int64) error {
	_, component, err := s.findComponent(ctx, modelID)
	if err != nil {
		return err
	}
	return component.Delete(ctx, modelID)
}

// findComponent returns the model's component type code and its Component.
func (s *Service) findComponent(ctx context.Context, modelID int64) (string, Component, error) {
	code, err := s.models.ComponentTypeCode(ctx, modelID)
	if err != nil {
		return "", nil, err
	}
	component, err := s.registry.Get(code)
	if err != nil {
		return "", nil, err
	}
	return code, component, nil
}
