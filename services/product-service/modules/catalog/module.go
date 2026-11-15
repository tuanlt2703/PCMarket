package catalog

import (
	"product-service/internal/infrastructure/database"
	"product-service/modules/catalog/componenttype"
	"product-service/modules/catalog/manufacturer"

	"github.com/gin-gonic/gin"
)

type CatalogModule struct {
	componentTypeHandler *componenttype.ComponentTypeHandler
	manufacturerHandler  *manufacturer.ManufacturerHandler
}

func NewModule(db *database.DB) *CatalogModule {
	componenttypeRepo := componenttype.NewComponentTypeRepository(db)
	componenttypeSvc := componenttype.NewComponentTypeService(componenttypeRepo)

	manufacturerRepo := manufacturer.NewManufacturerRepository(db)
	manufacturerSvc := manufacturer.NewManufacturerService(manufacturerRepo)

	return &CatalogModule{
		componentTypeHandler: componenttype.NewComponentTypeHandler(componenttypeSvc),
		manufacturerHandler:  manufacturer.NewManufacturerHandler(manufacturerSvc),
	}
}

func (m *CatalogModule) RegisterRoutes(admin, public *gin.RouterGroup) {
	ct := admin.Group("/component-types")
	ct.GET("", m.componentTypeHandler.GetAll)
	ct.GET("/:id", m.componentTypeHandler.GetByID)
	ct.POST("", m.componentTypeHandler.Create)
	ct.PUT("/:id", m.componentTypeHandler.Update)
	ct.DELETE("/:id", m.componentTypeHandler.Delete)

	mf := admin.Group("/manufacturers")
	mf.GET("", m.manufacturerHandler.GetAll)
	mf.GET("/:id", m.manufacturerHandler.GetByID)
	mf.POST("", m.manufacturerHandler.Create)
	mf.PUT("/:id", m.manufacturerHandler.Update)
	mf.DELETE("/:id", m.manufacturerHandler.Delete)
}
