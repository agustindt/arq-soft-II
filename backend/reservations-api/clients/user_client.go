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

// TODO hacer bien esta struct

func GetUserByID(baseURL string, userID uint) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, userID))
	if err != nil {
		return nil, fmt.Errorf("error calling users API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user not found or not authorized")
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error decoding user response: %w", err)
	}

	return &user, nil
}
