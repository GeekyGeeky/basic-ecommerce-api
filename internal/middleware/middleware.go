package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/GeekyGeeky/basic-ecommerce-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := authService.ParseToken(tokenString)
		if err != nil || !token.Valid {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
			return
		}

		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
			return
		}
		userId := int(userIdFloat)
		c.Set("user_id", userId)
		c.Next()
	}
}

func AdminMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		// Query the database to check if the user is an admin
		// db := c.MustGet("db").(*sql.DB) // Get the database connection from context
		var isAdmin bool
		err := authService.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", userId).Scan(&isAdmin)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check admin status"})
			return
		}

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			return
		}

		c.Next()
	}
}
