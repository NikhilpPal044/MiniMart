package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"minimart/config"
	"minimart/models"
)

// POST /api/orders
func PlaceOrder(c *gin.Context) {
	userID := c.GetUint("userID")

	var createdOrder models.Order

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var cartItems []models.CartItem

		// Lock this users cart rows during checkout.
		if err := tx.
			Clauses().
			Where("user_id = ?", userID).
			Find(&cartItems).Error; err != nil {
			return err
		}

		if len(cartItems) == 0 {
			return errors.New("cart is empty")
		}

		order := models.Order{
			UserID: userID,
			Status: "pending",
			Total:  0,
		}

		// Validate stock and calculate the total
		// using current database prices.
		for _, item := range cartItems {
			var product models.Product

			err := tx.
				Where("id = ?", item.ProductID).
				First(&product).Error

			if err != nil {
				return errors.New("a product in your cart no longer exists")
			}

			if item.Quantity < 1 || item.Quantity > product.Stock {
				return errors.New("insufficient stock for " + product.Name)
			}

			order.Total += product.Price * float64(item.Quantity)

			order.Items = append(order.Items, models.OrderItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})
		}

		// Create the order and order items.
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Decrease stock only when sufficient stock remains.
		for _, item := range cartItems {
			result := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return errors.New("stock changed during checkout; please try again")
			}
		}

		// Clear the cart only after order creation and stock updates.
		if err := tx.
			Where("user_id = ?", userID).
			Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		createdOrder = order

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Load the persisted order and its items for the response.
	if err := config.DB.
		Preload("Items").
		First(&createdOrder, createdOrder.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Order placed, but could not load order details",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order placed successfully",
		"order":   createdOrder,
	})
}

// GET /api/orders
func GetOrders(c *gin.Context) {
	userID := c.GetUint("userID")

	var orders []models.Order

	err := config.DB.
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}
