package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"arq-soft-II/backend/users-api/middleware"
	"arq-soft-II/backend/users-api/models"
	"arq-soft-II/backend/users-api/services"

	"github.com/gin-gonic/gin"
)

// UserController handles user related HTTP endpoints.
type UserController struct {
	userService *services.UserService
}

// NewUserController instantiates a UserController.
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// GetProfile returns the profile for the authenticated user.
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID, _, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	profile, err := ctrl.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve profile",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    profile,
	})
}

// UpdateProfile updates fields on the authenticated user's profile.
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID, _, _, _, exists := middleware.GetUserFromContext(c)
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

	updated, err := ctrl.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    updated,
	})
}

// ChangePassword updates the password for the authenticated user.
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID, _, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	err := ctrl.userService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User profile not found",
			})
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid current password",
				"message": "Current password is incorrect",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update password",
				"message": "Database error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// GetUserByID retrieves a public profile for an active user.
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	user, err := ctrl.userService.GetPublicUserByID(c.Request.Context(), uint(userID))
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User not found or inactive",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve user",
			"message": "Database error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// ListUsers returns active users with pagination.
func (ctrl *UserController) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := ctrl.userService.ListPublicUsers(c.Request.Context(), page, limit)
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
		},
	})
}

// UploadAvatar uploads and associates a new avatar for the authenticated user.
func (ctrl *UserController) UploadAvatar(c *gin.Context) {
	userID, _, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get uploaded file",
			"message": "Please provide an avatar file",
		})
		return
	}
	defer file.Close()

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

	const maxSize = 5 * 1024 * 1024
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": "Avatar must be smaller than 5MB",
		})
		return
	}

	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create upload directory",
			"message": "Internal server error",
		})
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), ext)
	filepathOnDisk := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filepathOnDisk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create file",
			"message": "Internal server error",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save file",
			"message": "Internal server error",
		})
		return
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	updated, previous, err := ctrl.userService.UpdateAvatar(c.Request.Context(), userID, avatarURL)
	if err != nil {
		os.Remove(filepathOnDisk)
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": "Database error occurred",
		})
		return
	}

	if previous != nil && *previous != "" {
		oldPath := strings.TrimPrefix(*previous, "/")
		if _, err := os.Stat(oldPath); err == nil {
			os.Remove(oldPath)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Avatar uploaded successfully",
		"avatar_url": avatarURL,
		"data":       updated,
	})
}

// DeleteAvatar removes the avatar association for the authenticated user.
func (ctrl *UserController) DeleteAvatar(c *gin.Context) {
	userID, _, _, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	updated, previous, err := ctrl.userService.DeleteAvatar(c.Request.Context(), userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": "User profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": "Database error occurred",
		})
		return
	}

	if previous != nil && *previous != "" {
		avatarPath := strings.TrimPrefix(*previous, "/")
		if _, err := os.Stat(avatarPath); err == nil {
			os.Remove(avatarPath)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Avatar deleted successfully",
		"data":    updated,
	})
}
