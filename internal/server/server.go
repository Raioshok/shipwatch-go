package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Raioshok/shipwatch-go/internal/monitor"
	"github.com/Raioshok/shipwatch-go/internal/store"
)

type API struct {
	checker   monitor.Checker
	history   *store.MemoryHistory
	endpoints []monitor.Endpoint
}

func NewAPI(endpoints []monitor.Endpoint, history *store.MemoryHistory) API {
	if history == nil {
		history = &store.MemoryHistory{}
	}
	return API{
		checker:   monitor.NewChecker(nil),
		history:   history,
		endpoints: endpoints,
	}
}

func (api API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", api.health)
	mux.HandleFunc("GET /endpoints", api.listEndpoints)
	mux.HandleFunc("POST /checks/run", api.runChecks)
	mux.HandleFunc("GET /checks", api.listChecks)
	mux.HandleFunc("GET /checks/latest", api.listChecks)
	mux.HandleFunc("GET /incidents", api.listIncidents)
	return mux
}

func (api API) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (api API) runChecks(w http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 10*time.Second)
	defer cancel()

	results := make([]monitor.Result, 0, len(api.endpoints))
	resultCh := make(chan monitor.Result, len(api.endpoints))
	var wg sync.WaitGroup
	for _, endpoint := range api.endpoints {
		wg.Add(1)
		go func(endpoint monitor.Endpoint) {
			defer wg.Done()
			resultCh <- api.checker.Check(ctx, endpoint)
		}(endpoint)
	}
	wg.Wait()
	close(resultCh)

	for result := range resultCh {
		api.history.Add(result)
		results = append(results, result)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"availability": monitor.Availability(results),
		"results":      results,
	})
}

func (api API) listEndpoints(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"endpoints": api.endpoints})
}

func (api API) listChecks(w http.ResponseWriter, _ *http.Request) {
	latest := api.history.LatestByEndpoint()
	results := make([]monitor.Result, 0, len(latest))
	for _, result := range latest {
		results = append(results, result)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"availability": monitor.Availability(results),
		"results":      results,
	})
}

func (api API) listIncidents(w http.ResponseWriter, _ *http.Request) {
	latest := api.history.LatestByEndpoint()
	results := make([]monitor.Result, 0, len(latest))
	for _, result := range latest {
		results = append(results, result)
	}
	writeJSON(w, http.StatusOK, map[string]any{"incidents": monitor.Incidents(results)})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
