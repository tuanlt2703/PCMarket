package app

import (
	"net/http"
	"product-service/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

func NewApp(db *infrastructure.DB) *gin.Engine {
	router := gin.Default()

	v1Api := router.Group("/api/v1")
	adminGroup := v1Api.Group("/admin")

	//temporary APIs to test server works
	v1Api.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello from public API") })
	adminGroup.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello from private admin API") })

	return router
}
