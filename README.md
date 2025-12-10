# Go Echo Boilerplate 🚀

A robust, production-ready Go application boilerplate built with [Echo](https://echo.labstack.com/). Designed for modularity, developer experience, and extensibility.

## ✨ Features

-   **Modular Service Architecture**:
    -   Easily add new services in `internal/services/modules`.
    -   Enable/Disable services via `config.yaml` or Environment variables.
    -   Pre-loaded with multiple examples:
        -   `Service A`: User management demo.
        -   `Service B`: Product management demo (Disabled by default).
        -   `Service C`: In-Memory Cache demo.
        -   `Service D`: Task management using **GORM** (SQLite/Postgres).

-   **Advanced Monitoring Dashboard 📊** (New!):
    -   **Web UI**: Built with Shadcn-Admin style (TailwindCSS + Alpine.js).
    -   **Dashboard**: Live traffic logs (colorful!), Service count, Infrastructure health.
    -   **Infrastructure Stats**: Real-time status of Redis, Kafka, Postgres, Cron.
    -   **System Info**: Hostname, IP, Disk Usage.
    -   **Endpoints**: List all registered API endpoints.
    -   **Cron Jobs**: View scheduled jobs and their execution times.
    -   **Config Viewer**: Inspect running configuration.
    -   **Tools**: Redis Key Scanner, Postgres Query Monitor, Kafka Topic Debugger.
    -   **Banner Editor**: Edit the startup ASCII art from the browser.

-   **💎 Fancy Logger**:
    -   Built on [Zerolog](https://github.com/rs/zerolog).
    -   Rich console output with colors and emojis for better DX.
    -   Structured JSON logging ready for production.

-   **🛡️ Robust Middleware**:
    -   **Permission Guard**: Demonstration of strict permission blocking (e.g., Block all `DELETE` requests).
    -   **Request Logger**: Beautiful HTTP request logging with latency and status codes.

-   **🧠 In-Memory Cache**:
    -   Thread-safe, generic Key-Value store (`pkg/cache`).
    -   Built-in TTL (Time-To-Live) support.

-   **⏰ Cron Jobs**:
    -   Integrated `robfig/cron`.
    -   Configurable via `config.yaml`.

-   **🏭 Infrastructure Ready**:
    -   **Redis**: Integrated with `go-redis`.
    -   **Kafka**: Integrated with `sarama`.
    -   **Postgres**: Integrated with `pgx` and `GORM`.

## 🚀 Getting Started

### Prerequisites

-   Go 1.22+

### Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/GabrielTenma/bp-go-def.git
    cd bp-go-def
    ```

2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Run the application**:
    ```bash
    go run cmd/app/main.go
    ```

4.  **Access Monitoring**:
    Open `http://localhost:9090` (Default password: `admin`).

### Configuration (`config.yaml`)

```yaml
app:
  name: "My Fancy Go App"
  debug: true
  banner_path: "banner.txt"

monitoring:
  enabled: true
  port: "9090"
  password: "admin"
  username: "admin"

cron:
  enabled: true
  jobs:
    health_check: "*/10 * * * * *"

services:
  enable_service_a: true
  enable_service_b: false
  enable_service_c: true
  enable_service_d: true # Task Service (GORM)
```

## 📚 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/api/v1/users` | (Service A) Get dummy users |
| **GET** | `/api/v1/products`| (Service B) Get dummy products |
| **GET** | `/api/v1/cache/:key` | (Service C) Get cached value |
| **POST** | `/tasks` | (Service D) Create task |
| **GET** | `/tasks` | (Service D) List tasks |

## 🛠️ Project Structure

```
├── cmd/
│   └── app/            # Main entry point
├── config/             # Configuration logic (Viper)
├── internal/
│   ├── middleware/     # Custom middlewares
│   ├── monitoring/     # Monitoring Server & Handlers
│   ├── server/         # Server entry, DI, and startup logic
│   └── services/       # Business Logic
│       ├── modules/    # Individual service implementations
│       └── registry/   # Service Registration logic
├── pkg/
│   ├── cache/          # Generic In-Memory Cache
│   ├── infrastructure/ # External Infrastructure (Redis, Kafka, Postgres, Cron)
│   ├── logger/         # Custom Logger wrapper
│   └── utils/          # System utilities
└── web/
    └── monitoring/     # Frontend assets for Monitoring Dashboard
```
