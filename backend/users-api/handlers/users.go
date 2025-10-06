package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"users-api/config"
	"users-api/middleware"
	"users-api/models"
	"users-api/utils"

	"github.com/gin-gonic/gin"
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

	var req models.UpdateProfileRequest
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

	// Actualizar los campos básicos
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	
	// Actualizar campos de perfil extendido
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.BirthDate != nil {
		user.BirthDate = req.BirthDate
	}
	if req.Location != nil {
		user.Location = req.Location
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.Height != nil {
		user.Height = req.Height
	}
	if req.Weight != nil {
		user.Weight = req.Weight
	}
	if req.SportsInterests != nil {
		user.SportsInterests = req.SportsInterests
	}
	if req.FitnessLevel != nil {
		user.FitnessLevel = req.FitnessLevel
	}
	
	// Actualizar enlaces sociales
	user.SocialLinks = req.SocialLinks

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
		"data":    user.ToPublicResponse(),
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

	// Convertir a response format (público)
	var userResponses []models.PublicUserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToPublicResponse())
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

// UploadAvatar maneja la subida de avatar del usuario
func UploadAvatar(c *gin.Context) {
	userID, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	// Obtener el archivo de la request
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get uploaded file",
			"message": "Please provide an avatar file",
		})
		return
	}
	defer file.Close()

	// Validar tipo de archivo
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid file type",
			"message": "Only JPEG, PNG, GIF and WebP images are allowed",
		})
		return
	}

	// Validar tamaño (max 5MB)
	const maxSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": "Avatar must be smaller than 5MB",
		})
		return
	}

	// Crear directorio de uploads si no existe
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create upload directory",
			"message": "Internal server error",
		})
		return
	}

	// Generar nombre único para el archivo
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Crear el archivo en el servidor
	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create file",
			"message": "Internal server error",
		})
		return
	}
	defer dst.Close()

	// Copiar el contenido del archivo
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save file",
			"message": "Internal server error",
		})
		return
	}

	// Generar URL del avatar
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	// Actualizar el usuario en la base de datos
	var user models.User
	db := config.GetDB()
	if err := db.First(&user, userID).Error; err != nil {
		// Eliminar archivo si no se puede actualizar la BD
		os.Remove(filepath)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User profile not found",
		})
		return
	}

	// Eliminar avatar anterior si existe
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		oldAvatarPath := strings.TrimPrefix(*user.AvatarURL, "/")
		if _, err := os.Stat(oldAvatarPath); err == nil {
			os.Remove(oldAvatarPath)
		}
	}

	// Actualizar URL del avatar
	user.AvatarURL = &avatarURL
	if err := db.Save(&user).Error; err != nil {
		// Eliminar archivo si no se puede actualizar la BD
		os.Remove(filepath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Avatar uploaded successfully",
		"avatar_url": avatarURL,
		"data":       user.ToResponse(),
	})
}

// DeleteAvatar elimina el avatar del usuario
func DeleteAvatar(c *gin.Context) {
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

	// Eliminar archivo del avatar si existe
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		avatarPath := strings.TrimPrefix(*user.AvatarURL, "/")
		if _, err := os.Stat(avatarPath); err == nil {
			os.Remove(avatarPath)
		}
	}

	// Limpiar URL del avatar en la BD
	user.AvatarURL = nil
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Avatar deleted successfully",
		"data":    user.ToResponse(),
	})
}
