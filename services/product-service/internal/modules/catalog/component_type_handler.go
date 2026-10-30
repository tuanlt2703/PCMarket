package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/httpx"
)

type ComponentTypeHandler struct {
	svc *ComponentTypeService
}

func NewComponentTypeHandler(svc *ComponentTypeService) *ComponentTypeHandler {
	return &ComponentTypeHandler{svc: svc}
}

func (h *ComponentTypeHandler) List(c *gin.Context) {
	active, ok := httpx.QueryBool(c, "is_active")
	if !ok {
		return
	}
	page := httpx.ParsePage(c)
	items, total, err := h.svc.List(c.Request.Context(), ComponentTypeFilter{Query: c.Query("q"), IsActive: active}, page)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, httpx.NewPageResult(items, page, total))
}

func (h *ComponentTypeHandler) Get(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	ct, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (h *ComponentTypeHandler) Create(c *gin.Context) {
	var in CreateComponentTypeRequest
	if !httpx.BindJSON(c, &in) {
		return
	}
	ct, err := h.svc.Create(c.Request.Context(), in.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, ct)
}

func (h *ComponentTypeHandler) Update(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	var in UpdateComponentTypeRequest
	if !httpx.BindJSON(c, &in) {
		return
	}
	ct, err := h.svc.Update(c.Request.Context(), id, in.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (h *ComponentTypeHandler) Delete(c *gin.Context) {
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
