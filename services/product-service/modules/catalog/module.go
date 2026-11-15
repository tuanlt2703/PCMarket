package catalog

import (
	"product-service/internal/infrastructure/database"
	"product-service/modules/catalog/componenttype"
	"product-service/modules/catalog/manufacturer"
	"product-service/modules/catalog/model"

	"github.com/gin-gonic/gin"
)

type CatalogModule struct {
	ctHandler *componenttype.ComponentTypeHandler
	mfHandler *manufacturer.ManufacturerHandler
	mdHandler *model.ModelHandler
}

func NewModule(db *database.DB) *CatalogModule {
	ctRepo := componenttype.NewComponentTypeRepository(db)
	ctSvc := componenttype.NewComponentTypeService(ctRepo)

	mfRepo := manufacturer.NewManufacturerRepository(db)
	mfSvc := manufacturer.NewManufacturerService(mfRepo)

	mdRepo := model.NewModelRepository(db)
	mdSvc := model.NewModelService(mdRepo)

	return &CatalogModule{
		ctHandler: componenttype.NewComponentTypeHandler(ctSvc),
		mfHandler: manufacturer.NewManufacturerHandler(mfSvc),
		mdHandler: model.NewModelHandler(mdSvc),
	}
}

func (cat *CatalogModule) RegisterRoutes(admin, public *gin.RouterGroup) {
	ct := admin.Group("/component-types")
	ct.GET("", cat.ctHandler.GetAll)
	ct.GET("/:id", cat.ctHandler.GetByID)
	ct.POST("", cat.ctHandler.Create)
	ct.PUT("/:id", cat.ctHandler.Update)
	ct.DELETE("/:id", cat.ctHandler.Delete)

	mf := admin.Group("/manufacturers")
	mf.GET("", cat.mfHandler.GetAll)
	mf.GET("/:id", cat.mfHandler.GetByID)
	mf.POST("", cat.mfHandler.Create)
	mf.PUT("/:id", cat.mfHandler.Update)
	mf.DELETE("/:id", cat.mfHandler.Delete)

	models := admin.Group("/models")
	models.GET("", cat.mdHandler.GetAll)
	models.GET("/:id", cat.mdHandler.GetByID)
	models.POST("", cat.mdHandler.Create)
	models.PUT("/:id", cat.mdHandler.Update)
	models.DELETE("/:id", cat.mdHandler.Delete)

}
