package controllers

import (
	"net/http"
	"strconv"

	"users-api/middleware"
	"users-api/repositories"
	"users-api/services"

	"github.com/gin-gonic/gin"
)

// AdminController handles administration endpoints.
type AdminController struct {
	adminService *services.AdminService
}

// NewAdminController instantiates an AdminController.
func NewAdminController(adminService *services.AdminService) *AdminController {
	return &AdminController{adminService: adminService}
}

type createRootRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	SecretKey string `json:"secret_key" binding:"required"`
}

type createUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	Role      string `json:"role" binding:"required,oneof=user moderator admin"`
}

type updateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user moderator admin super_admin"`
}

type updateUserStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// CreateRoot bootstraps the root account.
func (ctrl *AdminController) CreateRoot(c *gin.Context) {
	var req createRootRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	user, err := ctrl.adminService.CreateRoot(c.Request.Context(), services.CreateRootInput{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		SecretKey: req.SecretKey,
	})
	if err != nil {
		if rootErr, ok := err.(*services.RootExistsError); ok {
			existing := rootErr.Existing
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Root user already exists",
				"message": "A root user already exists in the system",
				"existing_root": gin.H{
					"id":       existing.ID,
					"username": existing.Username,
					"email":    existing.Email,
				},
			})
			return
		}

		switch err {
		case services.ErrInvalidSecretKey:
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Invalid secret key",
				"message": "Secret key required to create root user",
			})
		case services.ErrRootAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Root user already exists",
				"message": "A root user already exists in the system",
			})
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
				"error":   "Failed to create root user",
				"message": "Database error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Root user created successfully",
		"data":    user,
		"warning": "This is the only way to create a root user. Store these credentials safely!",
	})
}

// CreateUser creates a new account via admin.
func (ctrl *AdminController) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	user, err := ctrl.adminService.CreateUser(c.Request.Context(), services.CreateUserInput{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	})
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
		"message": "User created successfully",
		"data":    user,
	})
}

// ListAllUsers lists users with admin level detail.
func (ctrl *AdminController) ListAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := repositories.UserFilter{
		Role:   c.Query("role"),
		Status: c.Query("status"),
		Search: c.Query("search"),
		Page:   page,
		Limit:  limit,
	}

	result, err := ctrl.adminService.ListUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"message": "Database error occurred",
		})
		return
	}

	totalPages := (result.Total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data": gin.H{
			"users": result.Users,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       result.Total,
				"total_pages": totalPages,
			},
			"filters": gin.H{
				"role":   filter.Role,
				"status": filter.Status,
				"search": filter.Search,
			},
		},
	})
}

// UpdateUserRole modifies a user's role.
func (ctrl *AdminController) UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	var req updateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	role, _ := middleware.GetUserRole(c)

	updated, err := ctrl.adminService.UpdateUserRole(
		c.Request.Context(),
		services.UpdateUserRoleInput{
			UserID: uint(userID),
			Role:   req.Role,
		},
		role,
	)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User not found",
			})
		case services.ErrCannotModifyRoot:
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Cannot modify root user",
				"message": "Only root users can modify other root users",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update user role",
				"message": "Database error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"data":    updated,
	})
}

// UpdateUserStatus toggles a user's activation state.
func (ctrl *AdminController) UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	var req updateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	updated, err := ctrl.adminService.UpdateUserStatus(
		c.Request.Context(),
		services.UpdateUserStatusInput{
			UserID: uint(userID),
			Active: req.IsActive,
		},
	)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User not found",
			})
		case services.ErrCannotDeactivateRoot:
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Cannot deactivate root user",
				"message": "Root users cannot be deactivated",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update user status",
				"message": "Database error occurred",
			})
		}
		return
	}

	statusText := "activated"
	if !req.IsActive {
		statusText = "deactivated"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User " + statusText + " successfully",
		"data":    updated,
	})
}

// DeleteUser removes a user account.
func (ctrl *AdminController) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	deleted, err := ctrl.adminService.DeleteUser(c.Request.Context(), uint(userID))
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User not found",
			})
		case services.ErrCannotDeleteRoot:
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Cannot delete root user",
				"message": "Root users cannot be deleted",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to delete user",
				"message": "Database error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"deleted_user": gin.H{
			"id":       deleted.ID,
			"username": deleted.Username,
			"email":    deleted.Email,
		},
	})
}

// GetSystemStats returns aggregated statistics.
func (ctrl *AdminController) GetSystemStats(c *gin.Context) {
	stats, err := ctrl.adminService.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve system stats",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "System statistics retrieved successfully",
		"data":    stats,
	})
}
