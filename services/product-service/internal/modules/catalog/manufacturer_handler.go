package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/httpx"
)

type ManufacturerHandler struct {
	svc *ManufacturerService
}

func NewManufacturerHandler(svc *ManufacturerService) *ManufacturerHandler {
	return &ManufacturerHandler{svc: svc}
}

func (h *ManufacturerHandler) List(c *gin.Context) {
	active, ok := httpx.QueryBool(c, "is_active")
	if !ok {
		return
	}
	page := httpx.ParsePage(c)
	items, total, err := h.svc.List(c.Request.Context(), ManufacturerFilter{Query: c.Query("q"), IsActive: active}, page)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, httpx.NewPageResult(items, page, total))
}

func (h *ManufacturerHandler) Get(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	m, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *ManufacturerHandler) Create(c *gin.Context) {
	var in CreateManufacturerRequest
	if !httpx.BindJSON(c, &in) {
		return
	}
	m, err := h.svc.Create(c.Request.Context(), in.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *ManufacturerHandler) Update(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	var in UpdateManufacturerRequest
	if !httpx.BindJSON(c, &in) {
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, in.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *ManufacturerHandler) Delete(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
