package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestNewRouterRoutes ensures all modules register without Gin route conflicts
// and that the expected endpoints exist. No database access is needed.
func TestNewRouterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(nil)

	got := map[string]bool{}
	for _, r := range router.Routes() {
		got[r.Method+" "+r.Path] = true
	}

	for _, want := range []string{
		"GET /api/v1/admin/component-types",
		"POST /api/v1/admin/manufacturers",
		"GET /api/v1/admin/models",
		"PUT /api/v1/admin/models/:id",
		"GET /api/v1/admin/models/:id/specs",
		"PUT /api/v1/admin/models/:id/specs",
		"GET /api/v1/admin/lookups",
		"GET /api/v1/admin/lookups/:type",
		"DELETE /api/v1/admin/lookups/:type/:id",
		"GET /api/v1/lookups",
		"GET /api/v1/lookups/:type",
		"GET /api/v1/products",
		"GET /api/v1/products/:id",
	} {
		if !got[want] {
			t.Errorf("missing route %s", want)
		}
	}
}

func TestRequestValidationWithoutDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(nil)

	tests := []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/api/v1/lookups/not-a-type", http.StatusNotFound},
		{http.MethodGet, "/api/v1/admin/lookups/cpu-sockets/abc", http.StatusBadRequest},
		{http.MethodGet, "/api/v1/admin/models/0", http.StatusBadRequest},
		{http.MethodGet, "/api/v1/admin/models?is_active=maybe", http.StatusBadRequest},
		{http.MethodGet, "/api/v1/products/-1", http.StatusBadRequest},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(tt.method, tt.path, nil))
		if w.Code != tt.want {
			t.Errorf("%s %s = %d, want %d (%s)", tt.method, tt.path, w.Code, tt.want, w.Body.String())
		}
	}
}
