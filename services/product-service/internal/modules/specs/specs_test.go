package specs_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"product-service/internal/modules/specs"
	"product-service/internal/modules/specs/cpu"
	"product-service/internal/platform/apperror"
)

// modelResolverStub says every model has the given component type code.
type modelResolverStub struct{ code string }

func (m modelResolverStub) ComponentTypeCode(context.Context, int64) (string, error) {
	return m.code, nil
}

// gpuInput is an Input for another component type, to test the type check.
type gpuInput struct{}

func (gpuInput) ComponentTypeCode() string { return "gpu" }

type insideTransactionKey struct{}

// transactorStub marks ctx so the component can check it runs inside the "transaction".
type transactorStub struct {
	calls     int
	lastError error
}

func (tx *transactorStub) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx.calls++
	tx.lastError = fn(context.WithValue(ctx, insideTransactionKey{}, true))
	return tx.lastError
}

// cpuComponentStub uses the real CPU request decoder but records Upsert
// calls instead of writing to a database.
type cpuComponentStub struct {
	upsertInTransaction bool
	upsertInput         specs.Input
	err                 error
}

func (*cpuComponentStub) ComponentTypeCode() string { return cpu.ComponentTypeCode }

func (*cpuComponentStub) DecodeRequest(body []byte) (specs.Input, error) {
	patch, err := cpu.DecodeUpsertRequest(body)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

func (*cpuComponentStub) Get(context.Context, int64) (any, error)       { return nil, nil }
func (*cpuComponentStub) GetPublic(context.Context, int64) (any, error) { return nil, nil }
func (*cpuComponentStub) Delete(context.Context, int64) error           { return nil }

func (c *cpuComponentStub) Upsert(ctx context.Context, modelID int64, input specs.Input) (any, error) {
	c.upsertInput = input
	c.upsertInTransaction, _ = ctx.Value(insideTransactionKey{}).(bool)
	patch, ok := input.(cpu.Patch)
	if !ok {
		return nil, apperror.Invalid("expected typed CPU patch")
	}
	return cpu.Spec{ModelID: modelID, Attributes: patch.Apply(cpu.Attributes{})}, c.err
}

func TestHandlerDecodesBeforeCallingService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	component := &cpuComponentStub{}
	tx := &transactorStub{}
	service := specs.NewService(tx, modelResolverStub{"cpu"}, specs.NewRegistry(component))
	handler := specs.NewHandler(service)
	router := gin.New()
	router.PUT("/models/:id/specs", handler.Upsert)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/models/42/specs",
		strings.NewReader(`{"core_count":8,"thread_count":16,"memory_type_ids":[1,2]}`)))
	if w.Code != http.StatusOK || !component.upsertInTransaction || tx.calls != 1 {
		t.Fatalf("upsert = %d %s, tx=%+v, component=%+v", w.Code, w.Body, tx, component)
	}
	patch, ok := component.upsertInput.(cpu.Patch)
	if !ok || patch.CoreCount.Value == nil || *patch.CoreCount.Value != 8 || len(*patch.MemoryTypeIDs.Value) != 2 {
		t.Fatalf("component did not receive a typed patch: %+v", component.upsertInput)
	}

	// Bad bodies must be rejected before a transaction is opened.
	for body, want := range map[string]int{
		`{"unknown":1}`:    http.StatusBadRequest,
		`null`:             http.StatusBadRequest,
		`{} {}`:            http.StatusBadRequest,
		`{"core_count":0}`: http.StatusUnprocessableEntity,
	} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/models/42/specs", strings.NewReader(body)))
		if w.Code != want || tx.calls != 1 {
			t.Errorf("%s: status=%d want %d, tx calls=%d", body, w.Code, want, tx.calls)
		}
	}
}

func TestServiceChecksModelTypeAndPropagatesTransactionError(t *testing.T) {
	component := &cpuComponentStub{}
	registry := specs.NewRegistry(component)
	tx := &transactorStub{}
	service := specs.NewService(tx, modelResolverStub{"cpu"}, registry)

	if _, err := service.Upsert(context.Background(), 1, nil); !errors.Is(err, apperror.ErrInvalid) || component.upsertInput != nil {
		t.Fatalf("nil input must not reach component: %v", err)
	}
	if _, err := service.Upsert(context.Background(), 1, gpuInput{}); !errors.Is(err, apperror.ErrInvalid) || component.upsertInput != nil {
		t.Fatalf("mismatched input must not reach component: %v", err)
	}

	expected := errors.New("junction insert failed")
	component.err = expected
	if _, err := service.Upsert(context.Background(), 1, cpu.Patch{}); !errors.Is(err, expected) || !errors.Is(tx.lastError, expected) {
		t.Fatalf("transaction error lost: returned=%v, transaction=%v", err, tx.lastError)
	}

	unsupported := specs.NewService(tx, modelResolverStub{"gpu"}, registry)
	if _, err := unsupported.Get(context.Background(), 1); !errors.Is(err, apperror.ErrInvalid) {
		t.Fatalf("unsupported component type: %v", err)
	}
}
