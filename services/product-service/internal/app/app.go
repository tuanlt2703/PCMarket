package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"product-service/internal/modules/catalog"
	"product-service/internal/modules/discovery"
	"product-service/internal/modules/lookup"
	"product-service/internal/modules/specs"
	"product-service/internal/modules/specs/cpu"
	"product-service/internal/platform/database"
)

// Module is implemented by every business module. admin routes are for the
// Admin Website; public routes are for the Customer Website.
type Module interface {
	RegisterRoutes(admin, public *gin.RouterGroup)
}

func NewRouter(db *database.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	catalogModule := catalog.NewModule(db)
	specsModule := specs.NewModule(db, catalogModule.Models,
		cpu.New(db),
		// Register gpu.New(db), motherboard.New(db), ... here as they are implemented.
	)

	modules := []Module{
		catalogModule,
		specsModule,
		lookup.NewModule(db),
		discovery.NewModule(db, specReader{specsModule.Service}, discovery.NoopOfferReader{}),
	}

	v1 := router.Group("/api/v1")
	admin := v1.Group("/admin")
	for _, m := range modules {
		m.RegisterRoutes(admin, v1)
	}
	return router
}

// specReader adapts the Specs service to Discovery's SpecReader port using the
// public view, which hides inactive lookup values.
type specReader struct{ svc *specs.Service }

func (r specReader) SpecFor(ctx context.Context, modelID int64) (any, error) {
	s, err := r.svc.GetPublic(ctx, modelID)
	if err != nil {
		return nil, err
	}
	return s, nil
}
