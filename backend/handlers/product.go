package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"minimart/config"
	"minimart/models"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	if err := config.DB.Order("id DESC").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var product models.Product

	result := config.DB.First(&product, uint(id))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}
