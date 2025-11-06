package polymarketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetOpenInterest retrieves the open interest for markets
func (c *DataClient) GetOpenInterest(ctx context.Context, params *GetOpenInterestParams) ([]OpenInterest, error) {
	// Build query parameters
	queryParams := url.Values{}

	// Add optional parameters
	if len(params.Market) > 0 {
		queryParams.Set("market", strings.Join(params.Market, ","))
	}

	// Construct URL
	reqURL := fmt.Sprintf("%s/oi?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make open interest request: %w", err)
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Error)
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse successful response
	var openInterest []OpenInterest
	if err := json.Unmarshal(body, &openInterest); err != nil {
		return nil, fmt.Errorf("failed to decode open interest response: %w", err)
	}

	return openInterest, nil
}

// GetLiveVolume retrieves the live volume for an event
func (c *DataClient) GetLiveVolume(ctx context.Context, params *GetLiveVolumeParams) ([]LiveVolume, error) {
	if params.Id < 1 {
		return nil, fmt.Errorf("id must be >= 1")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("id", fmt.Sprintf("%d", params.Id))

	// Construct URL
	reqURL := fmt.Sprintf("%s/live-volume?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make live-volume request: %w", err)
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Error)
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse successful response
	var liveVolume []LiveVolume
	if err := json.Unmarshal(body, &liveVolume); err != nil {
		return nil, fmt.Errorf("failed to decode live-volume response: %w", err)
	}

	return liveVolume, nil
}
