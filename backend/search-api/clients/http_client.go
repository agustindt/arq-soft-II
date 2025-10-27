package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"search-api/models"
)

type ActivitiesClient struct {
	BaseURL string
	http    *http.Client
}

func NewActivitiesClient(base string) *ActivitiesClient {
	return &ActivitiesClient{BaseURL: base, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *ActivitiesClient) GetByID(id string) (*models.Activity, error) {
	url := fmt.Sprintf("%s/activities/%s", c.BaseURL, id)
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("activities-api status %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	var out models.Activity
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
