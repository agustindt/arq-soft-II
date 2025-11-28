package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (s *SolrClient) Search(query string, filters map[string]interface{}) (*SolrResponse, error) {
	selectURL := fmt.Sprintf("%s/%s/select", s.BaseURL, s.Collection)

	params := url.Values{}
	params.Set("q", query)

	// Filtros
	if category, ok := filters["category"].(string); ok && category != "" {
		params.Set("fq", fmt.Sprintf("category:%s", category))
	}

	if difficulty, ok := filters["difficulty"].(string); ok && difficulty != "" {
		if fq := params.Get("fq"); fq != "" {
			params.Set("fq", fmt.Sprintf("%s AND difficulty:%s", fq, difficulty))
		} else {
			params.Set("fq", fmt.Sprintf("difficulty:%s", difficulty))
		}
	}

	// Rango de precio
	if priceMin, ok := filters["price_min"].(float64); ok {
		if priceMax, ok := filters["price_max"].(float64); ok {
			if fq := params.Get("fq"); fq != "" {
				params.Set("fq", fmt.Sprintf("%s AND price:[%f TO %f]", fq, priceMin, priceMax))
			} else {
				params.Set("fq", fmt.Sprintf("price:[%f TO %f]", priceMin, priceMax))
			}
		}
	}

	// Tamaño por defecto
	size := 10
	if s, ok := filters["limit"].(int); ok && s > 0 {
		size = s
	}

	// Paginación
	if page, ok := filters["page"].(int); ok && page > 0 {
		start := (page - 1) * size
		params.Set("start", fmt.Sprintf("%d", start))
		params.Set("rows", fmt.Sprintf("%d", size))
	} else {
		params.Set("rows", fmt.Sprintf("%d", size))
	}

	// Ordenamiento
	if sort, ok := filters["sort"].(string); ok && sort != "" {
		sortParts := strings.Split(sort, "_")
		if len(sortParts) == 2 {
			direction := "asc"
			if strings.ToLower(sortParts[1]) == "desc" {
				direction = "desc"
			}
			params.Set("sort", fmt.Sprintf("%s %s", sortParts[0], direction))
		}
	}

	// Formato de respuesta
	params.Set("wt", "json")

	// Construir URL completa
	fullURL := fmt.Sprintf("%s?%s", selectURL, params.Encode())

	resp, err := s.Client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("solr returned status %d: %s", resp.StatusCode, string(body))
	}

	var solrResp SolrResponse
	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		return nil, fmt.Errorf("error decoding Solr response: %w", err)
	}

	return &solrResp, nil
}
