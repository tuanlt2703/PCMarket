package specs

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/httpx"
)

type Handler struct {
	svc      *Service
	registry *Registry
}

func NewHandler(svc *Service, registry *Registry) *Handler {
	return &Handler{svc: svc, registry: registry}
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	spec, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, spec)
}

func (h *Handler) Upsert(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	body, err := httpx.ReadBody(c)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	code, err := h.svc.ComponentTypeCode(c.Request.Context(), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	registration, err := h.registry.Get(code)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	input, err := registration.Decode(body)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	spec, err := h.svc.Upsert(c.Request.Context(), id, input)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, spec)
}

func (h *Handler) Delete(c *gin.Context) {
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
