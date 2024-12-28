package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/GeekyGeeky/basic-ecommerce-api/internal/auth"
	"github.com/GeekyGeeky/basic-ecommerce-api/internal/database"
	"github.com/GeekyGeeky/basic-ecommerce-api/internal/handlers"
	"github.com/GeekyGeeky/basic-ecommerce-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	/* setup auth service  */
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	authService := auth.NewAuthService(jwtKey, db)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "pong",
		})

	})

	router.POST("/auth/register", authService.Register)
	router.POST("/auth/login", authService.Login)

	protectedRoute := router.Group("/api")

	protectedRoute.Use(middleware.AuthMiddleware(authService))
	{

		/* user routes */
		protectedRoute.GET("/products", handlers.GetProducts(db))
		protectedRoute.POST("/orders", handlers.PlaceOrder(db))
		protectedRoute.GET("/orders", handlers.ListOrders(db))
		protectedRoute.PUT("/orders/:id/cancel", handlers.CancelOrder(db))

		/* admin routes */
		protectedRoute.Use(middleware.AdminMiddleware(authService))
		{
			protectedRoute.POST("/products", handlers.CreateProduct(db))
			protectedRoute.PUT("/products/:id", handlers.UpdateProduct(db))
			protectedRoute.DELETE("/products/:id", handlers.DeleteProduct(db))
			protectedRoute.PUT("/orders/:id/status", handlers.UpdateOrderStatus(db))
		}

	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router, // Use the Gin engine as the handler
	}

	// Create a context for shutdown handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Goroutine for handling shutdown signals
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill)
		<-quit // Block until a signal is received
		log.Println("Shutting down server...")

		// Attempt a graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
		log.Println("Server exiting")
	}()

	// Initial server startup message
	log.Println("Starting server on port 8080")

	// Start the server in a separate goroutine so it doesn't block
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen:", err)
	}

	// Wait for the shutdown signal or context deadline
	<-ctx.Done()
	log.Println("Server gracefully stopped")

}
