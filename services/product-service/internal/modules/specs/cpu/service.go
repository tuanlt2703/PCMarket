package cpu

import (
	"context"
	"errors"

	"product-service/internal/platform/apperror"
)

type SpecRepository interface {
	Get(context.Context, int64) (Spec, error)
	// GetPublic hides references to inactive lookup values.
	GetPublic(context.Context, int64) (Spec, error)
	// GetForUpdate locks the spec row for the rest of the transaction.
	GetForUpdate(context.Context, int64) (Spec, error)
	Upsert(context.Context, int64, Attributes) error
	Delete(context.Context, int64) error
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

// Upsert merges the patch into the stored specification (empty on first write),
// validates the result and saves it. The parent Specs service owns the transaction.
func (s *Service) Upsert(ctx context.Context, modelID int64, patch Patch) (Spec, error) {
	current, err := s.repo.GetForUpdate(ctx, modelID)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return Spec{}, err
	}
	merged := patch.Apply(current.Attributes)
	if err := validate(merged); err != nil {
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

func validate(f Attributes) error {
	if f.CoreCount != nil && f.ThreadCount != nil && *f.ThreadCount < *f.CoreCount {
		return apperror.Invalid("thread_count must be greater than or equal to core_count")
	}
	if f.BaseClockMHz != nil && f.MaxBoostClockMHz != nil && *f.MaxBoostClockMHz < *f.BaseClockMHz {
		return apperror.Invalid("max_boost_clock_mhz must be greater than or equal to base_clock_mhz")
	}
	if f.IntegratedGPUChipsetID != nil && (f.HasIntegratedGraphics == nil || !*f.HasIntegratedGraphics) {
		return apperror.Invalid("integrated_gpu_chipset_id requires has_integrated_graphics = true")
	}
	return nil
}
