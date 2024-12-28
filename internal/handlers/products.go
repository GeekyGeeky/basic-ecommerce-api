package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/GeekyGeeky/basic-ecommerce-api/internal/models"
)

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO products (name, description, price) VALUES (?, ?, ?)", product.Name, product.Description, product.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

func GetProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			products = append(products, p)
		}
		c.JSON(http.StatusOK, products)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = db.Exec("UPDATE products SET name = ?, description = ?, price = ? WHERE id = ?", product.Name, product.Description, product.Price, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		_, err = db.Exec("DELETE FROM products WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
	}
}
