package app

import (
	"net/http"
	"product-service/internal/infrastructure"
	"product-service/modules/catalog"

	"github.com/gin-gonic/gin"
)

func NewApp(db *infrastructure.DB) *gin.Engine {
	router := gin.Default()

	publicAPIV1 := router.Group("/api/v1")
	adminAPIV1 := publicAPIV1.Group("/admin")

	catalogModule := catalog.NewModule(db)

	//temporary APIs to test server works
	publicAPIV1.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello from public API") })
	adminAPIV1.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello from private admin API") })

	catalogModule.RegisterRoutes(adminAPIV1, publicAPIV1)

	return router
}
