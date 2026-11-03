package specs

import (
	"github.com/gin-gonic/gin"

	"product-service/internal/platform/database"
)

type Module struct {
	// Service is exposed so Discovery can read public specifications.
	Service *Service
	handler *Handler
}

// NewModule wires the specs module. components lists one Component per
// supported component type, e.g. cpu.NewComponent(db).
func NewModule(db *database.DB, models ModelTypeResolver, components ...Component) *Module {
	service := NewService(db, models, NewRegistry(components...))
	return &Module{Service: service, handler: NewHandler(service)}
}

// RegisterRoutes mounts GET/PUT/DELETE /models/:id/specs on the admin group.
// The payload shape depends on the model's component type.
func (m *Module) RegisterRoutes(admin, _ *gin.RouterGroup) {
	group := admin.Group("/models/:id/specs")
	group.GET("", m.handler.Get)
	group.PUT("", m.handler.Upsert)
	group.DELETE("", m.handler.Delete)
}
