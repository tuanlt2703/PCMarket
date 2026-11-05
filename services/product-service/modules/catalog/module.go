package catalog

import (
	"product-service/internal/infrastructure"
	"product-service/modules/catalog/componenttype"

	"github.com/gin-gonic/gin"
)

type CatalogModule struct {
	componentTypeHandler *componenttype.ComponentTypeHandler
}

func NewModule(db *infrastructure.DB) *CatalogModule {
	componenttypeRepo := componenttype.NewComponentTypeRepository(db)
	componenttypeSvc := componenttype.NewComponentTypeService(componenttypeRepo)

	return &CatalogModule{
		componentTypeHandler: componenttype.NewComponentTypeHandler(componenttypeSvc),
	}
}

func (m *CatalogModule) RegisterRoutes(admin, public *gin.RouterGroup) {
	ct := admin.Group("/component-types")
	ct.GET("", m.componentTypeHandler.GetAll)
}
