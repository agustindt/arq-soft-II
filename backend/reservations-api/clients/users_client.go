package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserValidationResponse struct {
	ID    uint   `json:"id"`
	Role  string `json:"role"`
	Valid bool   `json:"valid"`
}

type UserClient struct {
	BaseURL string
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{BaseURL: baseURL}
}

func (uc *UserClient) ValidateToken(token string) (*UserValidationResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/validate", uc.BaseURL) // TODO mover a .env

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 🔐 reenviamos el token al servicio de usuarios
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response from user service: %d", resp.StatusCode)
	}

	var userResp UserValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	return &userResp, nil
}
