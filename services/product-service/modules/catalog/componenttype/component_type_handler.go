package componenttype

import (
	"net/http"
	"strconv"

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

func (handler *ComponentTypeHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Error: incorrect id format")
	}

	component, err := handler.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error : internal server error")
		return
	}
	c.JSON(http.StatusOK, component)
}
