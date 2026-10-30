package lookup

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/httpx"
)

type itemRepositoryStub struct {
	table  Table
	input  Input
	update UpdateInput
	filter Filter
	page   *httpx.Page
	err    error
}

func (r *itemRepositoryStub) List(_ context.Context, table Table, filter Filter, page *httpx.Page) ([]Item, int64, error) {
	r.table, r.filter, r.page = table, filter, page
	return []Item{}, 0, r.err
}

func (r *itemRepositoryStub) Get(_ context.Context, table Table, id int64) (Item, error) {
	r.table = table
	return Item{ID: id}, r.err
}

func (r *itemRepositoryStub) Create(_ context.Context, table Table, input Input) (Item, error) {
	r.table, r.input = table, input
	return Item{ID: 1, Code: input.Code, Name: input.Name}, r.err
}

func (r *itemRepositoryStub) Update(_ context.Context, table Table, id int64, input UpdateInput) (Item, error) {
	r.table, r.update = table, input
	return Item{ID: id}, r.err
}

func (r *itemRepositoryStub) Delete(_ context.Context, table Table, _ int64) error {
	r.table = table
	return r.err
}

func TestServiceResolvesLookupTable(t *testing.T) {
	ctx := context.Background()
	repo := &itemRepositoryStub{}
	service := NewService(repo, NewRegistry())
	page := httpx.Page{Page: 1, PageSize: 20}
	_, _, err := service.List(ctx, "cpu-sockets", Filter{Query: " lga "}, page)
	if err != nil || repo.table.Name != "cpu_sockets" || repo.filter.Query != "lga" {
		t.Fatalf("list = %+v, %v", repo, err)
	}
	_, err = service.ListActive(ctx, "memory-types")
	if err != nil || repo.table.Name != "memory_types" || repo.filter.IsActive == nil || !*repo.filter.IsActive || repo.page != nil {
		t.Fatalf("active list = %+v, %v", repo, err)
	}
	_, err = service.Create(ctx, "cpu-sockets", Input{Code: " lga ", Name: " LGA "})
	if err != nil || repo.input.Code != "lga" || repo.input.Name != "LGA" {
		t.Fatalf("create = %+v, %v", repo, err)
	}
	for _, slug := range []string{"models", "cpu_sockets", "cpu-sockets; DROP TABLE models"} {
		_, err := service.Get(ctx, slug, 1)
		if !errors.Is(err, apperror.ErrNotFound) {
			t.Errorf("%q: got %v, want not found", slug, err)
		}
	}
	expected := errors.New("database unavailable")
	repo.err = expected
	if _, err := service.Get(ctx, "cpu-sockets", 1); !errors.Is(err, expected) {
		t.Fatalf("repository error lost: %v", err)
	}
}

func TestHandlerUsesServiceForLookupRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &itemRepositoryStub{}
	handler := NewHandler(NewService(repo, NewRegistry()))
	router := gin.New()
	router.POST("/lookups/:type", handler.Create)
	router.PUT("/lookups/:type/:id", handler.Update)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/lookups/cpu-sockets",
		strings.NewReader(`{"code":" lga1700 ","name":"LGA1700","is_active":false}`)))
	if w.Code != http.StatusCreated || repo.table.Name != "cpu_sockets" || repo.input.Code != "lga1700" ||
		repo.input.IsActive == nil || *repo.input.IsActive {
		t.Fatalf("create = %d %s, repository = %+v", w.Code, w.Body, repo)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/lookups/cpu-sockets/7",
		strings.NewReader(`{"name":" LGA 1700 ","description":null}`)))
	u := repo.update
	if w.Code != http.StatusOK || u.Code.Set || u.IsActive.Set || !u.Description.IsNull() ||
		u.Name.Value == nil || *u.Name.Value != "LGA 1700" {
		t.Fatalf("partial update = %d %s, input = %+v", w.Code, w.Body, u)
	}

	for body, want := range map[string]int{
		`{"name":""}`:        http.StatusUnprocessableEntity,
		`{"code":null}`:      http.StatusUnprocessableEntity,
		`{"is_active":null}`: http.StatusUnprocessableEntity,
		`{"code":"` + strings.Repeat("x", 51) + `"}`: http.StatusUnprocessableEntity,
		`{"unknown":1}`:     http.StatusBadRequest,
		`{"is_active":"y"}`: http.StatusBadRequest,
	} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/lookups/cpu-sockets/7", strings.NewReader(body)))
		if w.Code != want {
			t.Errorf("%s: status = %d, want %d (%s)", body, w.Code, want, w.Body)
		}
	}
}
