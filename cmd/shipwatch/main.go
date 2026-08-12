package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alteregoeth-ai/shipwatch-go/internal/monitor"
	"github.com/alteregoeth-ai/shipwatch-go/internal/server"
)

func main() {
	port := getenv("PORT", "8080")
	targets := strings.Split(getenv("TARGETS", "GitHub=https://github.com"), ",")

	endpoints := make([]monitor.Endpoint, 0, len(targets))
	for _, target := range targets {
		parts := strings.SplitN(target, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("invalid TARGETS entry %q, expected name=url", target)
		}
		endpoints = append(endpoints, monitor.Endpoint{
			Name:           strings.TrimSpace(parts[0]),
			URL:            strings.TrimSpace(parts[1]),
			Timeout:        3 * time.Second,
			ExpectedStatus: http.StatusOK,
		})
	}

	api := server.NewAPI(endpoints, nil)
	log.Printf("shipwatch-go listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, api.Handler()))
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
