package polymarketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HealthCheck performs a health check on the Polymarket Data API
// Returns "OK" if the API is healthy
func (c *Client) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	url := Endpoint + "/"

	resp, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to make health check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("health check failed with status %d: %s", resp.StatusCode, string(body))
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return nil, fmt.Errorf("failed to decode health check response: %w", err)
	}

	return &healthResp, nil
}
