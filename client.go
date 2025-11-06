package polymarketdata

import "net/http"

const Endpoint = "https://data-api.polymarket.com"

type DataClient struct {
	httpClient *http.Client
}

func NewDataClient(httpClient *http.Client) (*DataClient, error) {
	return &DataClient{
		httpClient: httpClient,
	}, nil
}
