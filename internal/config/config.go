package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Raioshok/shipwatch-go/internal/monitor"
)

type File struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	ExpectedStatus int    `json:"expected_status"`
}

func Load(path string) ([]monitor.Endpoint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}

	endpoints := make([]monitor.Endpoint, 0, len(file.Endpoints))
	for _, item := range file.Endpoints {
		endpoints = append(endpoints, monitor.Endpoint{
			Name:           item.Name,
			URL:            item.URL,
			Timeout:        time.Duration(item.TimeoutSeconds) * time.Second,
			ExpectedStatus: item.ExpectedStatus,
		})
	}
	return endpoints, nil
}

func FromTargets(targets string) ([]monitor.Endpoint, error) {
	parts := strings.Split(targets, ",")
	endpoints := make([]monitor.Endpoint, 0, len(parts))
	for _, target := range parts {
		pair := strings.SplitN(target, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid target %q, expected name=url", target)
		}
		endpoints = append(endpoints, monitor.Endpoint{
			Name:           strings.TrimSpace(pair[0]),
			URL:            strings.TrimSpace(pair[1]),
			Timeout:        3 * time.Second,
			ExpectedStatus: 200,
		})
	}
	return endpoints, nil
}
