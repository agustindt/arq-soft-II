package services

import (
	"encoding/json"
	"fmt"
	"log"
	"search-api/clients"
	"search-api/domain"
	"search-api/utils"
	"strings"
	"time"
)

type SearchService struct {
	cache      *utils.Cache
	solrClient *clients.SolrClient
	ttl        time.Duration
}

type SearchResult struct {
	Query      string                 `json:"query"`
	Results    []clients.SolrDocument `json:"results"`
	TotalFound int                    `json:"total_found"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	Timestamp  string                 `json:"timestamp"`
}

func NewSearchService(c *utils.Cache, solr *clients.SolrClient, ttl time.Duration) *SearchService {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}

	return &SearchService{
		cache:      c,
		solrClient: solr,
		ttl:        ttl,
	}
}

// Search realiza una bﾃｺsqueda en Solr con cache
func (s *SearchService) Search(query string, filters map[string]interface{}) (*SearchResult, error) {
	// Normalizar paginación y ordenamiento
	page := 1
	if p, ok := filters["page"].(int); ok && p > 0 {
		page = p
	}
	limit := 10
	if l, ok := filters["limit"].(int); ok && l > 0 && l <= 100 {
		limit = l
	}
	sort := ""
	if srt, ok := filters["sort"].(string); ok {
		sort = srt
	}
	filters["page"] = page
	filters["limit"] = limit
	if sort != "" {
		filters["sort"] = sort
	}

	// Generar key de cache basada en query y filtros
	cacheKey := s.generateCacheKey(query, filters)

	// Intentar obtener desde cachﾃｩ
	if data, err := s.cache.Get(cacheKey); err == nil {
		var result SearchResult
		if err := json.Unmarshal(data, &result); err == nil {
			log.Printf("笞｡ Cache hit: %s", query)
			return &result, nil
		}
	}

	// Si no estﾃ｡ en cachﾃｩ, buscar en Solr
	log.Printf("剥 Buscando en Solr: %s", query)

	solrResp, err := s.solrClient.Search(query, filters)
	if err != nil {
		return nil, fmt.Errorf("error searching in Solr: %w", err)
	}

	// Preparar resultado
	result := SearchResult{
		Query:      query,
		Results:    solrResp.Response.Docs,
		TotalFound: solrResp.Response.NumFound,
		Page:       page,
		Limit:      limit,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	// Guardar en cachﾃｩ
	data, _ := json.Marshal(result)
	if err := s.cache.Set(cacheKey, data, s.ttl); err != nil {
		log.Printf("笞・・ No se pudo guardar en cachﾃｩ: %v", err)
	}

	log.Printf("笨・Encontrados %d resultados para: %s", result.TotalFound, query)

	return &result, nil
}

// IndexActivity indexa una actividad en Solr
func (s *SearchService) IndexActivity(activity domain.Activity) error {
	log.Printf("統 Indexando actividad en Solr: %s (ID: %s)", activity.Name, activity.ID)

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

	log.Printf("笨・Actividad indexada correctamente: %s", activity.ID)

	// Invalidar cachﾃｩ relacionada
	s.invalidateCache()

	return nil
}

// DeleteActivity elimina una actividad de Solr
func (s *SearchService) DeleteActivity(activityID string) error {
	log.Printf("卵・・ Eliminando actividad de Solr: %s", activityID)

	err := s.solrClient.Delete(activityID)
	if err != nil {
		return fmt.Errorf("error deleting activity from Solr: %w", err)
	}

	log.Printf("笨・Actividad eliminada de Solr: %s", activityID)

	// Invalidar cachﾃｩ
	s.invalidateCache()

	return nil
}

// UpdateActivity actualiza una actividad en Solr (bﾃ｡sicamente es un re-index)
func (s *SearchService) UpdateActivity(activity domain.Activity) error {
	log.Printf("売 Actualizando actividad en Solr: %s (ID: %s)", activity.Name, activity.ID)

	// Para Solr, actualizar es lo mismo que indexar (sobreescribe)
	return s.IndexActivity(activity)
}

// generateCacheKey genera una key ﾃｺnica para el cache basada en query y filtros
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

	if limit, ok := filters["limit"].(int); ok {
		parts = append(parts, fmt.Sprintf("limit:%d", limit))
	}

	if sort, ok := filters["sort"].(string); ok && sort != "" {
		parts = append(parts, "sort:"+sort)
	}

	return strings.Join(parts, ":")
}

// invalidateCache limpia el cache (simplificado - en producciﾃｳn serﾃｭa mﾃ｡s selectivo)
func (s *SearchService) invalidateCache() {
	// En una implementaciﾃｳn real, aquﾃｭ invalidarﾃｭamos solo las keys relacionadas
	// Por ahora, dejamos que expiren naturalmente (30s TTL)
	log.Println("笞・・ Cache invalidation: entries will expire naturally")
}
