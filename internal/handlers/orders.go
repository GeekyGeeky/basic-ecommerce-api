package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/GeekyGeeky/basic-ecommerce-api/internal/models"

	"github.com/gin-gonic/gin"
)

// PlaceOrder (POST /orders)
func PlaceOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from context"})
			return
		}
		order.UserID = userID.(uint)

		query := `INSERT INTO orders (user_id, product_id, status) VALUES (?, ?, ?)`
		result, err := db.Exec(query, order.UserID, order.ProductID, order.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order: " + err.Error()})
			return
		}

		orderID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last inserted ID"})
			return
		}
		order.ID = uint(orderID)

		c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully", "order": order})
	}
}

// ListOrders (GET /orders)
func ListOrders(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from context"})
			return
		}

		query := `SELECT id, product_id, status, created_at FROM orders WHERE user_id = ?`
		rows, err := db.Query(query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
			return
		}
		defer rows.Close()

		var orders []models.Order
		for rows.Next() {
			var order models.Order
			err := rows.Scan(&order.ID, &order.ProductID, &order.Status, &order.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order data"})
				return
			}
			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

// CancelOrder (PUT /orders/:id/cancel)
func CancelOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from context"})
			return
		}

		query := `UPDATE orders SET status = 'Cancelled' WHERE id = ? AND user_id = ?`
		result, err := db.Exec(query, orderID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affected rows"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or not owned by user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
	}
}

// UpdateOrderStatus (PUT /orders/:id/status)
func UpdateOrderStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		var req struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from context"})
			return
		}
		query := `UPDATE orders SET status = ? WHERE id = ? AND user_id = ?`
		result, err := db.Exec(query, req.Status, orderID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affected rows"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or not owned by user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
	}
}
