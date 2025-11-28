package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetUserByIDWithContext consulta la Users API para obtener datos del usuario
func GetUserByIDWithContext(ctx context.Context, baseURL string, userID uint) (*User, error) {
	trimmed := strings.TrimRight(baseURL, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/users/%d", trimmed, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("error calling users API: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling users API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user not found or not authorized")
	}

	var response struct {
		User User `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding user response: %w", err)
	}

	return &response.User, nil
}
