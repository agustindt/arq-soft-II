package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Activity representa la estructura de una actividad desde la Activities API
type Activity struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MaxCapacity int      `json:"max_capacity"`
	Schedule    []string `json:"schedule"`
	Duration    int      `json:"duration"` // Duración en minutos
}

// ActivityData representa los datos de la actividad
type ActivityData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MaxCapacity int      `json:"max_capacity"`
	Schedule    []string `json:"schedule"`
	Duration    int      `json:"duration"` // Duración en minutos
}

// ActivityResponse representa la respuesta completa de la Activities API (con wrapper)
type ActivityResponse struct {
	Activity ActivityData `json:"activity"`
}

// GetActivityByID obtiene los detalles de una actividad desde la Activities API
func GetActivityByID(baseURL string, activityID string) (*Activity, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}
	if activityID == "" {
		return nil, fmt.Errorf("activityID cannot be empty")
	}

	// Construir URL del endpoint
	url := fmt.Sprintf("%s/activities/%s", baseURL, activityID)

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Realizar request GET
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling activities API: %w", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("activity not found (ID: %s)", activityID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("activities API returned status %d", resp.StatusCode)
	}

	// Decodificar respuesta JSON (con wrapper "activity")
	var activityResp ActivityResponse
	if err := json.NewDecoder(resp.Body).Decode(&activityResp); err != nil {
		return nil, fmt.Errorf("error decoding activity response: %w", err)
	}

	// Convertir a struct Activity
	activity := &Activity{
		ID:          activityResp.Activity.ID,
		Name:        activityResp.Activity.Name,
		MaxCapacity: activityResp.Activity.MaxCapacity,
		Schedule:    activityResp.Activity.Schedule,
		Duration:    activityResp.Activity.Duration,
	}

	return activity, nil
}
