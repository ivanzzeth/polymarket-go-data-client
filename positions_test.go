package polymarketdata

import (
	"context"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
)

func TestGetPositions(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetPositionsParams{
		User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		Limit: 10,
	}

	positions, err := client.GetPositions(context.Background(), params)
	if err != nil {
		t.Fatalf("GetPositions failed: %v", err)
	}

	t.Logf("GetPositions response (count: %d):", len(positions))
	for i, pos := range positions {
		if i < 3 { // Print first 3 for brevity
			t.Logf("Position[%d]: %+v", i, pos)
		}
	}
}

func TestGetClosedPositions(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetClosedPositionsParams{
		User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		Limit: 10,
	}

	positions, err := client.GetClosedPositions(context.Background(), params)
	if err != nil {
		t.Fatalf("GetClosedPositions failed: %v", err)
	}

	t.Logf("GetClosedPositions response (count: %d):", len(positions))
	for i, pos := range positions {
		if i < 3 {
			t.Logf("ClosedPosition[%d]: %+v", i, pos)
		}
	}
}

func TestGetPositionsValue(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetValueParams{
		User: "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
	}

	values, err := client.GetPositionsValue(context.Background(), params)
	if err != nil {
		t.Fatalf("GetPositionsValue failed: %v", err)
	}

	t.Logf("GetPositionsValue response (count: %d):", len(values))
	for i, value := range values {
		if i < 3 {
			t.Logf("UserValue[%d]: %+v", i, value)
		}
	}
}

func TestGetPositionsWithThreshold(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	threshold := decimal.NewFromFloat(10.0)
	params := &GetPositionsParams{
		User:          "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		SizeThreshold: &threshold,
		Limit:         5,
	}

	positions, err := client.GetPositions(context.Background(), params)
	if err != nil {
		t.Fatalf("GetPositions with threshold failed: %v", err)
	}

	t.Logf("GetPositions with threshold response (count: %d):", len(positions))
	for i, pos := range positions {
		if i < 2 {
			t.Logf("Position[%d]: %+v", i, pos)
		}
	}
}
