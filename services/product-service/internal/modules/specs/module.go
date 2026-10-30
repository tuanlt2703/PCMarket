package specs

import (
	"github.com/gin-gonic/gin"

	"product-service/internal/platform/database"
)

type Module struct {
	Service *Service
	handler *Handler
}

func NewModule(db *database.DB, models ModelTypeResolver, specs ...Registration) *Module {
	registry := NewRegistry(specs...)
	service := NewService(db, models, registry)
	return &Module{Service: service, handler: NewHandler(service, registry)}
}

// RegisterRoutes mounts GET/PUT/DELETE /models/:id/specs on the admin group.
// The payload shape depends on the model's component type.
func (m *Module) RegisterRoutes(admin, _ *gin.RouterGroup) {
	g := admin.Group("/models/:id/specs")
	g.GET("", m.handler.Get)
	g.PUT("", m.handler.Upsert)
	g.DELETE("", m.handler.Delete)
}
