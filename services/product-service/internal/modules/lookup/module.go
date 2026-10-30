package lookup

import (
	"github.com/gin-gonic/gin"

	"product-service/internal/platform/database"
)

type Module struct {
	handler *Handler
}

func NewModule(db *database.DB) *Module {
	return &Module{handler: NewHandler(NewService(NewRepository(db), NewRegistry()))}
}

// RegisterRoutes mounts:
//
//	admin:  GET /lookups (available types), CRUD /lookups/:type[/:id]
//	public: GET /lookups?types=a,b (active values grouped by type), GET /lookups/:type
func (m *Module) RegisterRoutes(admin, public *gin.RouterGroup) {
	h := m.handler

	admin.GET("/lookups", h.ListTypes)
	g := admin.Group("/lookups/:type")
	g.GET("", h.List)
	g.POST("", h.Create)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)

	public.GET("/lookups", h.PublicBatch)
	public.GET("/lookups/:type", h.PublicList)
}
