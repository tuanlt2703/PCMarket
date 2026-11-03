package specs

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

func (h *Handler) Get(c *gin.Context) {
	modelID, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	spec, err := h.svc.Get(c.Request.Context(), modelID)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, spec)
}

// Upsert handles PUT /models/:id/specs. The body format depends on the
// model's component type, so the body is decoded only after the type is known.
func (h *Handler) Upsert(c *gin.Context) {
	modelID, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	body, err := httpx.ReadBody(c)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	component, err := h.svc.ComponentFor(c.Request.Context(), modelID)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	input, err := component.DecodeRequest(body)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	spec, err := h.svc.Upsert(c.Request.Context(), modelID, input)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, spec)
}

func (h *Handler) Delete(c *gin.Context) {
	modelID, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), modelID); err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
