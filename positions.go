package polymarketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetPositions retrieves current positions for a user
func (c *Client) GetPositions(ctx context.Context, params *GetPositionsParams) ([]Position, error) {
	if params.User == "" {
		return nil, fmt.Errorf("user address is required")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("user", params.User)

	// Add optional parameters
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

	if params.SizeThreshold != nil {
		queryParams.Set("sizeThreshold", params.SizeThreshold.String())
	}

	if params.Redeemable != nil {
		queryParams.Set("redeemable", strconv.FormatBool(*params.Redeemable))
	}

	if params.Mergeable != nil {
		queryParams.Set("mergeable", strconv.FormatBool(*params.Mergeable))
	}

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

	if params.SortBy != "" {
		queryParams.Set("sortBy", string(params.SortBy))
	}

	if params.SortDirection != "" {
		queryParams.Set("sortDirection", string(params.SortDirection))
	}

	if params.Title != "" {
		if len(params.Title) > 100 {
			return nil, fmt.Errorf("title must not exceed 100 characters")
		}
		queryParams.Set("title", params.Title)
	}

	// Construct URL
	reqURL := fmt.Sprintf("%s/positions?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make positions request: %w", err)
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
	var positions []Position
	if err := json.Unmarshal(body, &positions); err != nil {
		return nil, fmt.Errorf("failed to decode positions response: %w", err)
	}

	return positions, nil
}

// GetClosedPositions fetches closed positions for a user
func (c *Client) GetClosedPositions(ctx context.Context, params *GetClosedPositionsParams) ([]ClosedPosition, error) {
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
	if len(params.Market) > 0 {
		queryParams.Set("market", strings.Join(params.Market, ","))
	}

	if params.Title != "" {
		if len(params.Title) > 100 {
			return nil, fmt.Errorf("title must not exceed 100 characters")
		}
		queryParams.Set("title", params.Title)
	}

	if len(params.EventId) > 0 {
		eventIds := make([]string, len(params.EventId))
		for i, id := range params.EventId {
			eventIds[i] = strconv.Itoa(id)
		}
		queryParams.Set("eventId", strings.Join(eventIds, ","))
	}

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

	if params.SortBy != "" {
		queryParams.Set("sortBy", string(params.SortBy))
	}

	if params.SortDirection != "" {
		queryParams.Set("sortDirection", string(params.SortDirection))
	}

	// Construct URL
	reqURL := fmt.Sprintf("%s/closed-positions?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make closed-positions request: %w", err)
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
	var closedPositions []ClosedPosition
	if err := json.Unmarshal(body, &closedPositions); err != nil {
		return nil, fmt.Errorf("failed to decode closed-positions response: %w", err)
	}

	return closedPositions, nil
}

// GetPositionsValue retrieves the total value of a user's positions
func (c *Client) GetPositionsValue(ctx context.Context, params *GetValueParams) ([]UserValue, error) {
	if params.User == "" {
		return nil, fmt.Errorf("user address is required")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("user", params.User)

	// Add optional parameters
	if len(params.Market) > 0 {
		queryParams.Set("market", strings.Join(params.Market, ","))
	}

	// Construct URL
	reqURL := fmt.Sprintf("%s/value?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make value request: %w", err)
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
	var userValues []UserValue
	if err := json.Unmarshal(body, &userValues); err != nil {
		return nil, fmt.Errorf("failed to decode value response: %w", err)
	}

	return userValues, nil
}
