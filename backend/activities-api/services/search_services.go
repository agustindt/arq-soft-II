package services

import (
	"encoding/json"
	"log"
	"time"

	"arq-soft-II/config/cache"
)

type SearchService struct {
	cache *cache.Cache
}

type SearchResult struct {
	Query     string   `json:"query"`
	Results   []string `json:"results"`
	Timestamp string   `json:"timestamp"`
}

func NewSearchService(c *cache.Cache) *SearchService {
	return &SearchService{cache: c}
}

func (s *SearchService) Search(query string) (*SearchResult, error) {
	cacheKey := "search:" + query

	// Intentar obtener desde caché
	if data, err := s.cache.Get(cacheKey); err == nil {
		var result SearchResult
		if err := json.Unmarshal(data, &result); err == nil {
			log.Println("⚡ Resultado recuperado desde caché:", query)
			return &result, nil
		}
	}

	// Si no está, simulamos búsqueda (Solr aún no integrado)
	log.Println("🔍 Realizando búsqueda simulada:", query)
	result := SearchResult{
		Query:     query,
		Results:   []string{"Actividad 1", "Actividad 2", "Actividad 3"},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Guardamos en caché
	data, _ := json.Marshal(result)
	if err := s.cache.Set(cacheKey, data, 30*time.Second); err != nil {
		log.Println("⚠️ No se pudo guardar en caché:", err)
	}

	return &result, nil
}
