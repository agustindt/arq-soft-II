package services

import (
	"arq-soft-II/backend/search-api/clients"
	"arq-soft-II/backend/search-api/domain"
	"arq-soft-II/config/cache"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type SearchService struct {
	cache      *cache.Cache
	solrClient *clients.SolrClient
}

type SearchResult struct {
	Query      string                 `json:"query"`
	Results    []clients.SolrDocument `json:"results"`
	TotalFound int                    `json:"total_found"`
	Timestamp  string                 `json:"timestamp"`
}

func NewSearchService(c *cache.Cache, solr *clients.SolrClient) *SearchService {
	return &SearchService{
		cache:      c,
		solrClient: solr,
	}
}

// Search realiza una búsqueda en Solr con cache
func (s *SearchService) Search(query string, filters map[string]interface{}) (*SearchResult, error) {
	// Generar key de cache basada en query y filtros
	cacheKey := s.generateCacheKey(query, filters)

	// Intentar obtener desde caché
	if data, err := s.cache.Get(cacheKey); err == nil {
		var result SearchResult
		if err := json.Unmarshal(data, &result); err == nil {
			log.Printf("⚡ Cache hit: %s", query)
			return &result, nil
		}
	}

	// Si no está en caché, buscar en Solr
	log.Printf("🔍 Buscando en Solr: %s", query)

	solrResp, err := s.solrClient.Search(query, filters)
	if err != nil {
		return nil, fmt.Errorf("error searching in Solr: %w", err)
	}

	// Preparar resultado
	result := SearchResult{
		Query:      query,
		Results:    solrResp.Response.Docs,
		TotalFound: solrResp.Response.NumFound,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	// Guardar en caché
	data, _ := json.Marshal(result)
	if err := s.cache.Set(cacheKey, data, 30*time.Second); err != nil {
		log.Printf("⚠️  No se pudo guardar en caché: %v", err)
	}

	log.Printf("✅ Encontrados %d resultados para: %s", result.TotalFound, query)

	return &result, nil
}

// IndexActivity indexa una actividad en Solr
func (s *SearchService) IndexActivity(activity domain.Activity) error {
	log.Printf("📝 Indexando actividad en Solr: %s (ID: %s)", activity.Name, activity.ID)

	// Convertir Activity a SolrDocument
	doc := clients.SolrDocument{
		ID:          activity.ID,
		Name:        activity.Name,
		Description: activity.Description,
		Category:    activity.Category,
		Difficulty:  activity.Difficulty,
		Location:    activity.Location,
		Price:       activity.Price,
		DateCreated: activity.CreatedAt.Format(time.RFC3339),
		Text: []string{
			activity.Name,
			activity.Description,
			activity.Location,
		},
	}

	err := s.solrClient.Index(doc)
	if err != nil {
		return fmt.Errorf("error indexing activity: %w", err)
	}

	log.Printf("✅ Actividad indexada correctamente: %s", activity.ID)

	// Invalidar caché relacionada
	s.invalidateCache()

	return nil
}

// DeleteActivity elimina una actividad de Solr
func (s *SearchService) DeleteActivity(activityID string) error {
	log.Printf("🗑️  Eliminando actividad de Solr: %s", activityID)

	err := s.solrClient.Delete(activityID)
	if err != nil {
		return fmt.Errorf("error deleting activity from Solr: %w", err)
	}

	log.Printf("✅ Actividad eliminada de Solr: %s", activityID)

	// Invalidar caché
	s.invalidateCache()

	return nil
}

// UpdateActivity actualiza una actividad en Solr (básicamente es un re-index)
func (s *SearchService) UpdateActivity(activity domain.Activity) error {
	log.Printf("🔄 Actualizando actividad en Solr: %s (ID: %s)", activity.Name, activity.ID)

	// Para Solr, actualizar es lo mismo que indexar (sobreescribe)
	return s.IndexActivity(activity)
}

// generateCacheKey genera una key única para el cache basada en query y filtros
func (s *SearchService) generateCacheKey(query string, filters map[string]interface{}) string {
	parts := []string{"search", query}

	if category, ok := filters["category"].(string); ok && category != "" {
		parts = append(parts, "cat:"+category)
	}

	if difficulty, ok := filters["difficulty"].(string); ok && difficulty != "" {
		parts = append(parts, "diff:"+difficulty)
	}

	if page, ok := filters["page"].(int); ok {
		parts = append(parts, fmt.Sprintf("page:%d", page))
	}

	if size, ok := filters["size"].(int); ok {
		parts = append(parts, fmt.Sprintf("size:%d", size))
	}

	return strings.Join(parts, ":")
}

// invalidateCache limpia el cache (simplificado - en producción sería más selectivo)
func (s *SearchService) invalidateCache() {
	// En una implementación real, aquí invalidaríamos solo las keys relacionadas
	// Por ahora, dejamos que expiren naturalmente (30s TTL)
	log.Println("⚠️  Cache invalidation: entries will expire naturally")
}
