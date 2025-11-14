package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SolrClient maneja la comunicación con Apache Solr
type SolrClient struct {
	BaseURL string
	Core    string
	Client  *http.Client
}

// SolrDocument representa un documento de actividad en Solr
type SolrDocument struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Difficulty  string   `json:"difficulty"`
	Location    string   `json:"location"`
	Price       float64  `json:"price"`
	DateCreated string   `json:"date_created"`
	Text        []string `json:"text,omitempty"` // campo multi-valor para búsqueda full-text
}

// SolrResponse representa la respuesta de una búsqueda en Solr
type SolrResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Response struct {
		NumFound int            `json:"numFound"`
		Start    int            `json:"start"`
		Docs     []SolrDocument `json:"docs"`
	} `json:"response"`
}

// NewSolrClient crea una nueva instancia del cliente Solr
func NewSolrClient(baseURL, core string) *SolrClient {
	return &SolrClient{
		BaseURL: baseURL,
		Core:    core,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Index indexa un documento en Solr
func (s *SolrClient) Index(doc SolrDocument) error {
	updateURL := fmt.Sprintf("%s/%s/update?commit=true", s.BaseURL, s.Core)

	// Solr espera un array de documentos
	payload := []SolrDocument{doc}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling document: %w", err)
	}

	req, err := http.NewRequest("POST", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Search realiza una búsqueda en Solr con filtros opcionales
func (s *SolrClient) Search(query string, filters map[string]interface{}) (*SolrResponse, error) {
	selectURL := fmt.Sprintf("%s/%s/select", s.BaseURL, s.Core)

	// Construir parámetros de búsqueda
	params := url.Values{}

	// Query principal
	if query == "" || query == "*" {
		params.Set("q", "*:*") // búsqueda de todos los documentos
	} else {
		// Búsqueda con wildcards: convertir a minúsculas y agregar wildcard al final
		// Los campos text_general en Solr ya convierten a lowercase
		wildcardQuery := strings.ToLower(query) + "*"
		params.Set("q", wildcardQuery)
		params.Set("defType", "edismax")
		params.Set("qf", "name^3 description^2 location category") // boost en name
	}

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

	// Paginación
	if page, ok := filters["page"].(int); ok && page > 0 {
		size := 10 // por defecto
		if s, ok := filters["size"].(int); ok {
			size = s
		}
		start := (page - 1) * size
		params.Set("start", fmt.Sprintf("%d", start))
		params.Set("rows", fmt.Sprintf("%d", size))
	} else {
		params.Set("rows", "10")
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

// Delete elimina un documento de Solr por ID
func (s *SolrClient) Delete(id string) error {
	updateURL := fmt.Sprintf("%s/%s/update?commit=true", s.BaseURL, s.Core)

	payload := map[string]interface{}{
		"delete": map[string]string{
			"id": id,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling delete command: %w", err)
	}

	req, err := http.NewRequest("POST", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// HealthCheck verifica que Solr esté accesible
func (s *SolrClient) HealthCheck() error {
	pingURL := fmt.Sprintf("%s/%s/admin/ping", s.BaseURL, s.Core)

	resp, err := s.Client.Get(pingURL)
	if err != nil {
		return fmt.Errorf("error pinging Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr health check failed with status %d", resp.StatusCode)
	}

	return nil
}
