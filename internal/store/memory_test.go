package store

import (
	"testing"
	"time"

	"github.com/Raioshok/shipwatch-go/internal/monitor"
)

func TestLatestByEndpointReturnsMostRecentResult(t *testing.T) {
	history := &MemoryHistory{}
	history.Add(monitor.Result{Name: "api", Healthy: false, CheckedAt: time.Date(2026, 8, 12, 10, 0, 0, 0, time.UTC)})
	history.Add(monitor.Result{Name: "api", Healthy: true, CheckedAt: time.Date(2026, 8, 12, 10, 1, 0, 0, time.UTC)})

	latest := history.LatestByEndpoint()

	if !latest["api"].Healthy {
		t.Fatalf("expected latest api result to be healthy")
	}
}

func TestNewMemoryHistoryLoadsPersistedResults(t *testing.T) {
	path := t.TempDir() + "/history.json"
	history, err := NewMemoryHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	history.Add(monitor.Result{Name: "api", Healthy: true, CheckedAt: time.Date(2026, 8, 12, 10, 0, 0, 0, time.UTC)})

	reloaded, err := NewMemoryHistory(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(reloaded.List()) != 1 {
		t.Fatalf("expected one persisted result")
	}
}
