package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MotivationAPI struct {
	url    string
	client *http.Client
}

func NewMotivationAPI(url string) *MotivationAPI {
	return &MotivationAPI{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (m *MotivationAPI) GetQuote(ctx context.Context) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", m.url, nil)

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("MotivationAPI - GetQuote - client.Do: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result []struct {
		Q string `json:"q"`
		A string `json:"a"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("MotivationAPI - GetQuote - decode: %w", err)
	}
	if len(result) > 0 {
		return result[0].Q + " - " + result[0].A, nil
	}
	return "keep pushing forward", nil
}
