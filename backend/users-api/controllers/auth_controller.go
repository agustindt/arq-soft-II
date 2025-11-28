package controllers

import (
	"net/http"
	"strings"

	"users-api/services"

	"github.com/gin-gonic/gin"
)

// AuthController wires HTTP requests to authentication services.
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController instantiates an AuthController.
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type registerRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register handles user registration.
func (ctrl *AuthController) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	result, err := ctrl.authService.Register(
		c.Request.Context(),
		services.RegisterInput{
			Email:     req.Email,
			Username:  req.Username,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	)
	if err != nil {
		switch err {
		case services.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Email already exists",
				"message": "A user with this email already exists",
			})
		case services.ErrUsernameAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Username already exists",
				"message": "A user with this username already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create user",
				"message": "Database error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": gin.H{
			"token": result.Token,
			"user":  result.User,
		},
	})
}

// Login handles user authentication.
func (ctrl *AuthController) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	result, err := ctrl.authService.Login(
		c.Request.Context(),
		services.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"message": "Email or password is incorrect",
			})
		case services.ErrAccountDisabled:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Account disabled",
				"message": "Your account has been disabled",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate token",
				"message": "Authentication error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"token": result.Token,
			"user":  result.User,
		},
	})
}

// RefreshToken issues a new token for a valid refresh request.
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
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

	newToken, err := ctrl.authService.Refresh(c.Request.Context(), tokenParts[1])
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
