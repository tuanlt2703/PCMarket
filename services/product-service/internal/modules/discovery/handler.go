package discovery

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Search(c *gin.Context) {
	page := httpx.ParsePage(c)
	items, total, err := h.svc.Search(c.Request.Context(), SearchFilter{
		Query:         c.Query("q"),
		ComponentType: c.Query("component_type"),
		Manufacturer:  c.Query("manufacturer"),
		Sort:          c.Query("sort"),
	}, page)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, httpx.NewPageResult(items, page, total))
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	product, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, product)
}
