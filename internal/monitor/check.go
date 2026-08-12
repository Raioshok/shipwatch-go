package monitor

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Endpoint struct {
	Name           string        `json:"name"`
	URL            string        `json:"url"`
	Timeout        time.Duration `json:"timeout"`
	ExpectedStatus int           `json:"expected_status"`
}

type Result struct {
	Name       string        `json:"name"`
	URL        string        `json:"url"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
	Healthy    bool          `json:"healthy"`
	CheckedAt  time.Time     `json:"checked_at"`
	Error      string        `json:"error,omitempty"`
}

type Incident struct {
	Endpoint  string    `json:"endpoint"`
	URL       string    `json:"url"`
	OpenedAt time.Time `json:"opened_at"`
	Reason   string    `json:"reason"`
}

type Checker struct {
	client *http.Client
	now    func() time.Time
}

func NewChecker(client *http.Client) Checker {
	if client == nil {
		client = http.DefaultClient
	}
	return Checker{client: client, now: time.Now}
}

func (c Checker) Check(ctx context.Context, endpoint Endpoint) Result {
	start := c.now()
	timeout := endpoint.Timeout
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	expectedStatus := endpoint.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = http.StatusOK
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.URL, nil)
	if err != nil {
		return failed(endpoint, start, c.now().Sub(start), err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return failed(endpoint, start, c.now().Sub(start), err)
	}
	defer response.Body.Close()

	latency := c.now().Sub(start)
	return Result{
		Name:       endpoint.Name,
		URL:        endpoint.URL,
		StatusCode: response.StatusCode,
		Latency:    latency,
		Healthy:    response.StatusCode == expectedStatus,
		CheckedAt:  start,
	}
}

func failed(endpoint Endpoint, checkedAt time.Time, latency time.Duration, err error) Result {
	if err == nil {
		err = errors.New("unknown error")
	}
	return Result{
		Name:      endpoint.Name,
		URL:       endpoint.URL,
		Latency:   latency,
		Healthy:   false,
		CheckedAt: checkedAt,
		Error:     err.Error(),
	}
}

func Availability(results []Result) float64 {
	if len(results) == 0 {
		return 0
	}
	healthy := 0
	for _, result := range results {
		if result.Healthy {
			healthy++
		}
	}
	return float64(healthy) / float64(len(results))
}

func Incidents(results []Result) []Incident {
	incidents := make([]Incident, 0)
	for _, result := range results {
		if result.Healthy {
			continue
		}
		reason := result.Error
		if reason == "" {
			reason = "unexpected status"
		}
		incidents = append(incidents, Incident{
			Endpoint:  result.Name,
			URL:       result.URL,
			OpenedAt: result.CheckedAt,
			Reason:   reason,
		})
	}
	return incidents
}
