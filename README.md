# Go Echo Boilerplate

A robust, production-ready Go application boilerplate built with [Echo](https://echo.labstack.com/). Designed for modularity, developer experience, and comprehensive monitoring with user management.

## Features

### Core Application
-   **Modular Service Architecture**: Easy service extension with selective enable/disable
-   **Configurable Authentication**: API key-based auth with permission controls
-   **Fancy Logger**: Rich console output with colors, emojis, and structured logging (Zerolog)
-   **Custom ASCII Banner**: Configurable startup banner
-   **In-Memory Cache**: Thread-safe, generic KV store with TTL support
-   **Hot Configuration**: Update config without restart

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
-   **� Live Metrics**: Real-time system stats (CPU, memory, disk, network)
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
git clone <repository-url>
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

## Demo

<p>
    <img src="https://s6.imgcdn.dev/YTovOt.gif" width="60%" alt="Demo">
</p>

## Monitoring Dashboard

![Monitoring Dashboard](https://s6.imgcdn.dev/YToc5C.png)

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

## Configuration

Edit `config.yaml`:

```yaml
app:
  name: "My Fancy Go App"
  debug: true
  env: "development"
  banner_path: "banner.txt"

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
  password: "admin"  # Initial password (change via UI!)
  title: "GoBP Admin"
  subtitle: "My Custom Subtitle"
  max_photo_size_mb: 2
  upload_dir: "web/monitoring/uploads"

redis:
  enabled: false
  address: "localhost:6379"

postgres:
  enabled: false
  host: "localhost"
  port: 5432

kafka:
  enabled: false
  brokers: ["localhost:9092"]

cron:
  enabled: true
  jobs:
    log_cleanup: "0 0 * * *"
    health_check: "*/10 * * * * *"
```

## Project Structure

```
.
├── cmd/app/              # Application entry point
├── config/               # Configuration logic
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
