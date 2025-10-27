package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"search-api/models"
)

type SolrRepository struct {
	base string
	core string
	http *http.Client
}

func NewSolrRepository(base, core string) *SolrRepository {
	return &SolrRepository{base: base, core: core, http: &http.Client{Timeout: 10 * time.Second}}
}

// Indexa/actualiza un documento
func (r *SolrRepository) Upsert(doc any) error {
	payload := map[string]any{"add": map[string]any{"doc": doc, "commitWithin": 1000}}
	b, _ := json.Marshal(payload)
	endpoint := fmt.Sprintf("%s/%s/update?commit=true", r.base, r.core)
	resp, err := r.http.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		rb, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr upsert %d: %s", resp.StatusCode, string(rb))
	}
	return nil
}

func (r *SolrRepository) DeleteByID(id string) error {
	payload := map[string]any{"delete": map[string]string{"id": id}}
	b, _ := json.Marshal(payload)
	endpoint := fmt.Sprintf("%s/%s/update?commit=true", r.base, r.core)
	resp, err := r.http.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		rb, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr delete %d: %s", resp.StatusCode, string(rb))
	}
	return nil
}

// Búsqueda con paginación, filtros y orden
func (r *SolrRepository) Search(q string, page, size int, sort string, filters url.Values) (models.SearchResult[map[string]any], error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	start := (page - 1) * size
	params := url.Values{}
	if q == "" {
		q = "*:*"
	}
	params.Set("q", q)
	params.Set("start", fmt.Sprint(start))
	params.Set("rows", fmt.Sprint(size))
	if sort != "" {
		params.Set("sort", sort)
	}
	for k, vs := range filters {
		for _, v := range vs {
			params.Add("fq", fmt.Sprintf("%s:%s", k, v))
		}
	}
	endpoint := fmt.Sprintf("%s/%s/select?%s", r.base, r.core, params.Encode())
	resp, err := r.http.Get(endpoint)
	if err != nil {
		return models.SearchResult[map[string]any]{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		rb, _ := io.ReadAll(resp.Body)
		return models.SearchResult[map[string]any]{}, fmt.Errorf("solr search %d: %s", resp.StatusCode, string(rb))
	}
	var raw struct {
		Response struct {
			NumFound int              `json:"numFound"`
			Docs     []map[string]any `json:"docs"`
		} `json:"response"`
	}
	b, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &raw); err != nil {
		return models.SearchResult[map[string]any]{}, err
	}
	return models.SearchResult[map[string]any]{Total: raw.Response.NumFound, Page: page, Size: size, Items: raw.Response.Docs}, nil
}
