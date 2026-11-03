package cpu

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/optional"
)

// specRepositoryStub keeps one stored spec in memory and records the last write.
type specRepositoryStub struct {
	stored        *Attributes // nil means no spec stored yet
	written       Attributes
	writtenModel  int64
	upsertCounter int
	err           error
}

func (r *specRepositoryStub) Get(_ context.Context, modelID int64) (Spec, error) {
	return Spec{ModelID: modelID, Attributes: r.written}, r.err
}

func (r *specRepositoryStub) GetPublic(ctx context.Context, modelID int64) (Spec, error) {
	return r.Get(ctx, modelID)
}

func (r *specRepositoryStub) GetForUpdate(_ context.Context, modelID int64) (Spec, error) {
	if r.stored == nil {
		return Spec{}, apperror.NotFound("no spec")
	}
	return Spec{ModelID: modelID, Attributes: *r.stored}, nil
}

func (r *specRepositoryStub) Upsert(_ context.Context, modelID int64, attributes Attributes) error {
	r.upsertCounter++
	r.writtenModel, r.written = modelID, attributes
	return r.err
}

func (r *specRepositoryStub) Delete(context.Context, int64) error {
	return r.err
}

func TestServiceUpsert(t *testing.T) {
	repo := &specRepositoryStub{}
	service := NewService(repo)
	patch := Patch{
		CoreCount:     optional.Of[int32](8),
		ThreadCount:   optional.Of[int32](16),
		MemoryTypeIDs: optional.Of([]int64{2}),
	}

	spec, err := service.Upsert(context.Background(), 42, patch)
	if err != nil || spec.ModelID != 42 || repo.writtenModel != 42 || repo.upsertCounter != 1 {
		t.Fatalf("upsert = %+v, %v; repo = %+v", spec, err, repo)
	}

	encoded, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]any
	if err := json.Unmarshal(encoded, &response); err != nil {
		t.Fatal(err)
	}
	if response["core_count"] != float64(8) || response["Attributes"] != nil {
		t.Fatalf("response must stay flat: %s", encoded)
	}

	patch.ThreadCount = optional.Of[int32](4)
	if _, err := service.Upsert(context.Background(), 42, patch); !errors.Is(err, apperror.ErrInvalid) || repo.upsertCounter != 1 {
		t.Fatalf("invalid input reached repository: %v, upserts=%d", err, repo.upsertCounter)
	}

	expected := errors.New("write failed")
	repo.err = expected
	patch.ThreadCount = optional.Of[int32](16)
	if _, err := service.Upsert(context.Background(), 42, patch); !errors.Is(err, expected) {
		t.Fatalf("write error lost: %v", err)
	}
}

func TestServiceMergesPatchIntoStoredSpec(t *testing.T) {
	cores, threads, socket := int32(8), int32(16), int64(3)
	repo := &specRepositoryStub{stored: &Attributes{CoreCount: &cores, ThreadCount: &threads, SocketID: &socket}}
	service := NewService(repo)

	if _, err := service.Upsert(context.Background(), 1, Patch{SocketID: optional.Null[int64]()}); err != nil {
		t.Fatal(err)
	}
	if repo.written.CoreCount == nil || *repo.written.CoreCount != 8 || repo.written.SocketID != nil {
		t.Fatalf("merged = %+v", repo.written)
	}

	// Cross-field rules are checked against the merged result, not just the patch.
	if _, err := service.Upsert(context.Background(), 1, Patch{ThreadCount: optional.Of[int32](4)}); !errors.Is(err, apperror.ErrInvalid) {
		t.Fatalf("thread_count below stored core_count: %v", err)
	}
}
