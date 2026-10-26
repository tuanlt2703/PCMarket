package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"product-service/internal/handlers"
	"product-service/internal/repositories"
	"product-service/internal/routes"
	"product-service/internal/services"
)

func NewApp(db *sql.DB) *gin.Engine {

	router := gin.Default()

	// Dependency Injection
	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	// Register routes
	routes.RegisterProductRoutes(router, productHandler)

	return router
}
