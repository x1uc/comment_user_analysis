package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_RateLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("test-cookie")
	rate := 200 * time.Millisecond
	c.SetRateLimit(rate)

	start := time.Now()
	
	// First request - should be immediate (mostly)
	_, err := c.Get(ts.URL)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	
	// Second request - should wait
	_, err = c.Get(ts.URL)
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	
	duration := time.Since(start)
	if duration < rate {
		t.Errorf("Expected duration >= %v, got %v", rate, duration)
	}
}
