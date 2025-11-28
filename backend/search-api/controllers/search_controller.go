package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchController struct {
	service SearchService
}

func (c *SearchController) HandleSearch(w http.ResponseWriter, r *http.Request) {
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

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		var limitNum int
		if _, err := fmt.Sscanf(l, "%d", &limitNum); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	} else if legacy := r.URL.Query().Get("size"); legacy != "" { // compatibilidad
		var legacyNum int
		if _, err := fmt.Sscanf(legacy, "%d", &legacyNum); err == nil && legacyNum > 0 && legacyNum <= 100 {
			limit = legacyNum
		}
	}
	filters["limit"] = limit

	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters["sort"] = sort
	}

	// Realizar búsqueda
	result, err := c.service.Search(query, filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
