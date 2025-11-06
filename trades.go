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

// GetTrades retrieves trades for a user or markets
func (c *DataClient) GetTrades(ctx context.Context, params *GetTradesParams) ([]Trade, error) {
	// Validate mutually exclusive parameters
	if len(params.Market) > 0 && len(params.EventId) > 0 {
		return nil, fmt.Errorf("market and eventId are mutually exclusive")
	}

	// Validate FilterType and FilterAmount must be provided together
	if (params.FilterType != "" && params.FilterAmount == nil) || (params.FilterType == "" && params.FilterAmount != nil) {
		return nil, fmt.Errorf("filterType and filterAmount must be provided together")
	}

	// Build query parameters
	queryParams := url.Values{}

	// Add optional parameters
	if params.Limit > 0 {
		if params.Limit > 10000 {
			return nil, fmt.Errorf("limit must be between 0 and 10000")
		}
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Offset > 0 {
		if params.Offset > 10000 {
			return nil, fmt.Errorf("offset must be between 0 and 10000")
		}
		queryParams.Set("offset", strconv.Itoa(params.Offset))
	}

	if params.TakerOnly != nil {
		queryParams.Set("takerOnly", strconv.FormatBool(*params.TakerOnly))
	}

	if params.FilterType != "" {
		queryParams.Set("filterType", string(params.FilterType))
		queryParams.Set("filterAmount", params.FilterAmount.String())
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

	if params.User != "" {
		queryParams.Set("user", params.User)
	}

	if params.Side != "" {
		queryParams.Set("side", string(params.Side))
	}

	// Construct URL
	reqURL := fmt.Sprintf("%s/trades?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make trades request: %w", err)
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
	var trades []Trade
	if err := json.Unmarshal(body, &trades); err != nil {
		return nil, fmt.Errorf("failed to decode trades response: %w", err)
	}

	return trades, nil
}

// GetTradedMarketsCount retrieves the total number of markets a user has traded
func (c *DataClient) GetTradedMarketsCount(ctx context.Context, params *GetTradedMarketsCountParams) (*TradedMarketsCount, error) {
	if params.User == "" {
		return nil, fmt.Errorf("user address is required")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("user", params.User)

	// Construct URL
	reqURL := fmt.Sprintf("%s/traded?%s", Endpoint, queryParams.Encode())

	// Make request
	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make traded request: %w", err)
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
	var tradedCount TradedMarketsCount
	if err := json.Unmarshal(body, &tradedCount); err != nil {
		return nil, fmt.Errorf("failed to decode traded response: %w", err)
	}

	return &tradedCount, nil
}
