package polymarketdata

import (
	"net/http"
	"testing"
)

func TestGetActivity(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetActivityParams{
		User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		Limit: 10,
	}

	activities, err := client.GetActivity(params)
	if err != nil {
		t.Fatalf("GetActivity failed: %v", err)
	}

	t.Logf("GetActivity response (count: %d):", len(activities))
	for i, activity := range activities {
		if i < 3 {
			t.Logf("Activity[%d]: %+v", i, activity)
		}
	}
}

func TestGetActivityWithType(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetActivityParams{
		User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		Type:  []ActivityType{ActivityTypeTrade},
		Limit: 5,
	}

	activities, err := client.GetActivity(params)
	if err != nil {
		t.Fatalf("GetActivity with type failed: %v", err)
	}

	t.Logf("GetActivity with type response (count: %d):", len(activities))
	for i, activity := range activities {
		if i < 2 {
			t.Logf("Activity[%d]: %+v", i, activity)
		}
	}
}

func TestGetActivityWithTimeRange(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &GetActivityParams{
		User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		Start: 1640000000, // Example timestamp
		Limit: 5,
	}

	activities, err := client.GetActivity(params)
	if err != nil {
		t.Fatalf("GetActivity with time range failed: %v", err)
	}

	t.Logf("GetActivity with time range response (count: %d):", len(activities))
	for i, activity := range activities {
		if i < 2 {
			t.Logf("Activity[%d]: %+v", i, activity)
		}
	}
}
