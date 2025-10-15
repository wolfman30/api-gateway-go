# api-gateway-go

Public REST API gateway for AI Twin reel orchestration. Accepts `POST /reels` requests, validates JWT auth, and publishes commands to SQS for downstream processing.

## Structure

- `cmd/server` — HTTP server entrypoint
- `internal/handlers` — route handlers (reels, runs)
- `internal/bus` — SQS publisher for reel commands
- `internal/models` — request/response types (aligned with OpenAPI schema from ai-twin-contracts)

## Local dev

### Running the server

```bash
go run cmd/server/main.go
```

Server listens on `:8081` by default.

### Environment variables

- `SQS_QUEUE_URL` — AWS SQS queue URL for publishing reel commands (defaults to stub if unset)

### Running tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run tests for a specific package
go test ./internal/handlers -v
```

### Testing without a live server

Unit tests use `httptest` to validate handler logic without binding to a port. See `internal/handlers/reels_test.go` for examples.

## API Endpoints

### `POST /reels`

Submit a new reel generation request.

**Request body**: JSON matching `CreateReelRequest` schema (see `internal/models/types.go`)

**Response**: `202 Accepted` with `{"runId": "uuid"}`

**Example**:
```bash
curl -X POST http://localhost:8081/reels \
  -H "Content-Type: application/json" \
  -d @../ai-twin-contracts/examples/create_reel_request.v1.json
```

### `GET /runs/{runId}`

Fetch the current status of a reel run.

**Response**: JSON with run status and step details

### `GET /health`

Health check endpoint.

**Response**: `200 OK` with `"OK"`

## Dependencies

- `github.com/google/uuid` — UUID generation for run IDs
