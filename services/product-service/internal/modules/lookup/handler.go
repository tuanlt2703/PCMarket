package lookup

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/httpx"
)

const maxTypesPerRequest = 20

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// --- Admin ---

func (h *Handler) ListTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"types": h.svc.ListTypes()})
}

func (h *Handler) List(c *gin.Context) {
	active, ok := httpx.QueryBool(c, "is_active")
	if !ok {
		return
	}
	page := httpx.ParsePage(c)
	items, total, err := h.svc.List(c.Request.Context(), c.Param("type"), Filter{Query: c.Query("q"), IsActive: active}, page)
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
	item, err := h.svc.Get(c.Request.Context(), c.Param("type"), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateRequest
	if !httpx.BindJSON(c, &request) {
		return
	}
	item, err := h.svc.Create(c.Request.Context(), c.Param("type"), request.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	var request UpdateRequest
	if !httpx.BindJSON(c, &request) {
		return
	}
	item, err := h.svc.Update(c.Request.Context(), c.Param("type"), id, request.Input())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), c.Param("type"), id); err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Public ---

func (h *Handler) PublicList(c *gin.Context) {
	items, err := h.svc.ListActive(c.Request.Context(), c.Param("type"))
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": toPublicItems(items)})
}

// PublicBatch lets a form or filter panel load several lookups in one request.
func (h *Handler) PublicBatch(c *gin.Context) {
	raw := strings.TrimSpace(c.Query("types"))
	if raw == "" {
		h.ListTypes(c)
		return
	}

	slugs := strings.Split(raw, ",")
	if len(slugs) > maxTypesPerRequest {
		httpx.WriteBadRequest(c, "too many lookup types requested")
		return
	}

	result := make(map[string][]PublicItem, len(slugs))
	for _, slug := range slugs {
		slug = strings.TrimSpace(slug)
		items, err := h.svc.ListActive(c.Request.Context(), slug)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}
		result[slug] = toPublicItems(items)
	}
	c.JSON(http.StatusOK, result)
}
