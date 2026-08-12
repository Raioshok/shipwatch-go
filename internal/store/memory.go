package store

import (
	"sort"
	"sync"

	"github.com/alteregoeth-ai/shipwatch-go/internal/monitor"
)

type MemoryHistory struct {
	mu      sync.RWMutex
	results []monitor.Result
}

func (h *MemoryHistory) Add(result monitor.Result) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.results = append(h.results, result)
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
