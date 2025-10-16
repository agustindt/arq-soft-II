package controllers

import (
	"encoding/json"
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
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Falta parámetro 'query'", http.StatusBadRequest)
		return
	}

	result, err := c.service.Search(query)
	if err != nil {
		http.Error(w, "Error al buscar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
