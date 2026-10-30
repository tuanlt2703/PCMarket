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
	typeID, ok := httpx.QueryInt64(c, "component_type_id")
	if !ok {
		return
	}
	mfID, ok := httpx.QueryInt64(c, "manufacturer_id")
	if !ok {
		return
	}
	active, ok := httpx.QueryBool(c, "is_active")
	if !ok {
		return
	}
	page := httpx.ParsePage(c)
	items, total, err := h.svc.List(c.Request.Context(), ProductModelFilter{
		Query: c.Query("q"), ComponentTypeID: typeID, ManufacturerID: mfID, IsActive: active,
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
	var in CreateProductModelRequest
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

func (h *ProductModelHandler) Update(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	var in UpdateProductModelRequest
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
