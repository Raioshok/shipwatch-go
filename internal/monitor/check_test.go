package monitor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckerReportsHealthyEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	checker := NewChecker(server.Client())
	result := checker.Check(context.Background(), Endpoint{
		Name:           "orders",
		URL:            server.URL,
		Timeout:        time.Second,
		ExpectedStatus: http.StatusAccepted,
	})

	if !result.Healthy {
		t.Fatalf("expected endpoint to be healthy: %+v", result)
	}
	if result.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, result.StatusCode)
	}
}

func TestAvailabilityCalculatesHealthyRatio(t *testing.T) {
	results := []Result{
		{Name: "a", Healthy: true},
		{Name: "b", Healthy: false},
		{Name: "c", Healthy: true},
	}

	if got := Availability(results); got != 2.0/3.0 {
		t.Fatalf("expected 2/3 availability, got %f", got)
	}
}

func TestIncidentsReturnsUnhealthyResults(t *testing.T) {
	incidents := Incidents([]Result{
		{Name: "api", URL: "https://api.example.test", Healthy: false, Error: "timeout"},
		{Name: "docs", URL: "https://docs.example.test", Healthy: true},
	})

	if len(incidents) != 1 {
		t.Fatalf("expected one incident, got %d", len(incidents))
	}
	if incidents[0].Reason != "timeout" {
		t.Fatalf("expected timeout reason, got %q", incidents[0].Reason)
	}
}
