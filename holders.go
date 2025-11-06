package polymarketdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetHolders retrieves top holders for markets
func (c *DataClient) GetHolders(params *GetHoldersParams) ([]MarketHolders, error) {
	if len(params.Market) == 0 {
		return nil, fmt.Errorf("market is required")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("market", strings.Join(params.Market, ","))

	// Add optional parameters
	if params.Limit > 0 {
		if params.Limit > 500 {
			return nil, fmt.Errorf("limit must be between 0 and 500")
		}
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.MinBalance > 0 {
		if params.MinBalance > 999999 {
			return nil, fmt.Errorf("minBalance must be between 0 and 999999")
		}
		queryParams.Set("minBalance", strconv.Itoa(params.MinBalance))
	}

	// Construct URL
	reqURL := fmt.Sprintf("%s/holders?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make holders request: %w", err)
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
	var holders []MarketHolders
	if err := json.Unmarshal(body, &holders); err != nil {
		return nil, fmt.Errorf("failed to decode holders response: %w", err)
	}

	return holders, nil
}
