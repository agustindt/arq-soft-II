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

	// Create valida y crea un nuevo Reserva
	Create(ctx context.Context, reserva domain.Reserva) (domain.Reserva, error)

	// GetByID obtiene un Reserva por su ID
	GetByID(ctx context.Context, id string) (domain.Reserva, error)

	// Update actualiza un Reserva existente
	Update(ctx context.Context, id string, reserva domain.Reserva) (domain.Reserva, error)

	// Delete elimina un Reserva por ID
	Delete(ctx context.Context, id string) error
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
	reservas, err := c.service.List(ctx.Request.Context())
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
func (c *ReservasController) DeleteReserva(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	err := c.service.Delete(ctx, id)
	if err != nil {
		if err.Error() == "Reserva not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Reserva not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete Reserva",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// 200 OK - Operación exitosa con contenido
// 201 Created - Recurso creado exitosamente
// 204 No Content - Operación exitosa sin contenido (típico para DELETE)
// 400 Bad Request - Error en los datos enviados por el cliente
// 404 Not Found - Recurso no encontrado
// 500 Internal Server Error - Error interno del servidor
// 501 Not Implemented - Funcionalidad no implementada (para TODOs)
