package discovery

import (
	"github.com/gin-gonic/gin"

	"product-service/internal/platform/database"
)

type Module struct {
	handler *Handler
}

func NewModule(db *database.DB, specs SpecReader, offers OfferReader) *Module {
	return &Module{handler: NewHandler(NewService(NewQueryRepository(db), specs, offers))}
}

// RegisterRoutes mounts the public storefront API:
//
//	GET /products?q=&component_type=cpu&manufacturer=intel&sort=name|newest&page=&page_size=
//	GET /products/:id
func (m *Module) RegisterRoutes(_, public *gin.RouterGroup) {
	public.GET("/products", m.handler.Search)
	public.GET("/products/:id", m.handler.Get)
}
