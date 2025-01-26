package users

import (
	"bytes"
	"context"
	"courses-api/domain/users"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HTTP struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewHTTP(config HTTPConfig) *HTTP {
	return &HTTP{
		baseURL: fmt.Sprintf("http://%s:%s", config.Host, config.Port),
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		apiKey: config.APIKey,
	}
}

func (h *HTTP) GetUserByID(ctx context.Context, id string) (users.User, error) {
	url := fmt.Sprintf("%s/users/%s", h.baseURL, id)
	log.Printf("HTTP Repository - GetUserByID - Making request to: %s\n", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("HTTP Repository - GetUserByID - Error creating request: %v\n", err)
		return users.User{}, fmt.Errorf("error creating request: %w", err)
	}

	// Agregar headers necesarios
	req.Header.Set("Content-Type", "application/json")
	if h.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.apiKey))
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("HTTP Repository - GetUserByID - Error making request: %v\n", err)
		return users.User{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("HTTP Repository - GetUserByID - Response status code: %d\n", resp.StatusCode)

	// Para debug, imprime el body de la respuesta
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("HTTP Repository - GetUserByID - Response body: %s\n", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return users.User{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Crear un nuevo reader con los bytes leídos
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var user users.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("HTTP Repository - GetUserByID - Error decoding response: %v\n", err)
		return users.User{}, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("HTTP Repository - GetUserByID - User retrieved: %+v\n", user)
	return user, nil
}

func (h *HTTP) ValidateAdminUser(ctx context.Context, userID string) (bool, error) {
	log.Printf("HTTP Repository - ValidateAdminUser - Starting validation for userID: %s\n", userID)

	user, err := h.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("HTTP Repository - ValidateAdminUser - Error getting user: %v\n", err)
		return false, fmt.Errorf("error getting user: %w", err)
	}

	log.Printf("HTTP Repository - ValidateAdminUser - User found: %+v\n", user)
	log.Printf("HTTP Repository - ValidateAdminUser - IsAdmin: %v\n", user.IsAdmin)

	return user.IsAdmin, nil
}
