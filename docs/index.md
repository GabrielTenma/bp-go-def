---
layout: default
title: Go Echo Boilerplate - Production-Ready Golang Framework
description: A powerful Go boilerplate with Echo framework, monitoring dashboard, TUI, and complete infrastructure integrations
---

# Go Echo Boilerplate

A robust, production-ready Go application boilerplate built with [Echo](https://echo.labstack.com/). Designed for modularity, developer experience, and comprehensive monitoring with user management.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21%2B-00ADD8.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

## Features

### Core Application
-   **Modular Service Architecture**: Easy service extension with selective enable/disable
-   **Configurable Authentication**: API key-based auth with permission controls
-   **Fancy Logger**: Rich console output with colors, emojis, and structured logging (Zerolog)
-   **Custom ASCII Banner**: Configurable startup banner
-   **In-Memory Cache**: Thread-safe, generic KV store with TTL support
-   **Hot Configuration**: Update config without restart

### Terminal Interface
-   **Interactive Boot**: Visual boot sequence with service status checks
-   **Live CLI Dashboard**: Real-time terminal-based monitoring (Bubble Tea)
-   **Responsive TUI**: Adaptive layouts for different terminal sizes

### Infrastructure Support
-   **Redis**: Key-value store integration
-   **PostgreSQL**: SQL database with GORM
-   **Kafka**: Message queue integration
-   **Cron Jobs**: Scheduled task execution

### Monitoring Dashboard (Shadcn-Admin Style)
-   **Modern UI**: Beautiful shadcn-inspired design with Lexend font
-   **Dark Mode**: Full light/dark theme support with persistent storage
-   **Custom Login**: Shadcn-admin styled login page with HTTP Basic Auth
-   **User Settings**: Profile customization, photo upload, password management
-   **⚡ Live Metrics**: Real-time system stats (CPU, memory, disk, network)
-   **Live Logs**: SSE-based log streaming with color-coded levels
-   **Config Editor**: In-browser YAML editing with backup/restore
-   **Service Manager**: View all endpoints with active status badges
-   **Infrastructure Tools**: Redis browser, Postgres monitor, Kafka debugger
-   **Cron Monitor**: View scheduled jobs and execution status

## Getting Started

### Prerequisites
- Go 1.21+
- (Optional) Redis, PostgreSQL, Kafka

### Installation

```bash
# Clone the repository
git clone https://github.com/GabrielTenma/bp-go-def.git
cd bp-go-def

# Install dependencies
go mod download

# Run the application
go run cmd/app/main.go
```

### First Access

1. Start application
2. Open `http://localhost:9090`
3. Login with default password: `admin`
4. **Important**: Change password via User Settings immediately!

## Backend Console

![Backend Console](https://github.com/GabrielTenma/bp-go-def/raw/main/.assets/Recording%202025-12-14%20223856.gif)

## Monitoring Dashboard

![Monitoring Dashboard](https://github.com/GabrielTenma/bp-go-def/raw/main/.assets/Recording%202025-12-14%20230230.gif)

### Login Page
- Custom shadcn-admin styled design
- HTTP Basic Auth integration
- Dark mode support
- Responsive layout

### User Settings
- **Profile Photo**: Upload JPG/PNG/GIF (max 2MB)
- **Display Name**: Custom username
- **Password Management**: Change monitoring password
- **SQLite Storage**: Secure local database

### Dashboard Features
- Service overview with live counts
- Infrastructure status indicators
- System information (hostname, IP, disk usage)
- Live log streaming
- Categorized sidebar navigation

### Tools
- **Config Editor**: Edit YAML, backup, restore
- **Redis Browser**: Scan keys, view values
- **Postgres Monitor**: Active sessions, top queries
- **Kafka Debugger**: Topic inspection
- **Banner Editor**: Update ASCII art

## Terminal User Interface (TUI)

The application includes a rich terminal interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

### Features
-   **Boot Sequence**: Animated startup process showing service initialization status.
-   **Live Dashboard**: Monitor system resources (CPU, RAM, Goroutines) directly in the terminal.
-   **Interactive**: Keyboard controls for navigation and quitting.

For detailed implementation documentation, see [TUI_IMPLEMENTATION.md](https://github.com/GabrielTenma/bp-go-def/blob/main/docs_wiki/TUI_IMPLEMENTATION.md).

## Configuration

Edit `config.yaml`:

```yaml
app:
  name: "My Fancy Go App"
  debug: true
  env: "development"
  banner_path: "banner.txt"
  startup_delay: 15       # seconds to display boot screen (0 to skip)
  quiet_startup: true     # suppress console logs (TUI only)
  enable_tui: true        # enable fancy TUI mode

server:
  port: "8080"

services:
  enable_service_a: true
  enable_service_b: false
  enable_service_c: true
  enable_service_d: false

auth:
  type: "apikey"
  secret: "super-secret-key"

monitoring:
  enabled: true
  port: "9090"
  password: "admin"
  obfuscate_api: true
  title: "GoBP Admin"
  subtitle: "My Kisah Emuach ❤️"
  max_photo_size_mb: 2
  upload_dir: "web/monitoring/uploads"
  
  minio:
    enabled: true
    endpoint: "localhost:9003"
    access_key: "minioadmin"
    secret_key: "minioadmin"
    use_ssl: false
    bucket: "main"

  external:
    services:
      - name: "Google"
        url: "https://google.com"
      - name: "Soundcloud"
        url: "https://soundcloud.com"
      - name: "Local API"
        url: "http://localhost:8080/health"

redis:
  enabled: false
  address: "localhost:6379"
  password: ""
  db: 0

postgres:
  enabled: true
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "Mypostgres01"
  dbname: "postgres"
  sslmode: "disable"

kafka:
  enabled: false
  brokers: 
    - "localhost:9092"
  topic: "my-topic"
  group_id: "my-group"

cron:
  enabled: true
  jobs:
    log_cleanup: "0 0 * * *"
    health_check: "*/10 * * * * *" # Every 10 seconds
```

## Project Structure

```
.
├── cmd/app/              # Application entry point
├── config/               # Configuration logic
├── docs_wiki/            # Documentation & Guides
├── internal/
├── internal/
│   ├── middleware/       # Auth & permission middleware
│   ├── monitoring/       # Monitoring dashboard
│   │   ├── database/     # SQLite user management
│   │   ├── handlers.go   # API endpoints
│   │   └── server.go     # Monitoring server
│   ├── server/           # Main server logic
│   └── services/         # Service implementations
├── pkg/
│   ├── infrastructure/   # Redis, Postgres, Kafka, Cron
│   ├── logger/           # Rich console logger
│   ├── tui/              # Terminal User Interface
│   └── utils/            # System utilities
├── web/monitoring/       # Monitoring UI
│   ├── assets/          # CSS, JS
│   ├── login.html       # Login page
│   ├── index.html       # Dashboard
│   └── uploads/         # User files
└── config.yaml          # Main configuration
```

## API Endpoints

### Main Application
- `GET /health` - Health check
- `GET /api/service-a` - Service A
- `GET /api/service-c` - Service C
- `DELETE /api/*` - Blocked by middleware

### Monitoring APIs (Protected)
- `GET /api/status` - System status
- `GET /api/endpoints` - List services
- `GET /api/config` - Get config
- `POST /api/config` - Update config
- `GET /api/user/settings` - User profile
- `POST /api/user/password` - Change password
- `POST /api/user/photo` - Upload photo

## Security

- HTTP Basic Auth for monitoring
- BCrypt password hashing
- SQLite user database
- File upload size limits
- API key authentication
- Permission-based access control

## Development

### Adding a Service

1. Create in `internal/services/`
2. Add config flag in `config.yaml`
3. Register in `internal/server/server.go`
4. Auto-appears in monitoring!

### Database

- User settings: `monitoring_users.db` (auto-created)
- Default user from `config.yaml` password
- Change password via UI

## License

MIT

---

**Built with ❤️ using Go, Echo, Alpine.js, Tailwind CSS**
