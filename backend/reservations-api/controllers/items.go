package controllers

import (
	"context"
	"net/http"
	"reservations/domain"

	"github.com/gin-gonic/gin"
)

// ReservasService define la lógica de negocio para Reservas
// Capa intermedia entre Controllers (HTTP) y Repository (datos)
// Responsabilidades: validaciones, transformaciones, reglas de negocio
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
// Responsabilidades:
// - Extraer datos del request (JSON, path params, query params)
// - Validar formato de entrada
// - Llamar al service correspondiente
// - Retornar respuesta HTTP adecuada
type ReservasController struct {
	service ReservasService // Inyección de dependencia
}

// NewReservasController crea una nueva instancia del controller
func NewReservasController(reservasService ReservasService) *ReservasController {
	return &ReservasController{
		service: reservasService,
	}
}

// ✅ IMPLEMENTADO - Ejemplo para que los estudiantes entiendan el patrón
func (c *ReservasController) GetReservas(ctx *gin.Context) {
	// 🔍 Llamar al service para obtener los datos
	reservas, err := c.service.List(ctx.Request.Context())
	if err != nil {
		// ❌ Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch Reservas",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con los datos
	ctx.JSON(http.StatusOK, gin.H{
		"reservas": reservas,
		"count":    len(reservas),
	})
}

// CreateReserva maneja POST /Reservas - Crea un nuevo Reserva
// Consigna 1: Recibir JSON, validar y crear Reserva
func (c *ReservasController) CreateReserva(ctx *gin.Context) {
	// Obtener el Reserva del body JSON
	var newReserva domain.Reserva
	if err := ctx.ShouldBindJSON(&newReserva); err != nil {
		// ❌ Error en los datos enviados por el cliente
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	Reserva, err := c.service.Create(ctx, newReserva)
	if err != nil {
		// ❌ Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create Reserva",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con el Reserva creado
	ctx.JSON(http.StatusCreated, gin.H{
		"Reserva": Reserva,
	})
}

// GetReservaByID maneja GET /Reservas/:id - Obtiene Reserva por ID
// Consigna 2: Extraer ID del path param, validar y buscar
func (c *ReservasController) GetReservaByID(ctx *gin.Context) {
	// Obtener el ID del path param
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
// Consigna 3: Extraer ID y datos, validar y actualizar
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
// Consigna 4: Extraer ID, validar y eliminar
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

	ctx.JSON(http.StatusNoContent, nil) // 204 No Content
}

// 📚 Notas sobre HTTP Status Codes
//
// 200 OK - Operación exitosa con contenido
// 201 Created - Recurso creado exitosamente
// 204 No Content - Operación exitosa sin contenido (típico para DELETE)
// 400 Bad Request - Error en los datos enviados por el cliente
// 404 Not Found - Recurso no encontrado
// 500 Internal Server Error - Error interno del servidor
// 501 Not Implemented - Funcionalidad no implementada (para TODOs)
//
// 💡 Tip: En una API real, sería buena práctica crear una función
// helper para manejar respuestas de error de manera consistente
