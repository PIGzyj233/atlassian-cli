# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
# Build both binaries (Windows needs .exe suffix)
go build -ldflags "-s -w -X main.version=dev" -o bin/jira.exe ./cmd/jira
go build -ldflags "-s -w -X main.version=dev" -o bin/confluence.exe ./cmd/confluence

# Run all tests (skip -race on Windows without CGO)
go test ./...

# Run a single package's tests
go test ./internal/api/...
go test ./pkg/jira/issue/...

# Run a single test by name
go test ./internal/api/... -run TestClient_Get

# Lint and format
golangci-lint run ./...
gofmt -s -w .
```

The Makefile exists but `make` may not be available in all environments. Use `go` commands directly.

## Architecture

Two independent CLI binaries (`cmd/jira/main.go`, `cmd/confluence/main.go`) share a common infrastructure layer:

```
cmd/{jira,confluence}/main.go     ← Entry points: create Factory, build root command, Execute()
        │
pkg/cmdutil/factory.go            ← Factory: holds Config, creates authenticated API clients
        │
pkg/{jira,confluence}/root.go     ← Root cobra commands, register all subcommand groups
pkg/{jira,confluence}/<feature>/   ← One package per feature (issue/, board/, page/, etc.)
        │
internal/api/client.go            ← HTTP client with retry (3x exponential backoff on 429/5xx)
internal/api/request.go           ← RequestBuilder: fluent query-param construction
internal/api/errors.go            ← APIError type with status code parsing
internal/api/pagination.go        ← Pagination helpers
        │
internal/config/                  ← YAML config (~/.config/atlassian-cli/config.yml)
internal/config/auth.go           ← Authenticator interface: BasicAuth, PATAuth
internal/config/hosts.go          ← Host resolution (flag → default → single-host)
        │
internal/output/                  ← Formatter interface: JSON, Table, Text
internal/convert/                 ← ADF ↔ storage format conversion (Confluence)
```

### Key patterns

- **Factory injection**: Every command receives a `*cmdutil.Factory`. Call `f.JiraClient(flagHost)` or `f.ConfluenceClient(flagHost)` to get an authenticated `*api.Client`.
- **Cloud vs Server**: `client.IsCloud()` detects `.atlassian.net` hosts. Jira Cloud uses REST API v3 (`/rest/api/3/`), Server/DC uses v2 (`/rest/api/2/`). Confluence always uses `/rest/api/`.
- **Command structure**: Each feature is a package under `pkg/jira/` or `pkg/confluence/`. A `NewCmd*` function returns a `*cobra.Command` with `RunE` for error handling. Subcommands are added in the feature package.
- **Output**: Commands call `output.Print(cmd, data)` which reads the `--output` flag and delegates to the appropriate formatter.
- **Global flags**: `--host`, `--output` (json|text|table), `--verbose`, `--no-color`.

### Config file format

YAML at `%APPDATA%/atlassian-cli/config.yml` (Windows) or `~/.config/atlassian-cli/config.yml` (Unix). Per-host entries specify instance type (cloud/server/datacenter), auth method (basic/pat), and per-service settings.

## Testing conventions

- Uses `github.com/stretchr/testify` (assert/require).
- E2E tests (`e2e_test.go`) spin up `httptest.NewServer` mock servers, write temp config, build the binary with `go build`, and execute it as a subprocess.
- The `exeName()` helper in e2e tests appends `.exe` on Windows via `runtime.GOOS`.

## Release

GoReleaser v2 config in `.goreleaser.yml`. CI triggers on `v*` tags via `.github/workflows/release.yml`. Builds for linux/darwin/windows × amd64/arm64 with CGO disabled.
