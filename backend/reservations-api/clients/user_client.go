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

func GetUserByID(baseURL string, userID uint) (*User, error) {
	// The users-api endpoint is /api/v1/users/:id
	url := fmt.Sprintf("%s/api/v1/users/%d", baseURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling users API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user not found or not authorized (status: %d)", resp.StatusCode)
	}

	// The response format is: {"message": "...", "data": {...}}
	var response struct {
		Message string `json:"message"`
		Data    User   `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding user response: %w", err)
	}

	return &response.Data, nil
}
