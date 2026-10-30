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

type modelResolver struct{ code string }

type gpuInput struct{}

func (gpuInput) ComponentTypeCode() string { return "gpu" }

func (m modelResolver) ComponentTypeCode(context.Context, int64) (string, error) {
	return m.code, nil
}

type txContextKey struct{}

type transactionStub struct {
	calls int
	err   error
}

func (tx *transactionStub) InTx(ctx context.Context, fn func(context.Context) error) error {
	tx.calls++
	tx.err = fn(context.WithValue(ctx, txContextKey{}, true))
	return tx.err
}

type cpuComponentStub struct {
	inTransaction bool
	input         specs.Input
	err           error
}

func (*cpuComponentStub) ComponentTypeCode() string                     { return "cpu" }
func (*cpuComponentStub) Get(context.Context, int64) (any, error)       { return nil, nil }
func (*cpuComponentStub) GetPublic(context.Context, int64) (any, error) { return nil, nil }
func (*cpuComponentStub) Delete(context.Context, int64) error           { return nil }
func (c *cpuComponentStub) Upsert(ctx context.Context, modelID int64, input specs.Input) (any, error) {
	c.input = input
	c.inTransaction, _ = ctx.Value(txContextKey{}).(bool)
	patch, ok := input.(cpu.Patch)
	if !ok {
		return nil, apperror.Invalid("expected typed CPU patch")
	}
	return cpu.Spec{ModelID: modelID, Attributes: patch.Apply(cpu.Attributes{})}, c.err
}

func TestHandlerDecodesBeforeCallingService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	component := &cpuComponentStub{}
	registration := cpu.New(nil)
	registration.Service = component
	registry := specs.NewRegistry(registration)
	tx := &transactionStub{}
	service := specs.NewService(tx, modelResolver{"cpu"}, registry)
	handler := specs.NewHandler(service, registry)
	router := gin.New()
	router.PUT("/models/:id/specs", handler.Upsert)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/models/42/specs",
		strings.NewReader(`{"core_count":8,"thread_count":16,"memory_type_ids":[1,2]}`)))
	if w.Code != http.StatusOK || !component.inTransaction || tx.calls != 1 {
		t.Fatalf("upsert = %d %s, tx=%+v, component=%+v", w.Code, w.Body, tx, component)
	}
	patch, ok := component.input.(cpu.Patch)
	if !ok || patch.CoreCount.Value == nil || *patch.CoreCount.Value != 8 || len(*patch.MemoryTypeIDs.Value) != 2 {
		t.Fatalf("service did not receive typed input: %+v", component.input)
	}
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
	registration := cpu.New(nil)
	registration.Service = component
	registry := specs.NewRegistry(registration)
	tx := &transactionStub{}
	service := specs.NewService(tx, modelResolver{"cpu"}, registry)
	if _, err := service.Upsert(context.Background(), 1, nil); !errors.Is(err, apperror.ErrInvalid) || component.input != nil {
		t.Fatalf("nil input must not reach component: %v", err)
	}
	if _, err := service.Upsert(context.Background(), 1, gpuInput{}); !errors.Is(err, apperror.ErrInvalid) || component.input != nil {
		t.Fatalf("mismatched input must not reach component: %v", err)
	}
	expected := errors.New("junction insert failed")
	component.err = expected
	if _, err := service.Upsert(context.Background(), 1, cpu.Patch{}); !errors.Is(err, expected) || !errors.Is(tx.err, expected) {
		t.Fatalf("transaction error lost: returned=%v, transaction=%v", err, tx.err)
	}
	unsupported := specs.NewService(tx, modelResolver{"gpu"}, registry)
	if _, err := unsupported.Get(context.Background(), 1); !errors.Is(err, apperror.ErrInvalid) {
		t.Fatalf("unsupported component type: %v", err)
	}
}
