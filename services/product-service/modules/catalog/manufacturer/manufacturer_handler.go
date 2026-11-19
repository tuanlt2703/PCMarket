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
	mfList, err := handler.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mfList)
}

func (handler *ManufacturerHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mf, err := handler.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mf)
}

func (handler *ManufacturerHandler) Create(c *gin.Context) {
	var createManufacturer CreateManufacturer
	err := c.BindJSON(&createManufacturer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mf, err := handler.svc.Create(createManufacturer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mf)
}

func (handler *ManufacturerHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateManufacturer UpdateManufacturer
	err = c.BindJSON(&updateManufacturer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mf, err := handler.svc.Update(id, updateManufacturer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mf)
}

func (handler *ManufacturerHandler) Delete(c *gin.Context) {
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
