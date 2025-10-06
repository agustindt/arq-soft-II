package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"users-api/config"
	"users-api/models"
	"users-api/utils"
)

// LoginRequest estructura para login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest estructura para registro
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginResponse estructura para respuesta de login
type LoginResponse struct {
	Token string                `json:"token"`
	User  models.UserResponse   `json:"user"`
}

// Register maneja el registro de nuevos usuarios
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	// Verificar si el email ya existe
	var existingUser models.User
	db := config.GetDB()
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Email already exists",
			"message": "A user with this email already exists",
		})
		return
	}

	// Verificar si el username ya existe
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Username already exists",
			"message": "A user with this username already exists",
		})
		return
	}

	// Hashear la contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process password",
			"message": "Internal server error",
		})
		return
	}

	// Crear el nuevo usuario
	user := models.User{
		Email:     strings.ToLower(req.Email),
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"message": "Database error occurred",
		})
		return
	}

	// Generar JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate token",
			"message": "Authentication error",
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": LoginResponse{
			Token: token,
			User:  user.ToResponse(),
		},
	})
}

// Login maneja el inicio de sesión
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	// Buscar el usuario por email
	var user models.User
	db := config.GetDB()
	if err := db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"message": "Email or password is incorrect",
		})
		return
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Account disabled",
			"message": "Your account has been disabled",
		})
		return
	}

	// Verificar la contraseña
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"message": "Email or password is incorrect",
		})
		return
	}

	// Generar JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate token",
			"message": "Authentication error",
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": LoginResponse{
			Token: token,
			User:  user.ToResponse(),
		},
	})
}

// RefreshToken maneja la renovación de tokens
func RefreshToken(c *gin.Context) {
	// Obtener el token del header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authorization header required",
			"message": "Please provide a valid JWT token",
		})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid authorization header format",
			"message": "Authorization header must be 'Bearer <token>'",
		})
		return
	}

	tokenString := tokenParts[1]

	// Renovar el token
	newToken, err := utils.RefreshJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Failed to refresh token",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data": gin.H{
			"token": newToken,
		},
	})
}
