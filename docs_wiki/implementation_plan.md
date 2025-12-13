# Go Echo Modular App Implementation Plan

## Goal Description
Create a production-ready, extensible Go application skeleton using `labstack/echo`. Key features include a modular service architecture (load services selectively), configurable middleware (Auth, etc.), and a developer-friendly console output with rich colors and icons.

## User Review Required
> [!NOTE]
> I will be using standard Go project layout norms.
> The "fancy console" will rely on standard ANSI escape codes (via libraries) for colors to ensure compatibility across terminals.

## Proposed Changes

### Core Structure
#### [NEW] `cmd/app`
- Entry point `main.go`.

#### [NEW] `config`
- `config.go`: Structs for loading configuration (likely from `.env` or YAML).
- Settings for toggling services (e.g., `EnableServiceA: true`).

#### [NEW] `pkg/logger`
- Custom logger wrapper.
- Features: Info, Error, Debug, Warn with icons and colors.

#### [NEW] `internal/server`
- Echo server initialization.
- Middleware registration logic.
- Dynamic route binding based on enabled services.

#### [NEW] `internal/services`
- Base Service interface (if needed) or simple modular structure.
- Example Service A & Service B.

### Middleware
#### [NEW] `internal/middleware`
- Configurable Auth middleware (e.g., JWT, API Key) with a switch in config.
- Logging middleware integration.

## Verification Plan

### Automated Tests
- Run `go run ./cmd/app` and verify startup logs.
- curl endpoints for enabled/disabled services to verify routing logic.

### Manual Verification
- Check console output for "fancy" formatting.
- Toggle services in config and restart to verify they load/unload.
