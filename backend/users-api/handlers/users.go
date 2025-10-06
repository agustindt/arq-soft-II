package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"users-api/config"
	"users-api/middleware"
	"users-api/models"
	"users-api/utils"
)

// GetProfile obtiene el perfil del usuario autenticado
func GetProfile(c *gin.Context) {
	userID, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	var user models.User
	db := config.GetDB()
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user.ToResponse(),
	})
}

// UpdateProfileRequest estructura para actualizar perfil
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateProfile actualiza el perfil del usuario autenticado
func UpdateProfile(c *gin.Context) {
	userID, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	var user models.User
	db := config.GetDB()
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User profile not found",
		})
		return
	}

	// Actualizar los campos
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    user.ToResponse(),
	})
}

// ChangePasswordRequest estructura para cambiar contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword cambia la contraseña del usuario autenticado
func ChangePassword(c *gin.Context) {
	userID, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	var user models.User
	db := config.GetDB()
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User profile not found",
		})
		return
	}

	// Verificar la contraseña actual
	if !utils.CheckPassword(req.CurrentPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid current password",
			"message": "Current password is incorrect",
		})
		return
	}

	// Hashear la nueva contraseña
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process new password",
			"message": "Internal server error",
		})
		return
	}

	// Actualizar la contraseña
	user.Password = hashedPassword
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update password",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// GetUserByID obtiene un usuario por ID (público)
func GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	var user models.User
	db := config.GetDB()
	if err := db.Where("is_active = ?", true).First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User not found or inactive",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data":    user.ToResponse(),
	})
}

// ListUsers obtiene lista de usuarios (pública, solo usuarios activos)
func ListUsers(c *gin.Context) {
	// Parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	db := config.GetDB()

	// Contar total de usuarios activos
	db.Model(&models.User{}).Where("is_active = ?", true).Count(&total)

	// Obtener usuarios con paginación
	if err := db.Where("is_active = ?", true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"message": "Database error occurred",
		})
		return
	}

	// Convertir a response format
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data": gin.H{
			"users": userResponses,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}
