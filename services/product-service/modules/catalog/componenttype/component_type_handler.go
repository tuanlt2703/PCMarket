package componenttype

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComponentTypeHandler struct {
	svc *ComponentTypeService
}

func NewComponentTypeHandler(svc *ComponentTypeService) *ComponentTypeHandler {
	return &ComponentTypeHandler{svc: svc}
}

func (handler *ComponentTypeHandler) GetAll(c *gin.Context) {
	componentList, err := handler.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error : internal server error")
		return
	}
	c.JSON(http.StatusOK, componentList)
}
