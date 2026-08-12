package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/alteregoeth-ai/shipwatch-go/internal/monitor"
	"github.com/alteregoeth-ai/shipwatch-go/internal/store"
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
	mux.HandleFunc("POST /checks/run", api.runChecks)
	mux.HandleFunc("GET /checks", api.listChecks)
	return mux
}

func (api API) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (api API) runChecks(w http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 10*time.Second)
	defer cancel()

	results := make([]monitor.Result, 0, len(api.endpoints))
	for _, endpoint := range api.endpoints {
		result := api.checker.Check(ctx, endpoint)
		api.history.Add(result)
		results = append(results, result)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"availability": monitor.Availability(results),
		"results":      results,
	})
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

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
