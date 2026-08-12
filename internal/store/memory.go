package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/Raioshok/shipwatch-go/internal/monitor"
)

type MemoryHistory struct {
	mu      sync.RWMutex
	results []monitor.Result
	path    string
}

func NewMemoryHistory(path string) (*MemoryHistory, error) {
	history := &MemoryHistory{path: path}
	if path == "" {
		return history, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return history, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return history, nil
	}
	if err := json.Unmarshal(data, &history.results); err != nil {
		return nil, err
	}
	return history, nil
}

func (h *MemoryHistory) Add(result monitor.Result) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.results = append(h.results, result)
	_ = h.persistLocked()
}

func (h *MemoryHistory) List() []monitor.Result {
	h.mu.RLock()
	defer h.mu.RUnlock()
	results := append([]monitor.Result(nil), h.results...)
	sort.Slice(results, func(i, j int) bool {
		return results[i].CheckedAt.After(results[j].CheckedAt)
	})
	return results
}

func (h *MemoryHistory) LatestByEndpoint() map[string]monitor.Result {
	h.mu.RLock()
	defer h.mu.RUnlock()

	latest := make(map[string]monitor.Result)
	for _, result := range h.results {
		existing, ok := latest[result.Name]
		if !ok || result.CheckedAt.After(existing.CheckedAt) {
			latest[result.Name] = result
		}
	}
	return latest
}

func (h *MemoryHistory) persistLocked() error {
	if h.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(h.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(h.results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, append(data, '\n'), 0o644)
}
