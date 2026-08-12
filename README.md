# ShipWatch Go

Go reliability monitor for checking business-critical endpoints and reporting availability. It demonstrates Go services, CLI configuration, HTTP monitoring, concurrency-safe in-memory storage, tests, Docker, and CI.

## Features

- Configurable endpoint checks through `TARGETS`
- JSON config loading through `CONFIG_FILE`
- Health-check API
- On-demand monitoring run
- Concurrent endpoint checks
- File-backed check history
- Latest status by endpoint
- Incident detection for unhealthy latest checks
- Availability calculation
- Standard-library implementation with focused tests

## Run locally

```bash
go run ./cmd/shipwatch
```

```bash
TARGETS="GitHub=https://github.com,Google=https://google.com" go run ./cmd/shipwatch
```

```bash
CONFIG_FILE=config/endpoints.json go run ./cmd/shipwatch
```

## API

- `GET /health`
- `GET /endpoints`
- `POST /checks/run`
- `GET /checks`
- `GET /checks/latest`
- `GET /incidents`

See `docs/API.md` for details.

## Test

```bash
go test ./...
```

## Portfolio talking points

- Go is a strong fit for infrastructure, backend services, cloud tooling, and SRE automation.
- The repo shows practical reliability engineering, not only syntax.
- The service uses only the standard library, which makes the core design easy to review.
