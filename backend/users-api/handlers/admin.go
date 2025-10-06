package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"users-api/config"
	"users-api/models"
	"users-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateRootRequest estructura para crear usuario root
type CreateRootRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	SecretKey string `json:"secret_key" binding:"required"` // Clave secreta para crear root
}

// CreateUserRequest estructura para crear usuario (admin)
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	Role      string `json:"role" binding:"required,oneof=user moderator admin"`
}

// UpdateUserRoleRequest estructura para actualizar rol de usuario
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user moderator admin super_admin"`
}

// UpdateUserStatusRequest estructura para activar/desactivar usuario
type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// CreateRoot crea el primer usuario root del sistema
func CreateRoot(c *gin.Context) {
	var req CreateRootRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	// Verificar clave secreta (en producción esto debería ser más seguro)
	expectedSecret := "SPORTS_PLATFORM_ROOT_2024"
	if req.SecretKey != expectedSecret {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Invalid secret key",
			"message": "Secret key required to create root user",
		})
		return
	}

	db := config.GetDB()

	// Verificar si ya existe un usuario root
	var existingRoot models.User
	if err := db.Where("role = ?", "root").First(&existingRoot).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Root user already exists",
			"message": "A root user already exists in the system",
			"existing_root": gin.H{
				"id":       existingRoot.ID,
				"username": existingRoot.Username,
				"email":    existingRoot.Email,
			},
		})
		return
	}

	// Verificar si el email ya existe
	var existingUser models.User
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

	// Crear usuario root
	rootUser := models.User{
		Email:         req.Email,
		Username:      req.Username,
		Password:      hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          "root",
		EmailVerified: true, // Root user se considera verificado automáticamente
		IsActive:      true,
	}

	if err := db.Create(&rootUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create root user",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Root user created successfully",
		"data": gin.H{
			"id":         rootUser.ID,
			"username":   rootUser.Username,
			"email":      rootUser.Email,
			"first_name": rootUser.FirstName,
			"last_name":  rootUser.LastName,
			"role":       rootUser.Role,
			"created_at": rootUser.CreatedAt,
		},
		"warning": "This is the only way to create a root user. Store these credentials safely!",
	})
}

// CreateUser crea un nuevo usuario (solo admins)
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	db := config.GetDB()

	// Verificar si el email ya existe
	var existingUser models.User
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

	// Crear usuario
	newUser := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    newUser.ToResponse(),
	})
}

// ListAllUsers lista todos los usuarios con información completa (solo admins)
func ListAllUsers(c *gin.Context) {
	// Parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Filtros
	role := c.Query("role")
	status := c.Query("status") // active, inactive, all
	search := c.Query("search") // buscar por email, username, nombre

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	db := config.GetDB()
	query := db.Model(&models.User{})

	// Aplicar filtros
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if status == "active" {
		query = query.Where("is_active = ?", true)
	} else if status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(email) LIKE ? OR LOWER(username) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Contar total
	query.Count(&total)

	// Obtener usuarios con paginación
	if err := query.
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

	// Convertir a response format (con información completa para admins)
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
			"filters": gin.H{
				"role":   role,
				"status": status,
				"search": search,
			},
		},
	})
}

// UpdateUserRole actualiza el rol de un usuario (solo admins)
func UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	db := config.GetDB()

	var user models.User
	if err := db.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User not found",
		})
		return
	}

	// No permitir cambiar el rol de usuarios root (solo root puede hacerlo)
	currentUserRole, exists := c.Get("userRole")
	if exists && user.Role == "root" && currentUserRole != "root" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Cannot modify root user",
			"message": "Only root users can modify other root users",
		})
		return
	}

	// Actualizar rol
	user.Role = req.Role
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user role",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"data":    user.ToResponse(),
	})
}

// UpdateUserStatus activa/desactiva un usuario (solo admins)
func UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	db := config.GetDB()

	var user models.User
	if err := db.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User not found",
		})
		return
	}

	// No permitir desactivar usuarios root
	if user.Role == "root" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Cannot deactivate root user",
			"message": "Root users cannot be deactivated",
		})
		return
	}

	// Actualizar estado
	user.IsActive = req.IsActive
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user status",
			"message": "Database error occurred",
		})
		return
	}

	statusText := "activated"
	if !req.IsActive {
		statusText = "deactivated"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User " + statusText + " successfully",
		"data":    user.ToResponse(),
	})
}

// DeleteUser elimina un usuario (solo root)
func DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	db := config.GetDB()

	var user models.User
	if err := db.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "User not found",
		})
		return
	}

	// No permitir eliminar usuarios root
	if user.Role == "root" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Cannot delete root user",
			"message": "Root users cannot be deleted",
		})
		return
	}

	// Eliminar usuario (soft delete sería mejor, pero por ahora hard delete)
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"deleted_user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetSystemStats obtiene estadísticas del sistema (solo admins)
func GetSystemStats(c *gin.Context) {
	db := config.GetDB()

	// Contar usuarios por rol
	var stats struct {
		TotalUsers    int64 `json:"total_users"`
		ActiveUsers   int64 `json:"active_users"`
		InactiveUsers int64 `json:"inactive_users"`
		RootUsers     int64 `json:"root_users"`
		AdminUsers    int64 `json:"admin_users"`
		ModeratorUsers int64 `json:"moderator_users"`
		RegularUsers  int64 `json:"regular_users"`
	}

	db.Model(&models.User{}).Count(&stats.TotalUsers)
	db.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)
	db.Model(&models.User{}).Where("is_active = ?", false).Count(&stats.InactiveUsers)
	db.Model(&models.User{}).Where("role = ?", "root").Count(&stats.RootUsers)
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.AdminUsers)
	db.Model(&models.User{}).Where("role = ?", "moderator").Count(&stats.ModeratorUsers)
	db.Model(&models.User{}).Where("role = ?", "user").Count(&stats.RegularUsers)

	c.JSON(http.StatusOK, gin.H{
		"message": "System statistics retrieved successfully",
		"data":    stats,
	})
}
