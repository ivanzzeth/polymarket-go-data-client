package polymarketdata

import (
	"context"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	client, err := NewDataClient(&http.Client{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	resp, err := client.HealthCheck(context.Background(), )
	if err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}

	t.Logf("HealthCheck response: %+v", resp)

	if resp.Data == "" {
		t.Error("Expected non-empty Data field")
	}
}
