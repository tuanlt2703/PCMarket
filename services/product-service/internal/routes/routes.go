package routes

import (
	"github.com/gin-gonic/gin"

	"product-service/internal/handlers"
)

func RegisterProductRoutes(router *gin.Engine, productHandler *handlers.ProductHandler) {

	products := router.Group("/products")
	{
		products.POST("", productHandler.Create)
		products.GET("", productHandler.GetAll)
		products.GET("/:id", productHandler.GetByID)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}
}
