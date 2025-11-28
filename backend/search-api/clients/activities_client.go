package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"search-api/domain"
)

type activityResponse struct {
	Activity domain.Activity `json:"activity"`
}

// GetActivityByID obtiene una actividad real desde activities-api vía HTTP.
func GetActivityByID(ctx context.Context, baseURL, id string) (domain.Activity, error) {
	trimmed := strings.TrimRight(baseURL, "/")
	url := fmt.Sprintf("%s/activities/%s", trimmed, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("crear request de actividad: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("llamar activities-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Activity{}, fmt.Errorf("activities-api status %d", resp.StatusCode)
	}

	var payload activityResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return domain.Activity{}, fmt.Errorf("decodificar actividad: %w", err)
	}

	return payload.Activity, nil
}
