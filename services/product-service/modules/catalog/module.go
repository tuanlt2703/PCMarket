package catalog

import (
	"product-service/internal/infrastructure/database"
	"product-service/modules/catalog/componenttype"

	"github.com/gin-gonic/gin"
)

type CatalogModule struct {
	componentTypeHandler *componenttype.ComponentTypeHandler
}

func NewModule(db *database.DB) *CatalogModule {
	componenttypeRepo := componenttype.NewComponentTypeRepository(db)
	componenttypeSvc := componenttype.NewComponentTypeService(componenttypeRepo)

	return &CatalogModule{
		componentTypeHandler: componenttype.NewComponentTypeHandler(componenttypeSvc),
	}
}

func (m *CatalogModule) RegisterRoutes(admin, public *gin.RouterGroup) {
	ct := admin.Group("/component-types")
	ct.GET("", m.componentTypeHandler.GetAll)
	ct.GET("/:id", m.componentTypeHandler.GetByID)
	ct.POST("", m.componentTypeHandler.Create)
	ct.PUT("/:id", m.componentTypeHandler.Update)
	ct.DELETE("/:id", m.componentTypeHandler.Delete)
}
