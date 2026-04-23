package handlers

import (
	"epbms/config"
	"epbms/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPerformers(c *gin.Context) {
	var performers []models.Performer

	err := config.DB.Find(&performers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get performers",
		})
		return
	}

	c.JSON(http.StatusOK, performers)
}

func GetPerformerByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid performer id",
		})
		return
	}

	var performer models.Performer

	err = config.DB.First(&performer, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "performer not found",
		})
		return
	}

	c.JSON(http.StatusOK, performer)
}

func CreatePerformer(c *gin.Context) {
	var performer models.Performer

	err := c.ShouldBindJSON(&performer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err = config.DB.Create(&performer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create performer",
		})
		return
	}

	c.JSON(http.StatusCreated, performer)
}

func UpdatePerformer(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid performer id",
		})
		return
	}

	var performer models.Performer

	err = config.DB.First(&performer, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "performer not found",
		})
		return
	}

	var updatedData models.Performer

	err = c.ShouldBindJSON(&updatedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	performer.Name = updatedData.Name
	performer.Category = updatedData.Category
	performer.Price = updatedData.Price
	performer.Description = updatedData.Description
	performer.AvailabilityStatus = updatedData.AvailabilityStatus

	err = config.DB.Save(&performer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update performer",
		})
		return
	}

	c.JSON(http.StatusOK, performer)
}

func DeletePerformer(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid performer id",
		})
		return
	}

	var performer models.Performer

	err = config.DB.First(&performer, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "performer not found",
		})
		return
	}

	err = config.DB.Delete(&performer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete performer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "performer deleted successfully",
	})
}
