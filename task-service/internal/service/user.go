package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rfashwall/task-service/internal/models"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

type HTTPUserService struct {
	BaseURL string
}

func NewHTTPUserService(baseURL string) *HTTPUserService {
	return &HTTPUserService{BaseURL: baseURL}
}

func (s *HTTPUserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	url := fmt.Sprintf("%s/users/%d", s.BaseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: status %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
