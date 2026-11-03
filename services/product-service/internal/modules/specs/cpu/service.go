package cpu

import (
	"context"
	"errors"

	"product-service/internal/platform/apperror"
)

// SpecRepository is the storage used by Service. *Repository implements it
// with PostgreSQL; tests use an in-memory stub.
type SpecRepository interface {
	// Get returns the admin view, including inactive lookup references.
	Get(ctx context.Context, modelID int64) (Spec, error)
	// GetPublic hides references to inactive lookup values.
	GetPublic(ctx context.Context, modelID int64) (Spec, error)
	// GetForUpdate locks the spec row until the transaction ends.
	GetForUpdate(ctx context.Context, modelID int64) (Spec, error)
	Upsert(ctx context.Context, modelID int64, attributes Attributes) error
	Delete(ctx context.Context, modelID int64) error
}

type Service struct {
	repo SpecRepository
}

func NewService(repo SpecRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, modelID int64) (Spec, error) {
	return s.repo.Get(ctx, modelID)
}

func (s *Service) GetPublic(ctx context.Context, modelID int64) (Spec, error) {
	return s.repo.GetPublic(ctx, modelID)
}

// Upsert merges patch into the stored specification, checks the result and
// saves it. It must run inside a transaction (opened by specs.Service) so the
// row stays locked between the read and the write.
func (s *Service) Upsert(ctx context.Context, modelID int64, patch Patch) (Spec, error) {
	current, err := s.repo.GetForUpdate(ctx, modelID)
	if errors.Is(err, apperror.ErrNotFound) {
		// First write for this model: start from an empty specification.
		current = Spec{}
	} else if err != nil {
		return Spec{}, err
	}

	merged := patch.Apply(current.Attributes)
	if err := checkCrossFieldRules(merged); err != nil {
		return Spec{}, err
	}
	if err := s.repo.Upsert(ctx, modelID, merged); err != nil {
		return Spec{}, err
	}
	return s.repo.Get(ctx, modelID)
}

func (s *Service) Delete(ctx context.Context, modelID int64) error {
	return s.repo.Delete(ctx, modelID)
}

// checkCrossFieldRules checks rules that involve more than one field. They
// run on the merged result, because one side may come from the stored spec.
func checkCrossFieldRules(a Attributes) error {
	if a.CoreCount != nil && a.ThreadCount != nil && *a.ThreadCount < *a.CoreCount {
		return apperror.Invalid("thread_count must be greater than or equal to core_count")
	}
	if a.BaseClockMHz != nil && a.MaxBoostClockMHz != nil && *a.MaxBoostClockMHz < *a.BaseClockMHz {
		return apperror.Invalid("max_boost_clock_mhz must be greater than or equal to base_clock_mhz")
	}
	hasIntegratedGraphics := a.HasIntegratedGraphics != nil && *a.HasIntegratedGraphics
	if a.IntegratedGPUChipsetID != nil && !hasIntegratedGraphics {
		return apperror.Invalid("integrated_gpu_chipset_id requires has_integrated_graphics = true")
	}
	return nil
}
