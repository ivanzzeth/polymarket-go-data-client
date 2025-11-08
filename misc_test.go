package polymarketdata

import (
	"context"
	"net/http"
	"testing"
)

func TestGetOpenInterest(t *testing.T) {
	client, err := NewClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetOpenInterestParams{
		Market: []string{"0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"},
	}

	openInterest, err := client.GetOpenInterest(context.Background(), params)
	if err != nil {
		t.Fatalf("GetOpenInterest failed: %v", err)
	}

	t.Logf("GetOpenInterest response (count: %d):", len(openInterest))
	for i, oi := range openInterest {
		if i < 3 {
			t.Logf("OpenInterest[%d]: %+v", i, oi)
		}
	}
}

func TestGetOpenInterestAll(t *testing.T) {
	client, err := NewClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetOpenInterestParams{}

	openInterest, err := client.GetOpenInterest(context.Background(), params)
	if err != nil {
		t.Fatalf("GetOpenInterest all failed: %v", err)
	}

	t.Logf("GetOpenInterest all response (count: %d):", len(openInterest))
	for i, oi := range openInterest {
		if i < 3 {
			t.Logf("OpenInterest[%d]: %+v", i, oi)
		}
	}
}

func TestGetLiveVolume(t *testing.T) {
	client, err := NewClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Using an example event ID
	params := &GetLiveVolumeParams{
		Id: 1,
	}

	liveVolume, err := client.GetLiveVolume(context.Background(), params)
	if err != nil {
		t.Fatalf("GetLiveVolume failed: %v", err)
	}

	t.Logf("GetLiveVolume response (count: %d):", len(liveVolume))
	for i, lv := range liveVolume {
		if i < 2 {
			t.Logf("LiveVolume[%d]: Total=%s, Markets count=%d", i, lv.Total.String(), len(lv.Markets))
			for j, market := range lv.Markets {
				if j < 3 {
					t.Logf("  Market[%d]: %+v", j, market)
				}
			}
		}
	}
}

func TestGetLiveVolumeMultipleEvents(t *testing.T) {
	client, err := NewClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with different event IDs
	for _, eventId := range []int{1, 10, 100} {
		params := &GetLiveVolumeParams{
			Id: eventId,
		}

		liveVolume, err := client.GetLiveVolume(context.Background(), params)
		if err != nil {
			t.Logf("GetLiveVolume for event %d failed: %v", eventId, err)
			continue
		}

		t.Logf("GetLiveVolume for event %d response (count: %d):", eventId, len(liveVolume))
		if len(liveVolume) > 0 {
			t.Logf("  First LiveVolume: Total=%s, Markets count=%d", liveVolume[0].Total.String(), len(liveVolume[0].Markets))
		}
	}
}
