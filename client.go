package polymarketdata

import (
	"context"
	"net/http"
)

const Endpoint = "https://data-api.polymarket.com"

type DataClient struct {
	httpClient *http.Client
}

func NewDataClient(httpClient *http.Client) (*DataClient, error) {
	return &DataClient{
		httpClient: httpClient,
	}, nil
}

// doRequest is a helper method to make HTTP requests with context
func (c *DataClient) doRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}
