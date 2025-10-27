package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"search-api/models"
	"search-api/repository"
)

type SearchService struct {
	repo  *repository.SolrRepository
	cache *DualCache
}

func NewSearchService(r *repository.SolrRepository, c *DualCache) *SearchService {
	return &SearchService{repo: r, cache: c}
}

func normalizeKey(q string, page, size int, sort string, filters url.Values) string {
	return fmt.Sprintf("q=%s|p=%d|s=%d|o=%s|f=%s", q, page, size, sort, filters.Encode())
}

func (s *SearchService) Search(q string, page, size int, sort string, filters url.Values) (models.SearchResult[map[string]any], error) {
	key := "search:" + normalizeKey(strings.TrimSpace(q), page, size, sort, filters)
	if str, ok := s.cache.Get(key); ok {
		var cached models.SearchResult[map[string]any]
		_ = json.Unmarshal([]byte(str), &cached)
		return cached, nil
	}
	res, err := s.repo.Search(q, page, size, sort, filters)
	if err != nil {
		return res, err
	}
	b, _ := json.Marshal(res)
	s.cache.Set(key, string(b))
	return res, nil
}

// Invalidar cache rápida para búsquedas repetidas
func (s *SearchService) Invalidate() { s.cache.InvalidatePrefix("search:") }
