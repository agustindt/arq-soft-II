package controllers

import (
	"activities-api/domain"
	"activities-api/services"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ActivitiesService interface {
	List(ctx context.Context) ([]domain.Activity, error)
	ListAll(ctx context.Context) ([]domain.Activity, error)
	Create(ctx context.Context, activity domain.Activity) (domain.Activity, error)
	GetByID(ctx context.Context, id string) (domain.Activity, error)
	Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error)
	Delete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	ToggleActive(ctx context.Context, id string) (domain.Activity, error)
	GetByCategory(ctx context.Context, category string) ([]domain.Activity, error)
}

// ActivitiesController maneja las peticiones HTTP para Activities
type ActivitiesController struct {
	service ActivitiesService
}

// NewActivitiesController crea una nueva instancia del controller
func NewActivitiesController(activitiesService ActivitiesService) *ActivitiesController {
	return &ActivitiesController{
		service: activitiesService,
	}
}

// GetActivities maneja GET /activities
func (c *ActivitiesController) GetActivities(ctx *gin.Context) {
	activities, err := c.service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch activities",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"count":      len(activities),
	})
}

// GetAllActivities maneja GET /activities/all
func (c *ActivitiesController) GetAllActivities(ctx *gin.Context) {
	activities, err := c.service.ListAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch all activities",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"count":      len(activities),
	})
}

// CreateActivity maneja POST /activities - Crea una nueva actividad
func (c *ActivitiesController) CreateActivity(ctx *gin.Context) {
	var newActivity domain.Activity
	if err := ctx.ShouldBindJSON(&newActivity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Obtener user_id del contexto (seteado por el middleware)
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			newActivity.CreatedBy = uid
		}
	}

	activity, err := c.service.Create(ctx.Request.Context(), newActivity)
	if err != nil {

		if errors.Is(err, services.ErrOwnerNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Owner not found or unauthorized",
			})
			return
		}

		if errors.Is(err, services.ErrOwnerForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "You are not allowed to modify this resource",
			})
			return
		}

		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "required") {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{
			"error":   "Failed to create activity",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"activity": activity,
		"message":  "Activity created successfully",
	})
}

// GetActivityByID maneja GET /activities/:id - Obtiene actividad por ID
func (c *ActivitiesController) GetActivityByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	activity, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch activity",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activity": activity,
	})
}

// UpdateActivity maneja PUT /activities/:id
func (c *ActivitiesController) UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	var toUpdate domain.Activity
	if err := ctx.ShouldBindJSON(&toUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	updatedActivity, err := c.service.Update(ctx.Request.Context(), id, toUpdate)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}

		if errors.Is(err, services.ErrOwnerNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Owner not found or unauthorized"})
			return
		}

		if errors.Is(err, services.ErrOwnerForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to modify this resource"})
			return
		}

		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "required") {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{
			"error":   "Failed to update activity",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activity": updatedActivity,
		"message":  "Activity updated successfully",
	})
}

// DeleteActivity maneja DELETE /activities/:id
func (c *ActivitiesController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	// Hard delete - eliminación permanente
	err := c.service.HardDelete(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}

		if errors.Is(err, services.ErrOwnerNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Owner not found or unauthorized"})
			return
		}

		if errors.Is(err, services.ErrOwnerForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to modify this resource"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete activity",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Activity permanently deleted",
	})
}

// ToggleActiveActivity maneja PATCH /activities/:id/toggle
func (c *ActivitiesController) ToggleActiveActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	toggled, err := c.service.ToggleActive(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}

		if errors.Is(err, services.ErrOwnerNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Owner not found or unauthorized"})
			return
		}

		if errors.Is(err, services.ErrOwnerForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to modify this resource"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to toggle activity status",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activity": toggled,
		"message":  "Activity state toggled successfully",
	})
}
