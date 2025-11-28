package controllers

import (
	"arq-soft-II/backend/reservations-api/domain"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReservasService interface {
	// List retorna todos los Reservas (sin filtros por ahora)
	List(ctx context.Context) ([]domain.Reserva, error)

	// ListByUserID retorna las reservas de un usuario específico
	ListByUserID(ctx context.Context, userID int) ([]domain.Reserva, error)

	// Create valida y crea un nuevo Reserva
	Create(ctx context.Context, reserva domain.Reserva) (domain.Reserva, error)

	// GetByID obtiene un Reserva por su ID
	GetByID(ctx context.Context, id string) (domain.Reserva, error)

	// Update actualiza un Reserva existente
	Update(ctx context.Context, id string, reserva domain.Reserva) (domain.Reserva, error)

	// Delete elimina un Reserva por ID
	Delete(ctx context.Context, id string) error

	// GetScheduleAvailability retorna la disponibilidad de cada horario para una actividad
	GetScheduleAvailability(ctx context.Context, activityID string, date time.Time) (map[string]int, error)
}

// ReservasController maneja las peticiones HTTP para Reservas
type ReservasController struct {
	service ReservasService
}

// NewReservasController crea una nueva instancia del controller
func NewReservasController(reservasService ReservasService) *ReservasController {
	return &ReservasController{
		service: reservasService,
	}
}

func (c *ReservasController) GetReservas(ctx *gin.Context) {
	// Obtener información del usuario del contexto (seteada por AuthRequired middleware)
	userID, existsUserID := ctx.Get("user_id")
	userRole, existsUserRole := ctx.Get("user_role")

	if !existsUserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Obtener el parámetro de consulta 'scope' (opcional)
	// scope=mine -> siempre retorna solo las reservas del usuario actual (incluso si es admin)
	// scope=all -> retorna todas las reservas (solo para admins)
	// sin scope -> comportamiento por defecto (admins ven todo, usuarios ven solo las suyas)
	scope := ctx.Query("scope")

	var reservas []domain.Reserva
	var err error

	// Verificar si es admin
	isAdmin := false
	if existsUserRole {
		if role, ok := userRole.(string); ok {
			isAdmin = role == "admin" || role == "super_admin" || role == "root"
		}
	}

	// Si scope=mine, forzar a retornar solo las reservas del usuario actual
	if scope == "mine" {
		if uid, ok := userID.(uint); ok {
			reservas, err = c.service.ListByUserID(ctx.Request.Context(), int(uid))
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}
	} else {
		// Comportamiento por defecto
		// Si es admin, retornar todas las reservas
		if isAdmin {
			reservas, err = c.service.List(ctx.Request.Context())
		} else {
			// Si es usuario normal, retornar solo sus reservas
			if uid, ok := userID.(uint); ok {
				reservas, err = c.service.ListByUserID(ctx.Request.Context(), int(uid))
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Invalid user ID format",
				})
				return
			}
		}
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch Reservas",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reservas": reservas,
		"count":    len(reservas),
	})
}

// CreateReserva maneja POST /Reservas - Crea un nuevo Reserva
func (c *ReservasController) CreateReserva(ctx *gin.Context) {
	var newReserva domain.Reserva
	if err := ctx.ShouldBindJSON(&newReserva); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Obtener user_id del contexto (seteado por el middleware AuthRequired)
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			// Agregar el user_id a la reserva
			newReserva.UsersID = []int{int(uid)}
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Establecer timestamps si no están presentes
	if newReserva.CreatedAt.IsZero() {
		newReserva.CreatedAt = time.Now()
	}
	if newReserva.UpdatedAt.IsZero() {
		newReserva.UpdatedAt = time.Now()
	}
	if newReserva.Status == "" {
		newReserva.Status = "Pendiente"
	}

	Reserva, err := c.service.Create(ctx.Request.Context(), newReserva)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create Reserva",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"Reserva": Reserva,
	})
}

// GetReservaByID maneja GET /Reservas/:id - Obtiene Reserva por ID
func (c *ReservasController) GetReservaByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	Reserva, err := c.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "Reserva not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Reserva not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch Reserva",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Reserva": Reserva,
	})
}

// UpdateReserva maneja PUT /Reservas/:id - Actualiza Reserva existente
func (c *ReservasController) UpdateReserva(ctx *gin.Context) {
	var toUpdate domain.Reserva
	err := ctx.ShouldBindJSON(&toUpdate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	updatedReserva, err := c.service.Update(ctx, id, toUpdate)
	if err != nil {
		if err.Error() == "Reserva not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Reserva not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update Reserva",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Reserva": updatedReserva,
	})
}

// DeleteReserva maneja DELETE /Reservas/:id - Elimina Reserva por ID
// Usuarios pueden eliminar sus propias reservas, admins pueden eliminar cualquiera
func (c *ReservasController) DeleteReserva(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	// Obtener información del usuario del contexto
	userID, existsUserID := ctx.Get("user_id")
	userRole, existsUserRole := ctx.Get("user_role")

	if !existsUserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Obtener la reserva para verificar permisos
	reserva, err := c.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "Reserva not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Reserva not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch Reserva",
			"details": err.Error(),
		})
		return
	}

	// Verificar permisos: el usuario debe ser dueño de la reserva o ser admin
	uid, _ := userID.(uint)
	isOwner := false
	for _, userIDInReserva := range reserva.UsersID {
		if int(uid) == userIDInReserva {
			isOwner = true
			break
		}
	}

	// Si no es el dueño, verificar si es admin
	isAdmin := false
	if existsUserRole {
		if role, ok := userRole.(string); ok {
			isAdmin = role == "admin" || role == "super_admin" || role == "root"
		}
	}

	if !isOwner && !isAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own reservations",
		})
		return
	}

	// Eliminar la reserva
	err = c.service.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete Reserva",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Reserva eliminada exitosamente",
	})
}

// GetScheduleAvailability maneja GET /activities/:id/availability - Obtiene disponibilidad por horario
func (c *ReservasController) GetScheduleAvailability(ctx *gin.Context) {
	activityID := ctx.Param("id")
	if activityID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Activity ID parameter is required",
		})
		return
	}

	// Obtener fecha del query parameter (formato YYYY-MM-DD)
	dateStr := ctx.Query("date")
	if dateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Date query parameter is required (format: YYYY-MM-DD)",
		})
		return
	}

	// Parsear la fecha
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid date format",
			"details": "Date must be in YYYY-MM-DD format",
		})
		return
	}

	// Obtener disponibilidad
	availability, err := c.service.GetScheduleAvailability(ctx.Request.Context(), activityID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch schedule availability",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activity_id":  activityID,
		"date":         dateStr,
		"availability": availability,
	})
}

// 200 OK - Operación exitosa con contenido
// 201 Created - Recurso creado exitosamente
// 204 No Content - Operación exitosa sin contenido (típico para DELETE)
// 400 Bad Request - Error en los datos enviados por el cliente
// 404 Not Found - Recurso no encontrado
// 500 Internal Server Error - Error interno del servidor
// 501 Not Implemented - Funcionalidad no implementada (para TODOs)
