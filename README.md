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

-   **Advanced Configuration**:
    -   Powered by [Viper](https://github.com/spf13/viper).
    -   Supports `config.yaml`, Environment Variables, and defaults.
    -   Customizable Server Port, Auth secrets, etc.

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
    -   Clean simple API: `Set(key, value, ttl)`, `Get(key)`.

-   **🐰 Customizable Banner**:
    -   Loads ASCII art from `banner.txt` at startup.
    -   Make your terminal startup fun!

-   **🏭 Infrastructure Ready**:
    -   **Redis**: Integrated with `go-redis` (v9). Includes helpers for `Set`, `Get`, `Delete`, `Replace`, `GetInfo`.
    -   **Kafka**: Integrated with `sarama` (Consumer & Producer). Includes `Consume` helper.
    -   **Postgres**: Integrated with `pgx` (v5). Includes `Select`, `Insert`, `Update`, `Delete` helpers.
    -   All infrastructure is optional and disabled by default.

## 🚀 Getting Started

### Prerequisites

-   Go 1.20+

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

### Configuration (`config.yaml`)

```yaml
app:
  name: "My Fancy Go App"
  debug: true
  banner_path: "banner.txt"

services:
  enable_service_a: true  # Users
  enable_service_b: false # Products
  enable_service_c: true  # Cache

# Infrastructure (Disabled by default)
redis:
  enabled: false
  address: "localhost:6379"
  password: ""
  db: 0

kafka:
  enabled: false
  brokers: 
    - "localhost:9092"
  topic: "my-topic"
  group_id: "my-group"

postgres:
  enabled: false
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "mydb"
  sslmode: "disable"
```

## 📚 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/api/v1/users` | (Service A) Get dummy users |
| **GET** | `/api/v1/products`| (Service B) Get dummy products (404 if disabled) |
| **GET** | `/api/v1/cache/:key` | (Service C) Get cached value |
| **POST** | `/api/v1/cache/:key` | (Service C) Set cached value |

## 🛠️ Project Structure

```
├── cmd/
│   └── app/            # Main entry point
├── config/             # Configuration logic (Viper)
├── internal/
│   ├── middleware/     # Custom middlewares (Auth, Logger, etc.)
│   ├── server/         # Server entry, DI, and startup logic
│   └── services/       # Business Logic
│       ├── modules/    # Individual service implementations
│       └── registry    # Service Registration logic
├── pkg/
│   ├── cache/          # Generic In-Memory Cache
│   ├── infrastructure/ # External Infrastructure (Redis, Kafka, Postgres)
│   └── logger/         # Custom Logger wrapper
└── docs/               # Architecture documentation
```
