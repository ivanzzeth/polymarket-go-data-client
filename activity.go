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

// GetActivity retrieves on-chain activity for a user
func (c *DataClient) GetActivity(params *GetActivityParams) ([]Activity, error) {
	if params.User == "" {
		return nil, fmt.Errorf("user address is required")
	}

	// Validate mutually exclusive parameters
	if len(params.Market) > 0 && len(params.EventId) > 0 {
		return nil, fmt.Errorf("market and eventId are mutually exclusive")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("user", params.User)

	// Add optional parameters
	if params.Limit > 0 {
		if params.Limit > 500 {
			return nil, fmt.Errorf("limit must be between 0 and 500")
		}
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Offset > 0 {
		if params.Offset > 10000 {
			return nil, fmt.Errorf("offset must be between 0 and 10000")
		}
		queryParams.Set("offset", strconv.Itoa(params.Offset))
	}

	if len(params.Market) > 0 {
		queryParams.Set("market", strings.Join(params.Market, ","))
	}

	if len(params.EventId) > 0 {
		eventIds := make([]string, len(params.EventId))
		for i, id := range params.EventId {
			eventIds[i] = strconv.Itoa(id)
		}
		queryParams.Set("eventId", strings.Join(eventIds, ","))
	}

	if len(params.Type) > 0 {
		types := make([]string, len(params.Type))
		for i, t := range params.Type {
			types[i] = string(t)
		}
		queryParams.Set("type", strings.Join(types, ","))
	}

	if params.Start > 0 {
		queryParams.Set("start", strconv.FormatInt(params.Start, 10))
	}

	if params.End > 0 {
		queryParams.Set("end", strconv.FormatInt(params.End, 10))
	}

	if params.SortBy != "" {
		queryParams.Set("sortBy", string(params.SortBy))
	}

	if params.SortDirection != "" {
		queryParams.Set("sortDirection", string(params.SortDirection))
	}

	if params.Side != "" {
		queryParams.Set("side", string(params.Side))
	}

	// Construct URL
	reqURL := fmt.Sprintf("%s/activity?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make activity request: %w", err)
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
	var activities []Activity
	if err := json.Unmarshal(body, &activities); err != nil {
		return nil, fmt.Errorf("failed to decode activity response: %w", err)
	}

	return activities, nil
}
