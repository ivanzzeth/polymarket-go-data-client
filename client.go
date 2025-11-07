package polymarketdata

import (
	"context"
	"net/http"
)

const Endpoint = "https://data-api.polymarket.com"

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) (*Client, error) {
	return &Client{
		httpClient: httpClient,
	}, nil
}

// doRequest is a helper method to make HTTP requests with context
func (c *Client) doRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}
