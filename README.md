# api-gateway-go

Public REST API gateway for AI Twin reel orchestration. Accepts `POST /reels` requests, validates JWT auth, and publishes commands to SQS for downstream processing.

## Structure

- `cmd/server` — HTTP server entrypoint
- `internal/handlers` — route handlers (reels, runs)
- `internal/bus` — SQS publisher for reel commands
- `internal/models` — request/response types (generated from OpenAPI)

## Local dev

```bash
go run cmd/server/main.go
```

Server listens on `:8080` by default.
