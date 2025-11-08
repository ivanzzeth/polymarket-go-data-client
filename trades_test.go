package polymarketdata

import (
	"context"
	"net/http"
	"testing"
)

func TestGetTrades(t *testing.T) {
	client, err := NewClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetTradesParams{
		User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		Limit: 10,
	}

	trades, err := client.GetTrades(context.Background(), params)
	if err != nil {
		t.Fatalf("GetTrades failed: %v", err)
	}

	t.Logf("GetTrades response (count: %d):", len(trades))
	for i, trade := range trades {
		if i < 3 {
			t.Logf("Trade[%d]: %+v", i, trade)
		}
	}
}

func TestGetTradesWithSide(t *testing.T) {
	client, err := NewClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetTradesParams{
		User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		Side:  TradeSideBuy,
		Limit: 5,
	}

	trades, err := client.GetTrades(context.Background(), params)
	if err != nil {
		t.Fatalf("GetTrades with side failed: %v", err)
	}

	t.Logf("GetTrades with side response (count: %d):", len(trades))
	for i, trade := range trades {
		if i < 2 {
			t.Logf("Trade[%d]: %+v", i, trade)
		}
	}
}

func TestGetTradedMarketsCount(t *testing.T) {
	client, err := NewClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetTradedMarketsCountParams{
		User: "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
	}

	count, err := client.GetTradedMarketsCount(context.Background(), params)
	if err != nil {
		t.Fatalf("GetTradedMarketsCount failed: %v", err)
	}

	t.Logf("GetTradedMarketsCount response: %+v", count)
}
