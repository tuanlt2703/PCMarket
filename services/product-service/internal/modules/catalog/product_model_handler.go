package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/httpx"
)

type ProductModelHandler struct {
	svc *ProductModelService
}

func NewProductModelHandler(svc *ProductModelService) *ProductModelHandler {
	return &ProductModelHandler{svc: svc}
}

func (h *ProductModelHandler) List(c *gin.Context) {
	typeID, ok := httpx.QueryID(c, "component_type_id")
	if !ok {
		return
	}
	manufacturerID, ok := httpx.QueryID(c, "manufacturer_id")
	if !ok {
		return
	}
	active, ok := httpx.QueryBool(c, "is_active")
	if !ok {
		return
	}
	page := httpx.ParsePage(c)
	items, total, err := h.svc.List(c.Request.Context(), ProductModelFilter{
		Query: c.Query("q"), ComponentTypeID: typeID, ManufacturerID: manufacturerID, IsActive: active,
	}, page)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, httpx.NewPageResult(items, page, total))
}

func (h *ProductModelHandler) Get(c *gin.Context) {
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

func (h *ProductModelHandler) Create(c *gin.Context) {
	var request CreateProductModelRequest
	if !httpx.BindJSON(c, &request) {
		return
	}
	m, err := h.svc.Create(c.Request.Context(), request.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *ProductModelHandler) Update(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	var request UpdateProductModelRequest
	if !httpx.BindJSON(c, &request) {
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, request.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *ProductModelHandler) Delete(c *gin.Context) {
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
