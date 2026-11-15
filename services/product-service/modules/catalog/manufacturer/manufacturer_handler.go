package manufacturer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManufacturerHandler struct {
	svc *ManufacturerService
}

func NewManufacturerHandler(svc *ManufacturerService) *ManufacturerHandler {
	return &ManufacturerHandler{svc: svc}
}

func (handler *ManufacturerHandler) GetAll(c *gin.Context) {
	componentList, err := handler.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, componentList)
}

func (handler *ManufacturerHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	component, err := handler.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, component)
}

func (handler *ManufacturerHandler) Create(c *gin.Context) {
	var newComponent CreateManufacturer
	err := c.BindJSON(&newComponent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	component, err := handler.svc.Create(newComponent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, component)
}

func (handler *ManufacturerHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var updateComponent UpdateManufacturer
	err = c.BindJSON(&updateComponent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	component, err := handler.svc.Update(id, updateComponent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, component)
}

func (handler *ManufacturerHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = handler.svc.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}
