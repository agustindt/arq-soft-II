package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"arq-soft-II/backend/search-api/services"
)

type SearchController struct {
	service *services.SearchService
}

func NewSearchController(s *services.SearchService) *SearchController {
	return &SearchController{service: s}
}

func (c *SearchController) HandleSearch(w http.ResponseWriter, r *http.Request) {
	// Parámetros de búsqueda - aceptar tanto 'q' como 'query' para compatibilidad
	query := r.URL.Query().Get("q")
	if query == "" {
		query = r.URL.Query().Get("query")
	}

	// Filtros opcionales
	filters := make(map[string]interface{})

	if category := r.URL.Query().Get("category"); category != "" {
		filters["category"] = category
	}

	if difficulty := r.URL.Query().Get("difficulty"); difficulty != "" {
		filters["difficulty"] = difficulty
	}

	// Paginación
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		var pageNum int
		if _, err := fmt.Sscanf(p, "%d", &pageNum); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	filters["page"] = page

	size := 10
	if s := r.URL.Query().Get("size"); s != "" {
		var sizeNum int
		if _, err := fmt.Sscanf(s, "%d", &sizeNum); err == nil && sizeNum > 0 && sizeNum <= 100 {
			size = sizeNum
		}
	}
	filters["size"] = size

	// Realizar búsqueda
	result, err := c.service.Search(query, filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
