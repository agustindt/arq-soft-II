package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetUserByID consulta la Users API para obtener datos del usuario
func GetUserByID(baseURL string, userID uint) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%d", baseURL, userID))
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
