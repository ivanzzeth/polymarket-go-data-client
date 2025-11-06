package polymarketdata

import (
	"context"
	"net/http"
	"testing"
)

func TestGetHolders(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Using an example market condition ID
	params := &GetHoldersParams{
		Market: []string{"0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"},
		Limit:  10,
	}

	holders, err := client.GetHolders(context.Background(), params)
	if err != nil {
		t.Fatalf("GetHolders failed: %v", err)
	}

	t.Logf("GetHolders response (count: %d):", len(holders))
	for i, marketHolder := range holders {
		if i < 2 {
			t.Logf("MarketHolders[%d]: Token=%s, Holders count=%d", i, marketHolder.Token, len(marketHolder.Holders))
			for j, holder := range marketHolder.Holders {
				if j < 2 {
					t.Logf("  Holder[%d]: %+v", j, holder)
				}
			}
		}
	}
}

func TestGetHoldersWithMinBalance(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetHoldersParams{
		Market:     []string{"0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"},
		MinBalance: 100,
		Limit:      5,
	}

	holders, err := client.GetHolders(context.Background(), params)
	if err != nil {
		t.Fatalf("GetHolders with min balance failed: %v", err)
	}

	t.Logf("GetHolders with min balance response (count: %d):", len(holders))
	for i, marketHolder := range holders {
		if i < 1 {
			t.Logf("MarketHolders[%d]: Token=%s, Holders count=%d", i, marketHolder.Token, len(marketHolder.Holders))
			for j, holder := range marketHolder.Holders {
				if j < 2 {
					t.Logf("  Holder[%d]: %+v", j, holder)
				}
			}
		}
	}
}
