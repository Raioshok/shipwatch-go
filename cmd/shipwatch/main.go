package main

import (
	"log"
	"os"
	"strings"

	"github.com/alteregoeth-ai/shipwatch-go/internal/config"
	"github.com/alteregoeth-ai/shipwatch-go/internal/monitor"
	"github.com/alteregoeth-ai/shipwatch-go/internal/server"
	"github.com/alteregoeth-ai/shipwatch-go/internal/store"
)

func main() {
	port := getenv("PORT", "8080")
	endpoints, err := loadEndpoints()
	if err != nil {
		log.Fatal(err)
	}
	history, err := store.NewMemoryHistory(getenv("HISTORY_FILE", "data/history.json"))
	if err != nil {
		log.Fatal(err)
	}

	api := server.NewAPI(endpoints, history)
	log.Printf("shipwatch-go listening on http://localhost:%s", port)
	log.Fatal(httpListenAndServe(":"+port, api.Handler()))
}

func loadEndpoints() ([]monitor.Endpoint, error) {
	if path := strings.TrimSpace(os.Getenv("CONFIG_FILE")); path != "" {
		return config.Load(path)
	}
	return config.FromTargets(getenv("TARGETS", "GitHub=https://github.com"))
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
