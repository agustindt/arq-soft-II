package controllers

import (
	"arq-soft-II/backend/activities-api/domain"
	"context"
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

// GetActivities maneja GET /activities - Lista actividades activas
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

// GetAllActivities maneja GET /activities/all - Lista todas las actividades (admin only)
func (c *ActivitiesController) GetAllActivities(ctx *gin.Context) {
	activities, err := c.service.ListAll(ctx.Request.Context())
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

// UpdateActivity maneja PUT /activities/:id - Actualiza actividad existente
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

// DeleteActivity maneja DELETE /activities/:id - Elimina actividad por ID
func (c *ActivitiesController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	err := c.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete activity",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Activity deleted successfully",
	})
}

// ToggleActiveActivity maneja PATCH /activities/:id/toggle - Activa/desactiva actividad
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

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to toggle activity status",
			"details": err.Error(),
		})
		return
	}

	status := "deactivated"
	if toggled.IsActive {
		status = "activated"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activity": toggled,
		"message":  "Activity " + status + " successfully",
	})
}

// GetActivitiesByCategory maneja GET /activities/category/:category - Filtra por categoría
func (c *ActivitiesController) GetActivitiesByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	if category == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Category parameter is required",
		})
		return
	}

	activities, err := c.service.GetByCategory(ctx.Request.Context(), category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch activities by category",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"category":   category,
		"activities": activities,
		"count":      len(activities),
	})
}

// HealthCheck maneja GET /healthz - Health check endpoint
func (c *ActivitiesController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "activities-api",
	})
}
