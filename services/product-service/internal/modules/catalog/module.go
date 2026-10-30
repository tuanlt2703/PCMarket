package catalog

import (
	"github.com/gin-gonic/gin"

	"product-service/internal/platform/database"
)

type Module struct {
	// Models is exposed so the Specs module can resolve a model's component type.
	Models *ProductModelService

	componentTypes *ComponentTypeHandler
	manufacturers  *ManufacturerHandler
	models         *ProductModelHandler
}

func NewModule(db *database.DB) *Module {
	models := NewProductModelService(NewProductModelRepository(db))
	return &Module{
		Models:         models,
		componentTypes: NewComponentTypeHandler(NewComponentTypeService(NewComponentTypeRepository(db))),
		manufacturers:  NewManufacturerHandler(NewManufacturerService(NewManufacturerRepository(db))),
		models:         NewProductModelHandler(models),
	}
}

// RegisterRoutes mounts the admin Catalog API. Storefront reads go through Discovery.
func (m *Module) RegisterRoutes(admin, _ *gin.RouterGroup) {
	ct := admin.Group("/component-types")
	ct.GET("", m.componentTypes.List)
	ct.POST("", m.componentTypes.Create)
	ct.GET("/:id", m.componentTypes.Get)
	ct.PUT("/:id", m.componentTypes.Update)
	ct.DELETE("/:id", m.componentTypes.Delete)

	mf := admin.Group("/manufacturers")
	mf.GET("", m.manufacturers.List)
	mf.POST("", m.manufacturers.Create)
	mf.GET("/:id", m.manufacturers.Get)
	mf.PUT("/:id", m.manufacturers.Update)
	mf.DELETE("/:id", m.manufacturers.Delete)

	models := admin.Group("/models")
	models.GET("", m.models.List)
	models.POST("", m.models.Create)
	models.GET("/:id", m.models.Get)
	models.PUT("/:id", m.models.Update)
	models.DELETE("/:id", m.models.Delete)
}
