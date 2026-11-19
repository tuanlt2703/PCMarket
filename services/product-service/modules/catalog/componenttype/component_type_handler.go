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
	ctList, err := handler.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ctList)
}

func (handler *ComponentTypeHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ct, err := handler.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (handler *ComponentTypeHandler) Create(c *gin.Context) {
	var createComponentType CreateComponentType
	err := c.BindJSON(&createComponentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ct, err := handler.svc.Create(createComponentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ct)
}

func (handler *ComponentTypeHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateComponentType UpdateComponentType
	err = c.BindJSON(&updateComponentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ct, err := handler.svc.Update(id, updateComponentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (handler *ComponentTypeHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = handler.svc.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}
