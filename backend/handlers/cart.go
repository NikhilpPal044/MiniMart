package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"minimart/config"
	"minimart/models"
)

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// GET /api/cart
func GetCart(c *gin.Context) {
	userID := c.GetUint("userID")

	var items []models.CartItem

	err := config.DB.
		Preload("Product").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&items).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

// POST /api/cart
func AddToCart(c *gin.Context) {
	userID := c.GetUint("userID")

	var req AddToCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Valid product_id and quantity are required",
		})
		return
	}

	var product models.Product

	err := config.DB.First(&product, req.ProductID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch product",
		})
		return
	}

	if req.Quantity > product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Requested quantity exceeds available stock",
		})
		return
	}

	var item models.CartItem

	err = config.DB.
		Where("user_id = ? AND product_id = ?", userID, req.ProductID).
		First(&item).Error

	if err == nil {
		newQuantity := item.Quantity + req.Quantity

		if newQuantity > product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Requested quantity exceeds available stock",
			})
			return
		}

		item.Quantity = newQuantity

		if err := config.DB.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not update cart",
			})
			return
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		item = models.CartItem{
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}

		if err := config.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not add product to cart",
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not check existing cart item",
		})
		return
	}

	config.DB.Preload("Product").First(&item, item.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product added to cart",
		"item":    item,
	})
}

// PATCH /api/cart/:id
func UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("userID")

	var req UpdateCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quantity must be at least 1",
		})
		return
	}

	var item models.CartItem

	err := config.DB.
		Where("id = ? AND user_id = ?", c.Param("id"), userID).
		First(&item).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cart item not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch cart item",
		})
		return
	}

	var product models.Product

	if err := config.DB.First(&product, item.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	if req.Quantity > product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Requested quantity exceeds available stock",
		})
		return
	}

	item.Quantity = req.Quantity

	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update cart item",
		})
		return
	}

	config.DB.Preload("Product").First(&item, item.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated successfully",
		"item":    item,
	})
}

// DELETE /api/cart/:id
func RemoveCartItem(c *gin.Context) {
	userID := c.GetUint("userID")

	result := config.DB.
		Where("id = ? AND user_id = ?", c.Param("id"), userID).
		Delete(&models.CartItem{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not remove cart item",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cart item not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product removed from cart",
	})
}

