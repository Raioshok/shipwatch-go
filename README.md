# ShipWatch Go

Small Go reliability service that checks configured endpoints concurrently, keeps bounded history, and derives incidents from the latest failed state. It uses the standard library so the concurrency and storage choices stay visible.

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

## Design Notes

- Goroutines fan checks out while a result channel gathers typed outcomes.
- Latest status and historical checks have separate query paths.
- File history is portable and inspectable; a production monitor would use a time-series store and durable scheduler.
