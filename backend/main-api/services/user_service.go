package services

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// UserService validates owners via external HTTP API
type UserService interface {
	Exists(ctx context.Context, userID string) (bool, error)
}

type httpUserService struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPUserService constructs a UserService
func NewHTTPUserService(baseURL string) UserService {
	return &httpUserService{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *httpUserService) Exists(ctx context.Context, userID string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", h.baseURL, userID), nil)
	if err != nil {
		return false, err
	}
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, nil
}
