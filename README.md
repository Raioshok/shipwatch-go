# ShipWatch Go

Go reliability monitor for checking business-critical endpoints and reporting availability. It demonstrates Go services, CLI configuration, HTTP monitoring, concurrency-safe in-memory storage, tests, Docker, and CI.

## Features

- Configurable endpoint checks through `TARGETS`
- Health-check API
- On-demand monitoring run
- Latest status by endpoint
- Availability calculation
- Standard-library implementation with focused tests

## Run locally

```bash
go run ./cmd/shipwatch
```

```bash
TARGETS="GitHub=https://github.com,Google=https://google.com" go run ./cmd/shipwatch
```

## API

- `GET /health`
- `POST /checks/run`
- `GET /checks`

## Test

```bash
go test ./...
```

## Portfolio talking points

- Go is a strong fit for infrastructure, backend services, cloud tooling, and SRE automation.
- The repo shows practical reliability engineering, not only syntax.
- The service uses only the standard library, which makes the core design easy to review.
