package auth

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	// "github.com/GeekyGeeky/basic-ecommerce-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtKey []byte
	DB     *sql.DB
}

func NewAuthService(jwtKey []byte, db *sql.DB) *AuthService {
	return &AuthService{jwtKey: jwtKey, DB: db}
}

func (s *AuthService) GenerateJWT(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		// "user_id":  user.ID,
		// "is_admin": user.IsAdmin,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(s.jwtKey)
	return tokenString, err
}

func (s *AuthService) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	return token, err
}

type RegisterRequest struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (h *AuthService) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	query := "INSERT INTO users (password, email) VALUES (?, ?)"
	_, err = h.DB.Exec(query, hashedPassword, req.Email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthService) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedPassword string
	var userId uint
	query := "SELECT id, password FROM users WHERE email = ?"
	err := h.DB.QueryRow(query, req.Email).Scan(&userId, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := h.GenerateJWT(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "message": "Login successful"})
}
