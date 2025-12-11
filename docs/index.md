---
layout: default
title: Go Echo Boilerplate - Production-Ready Golang Framework
description: A powerful Go boilerplate with Echo framework, monitoring dashboard, and complete infrastructure integrations
---

# 🚀 Go Echo Boilerplate

A robust, production-ready Go application boilerplate built with [Echo](https://echo.labstack.com/). Designed for modularity, exceptional developer experience, and comprehensive monitoring.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21%2B-00ADD8.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

---

## ✨ Key Features

### 🏗️ Modular & Flexible Architecture
- **Modular Services**: Easily add new services with per-service enable/disable system
- **Clean Architecture**: Clear separation between handlers, services, and infrastructure layers
- **Plugin System**: Activate/deactivate features through configuration without code changes

### 📊 Beautiful Monitoring Dashboard
- **Modern UI**: Shadcn-admin inspired design with Lexend font
- **Dark Mode**: Light/dark theme support with persistent storage
- **Real-time Metrics**: Live monitoring of CPU, memory, disk, and network
- **Live Log Streaming**: Server-Sent Events (SSE) for real-time log streaming
- **User Management**: Profile customization, photo upload, password management

### 🔧 Complete Infrastructure Integration
- **Redis**: Key-value store with connection pooling
- **PostgreSQL**: SQL database with GORM ORM
- **Kafka**: Message queue for event-driven architecture
- **Cron Jobs**: Scheduled tasks with monitoring
- **MinIO**: Object storage for file uploads

### 📦 Standardized Request/Response Patterns
- **Consistent API**: Uniform response format across all endpoints
- **Auto Validation**: Automatic request validation with clear error messages
- **Built-in Pagination**: Pagination support with complete metadata
- **Type-Safe**: Uses structs for request/response instead of maps

### 🎨 Premium Developer Experience
- **Rich Logger**: Console output with colors, emojis, and structure (Zerolog)
- **Custom ASCII Banner**: Customizable banner on startup
- **Hot Config Reload**: Update configuration without application restart
- **Fancy Errors**: Informative and easy-to-debug error messages

### 🔒 Security First
- **API Key Authentication**: API key-based auth with permission control
- **HTTP Basic Auth**: For monitoring dashboard
- **BCrypt Password**: Secure password hashing
- **Permission Middleware**: Permission-based access control

---

## 🚀 Quick Start

```bash
# Clone repository
git clone https://github.com/GabrielTenma/bp-go-def.git
cd bp-go-def

# Install dependencies
go mod download

# Run the application
go run cmd/app/main.go
```

### First Access
1. Start the application
2. Open `http://localhost:9090` for monitoring
3. Login with default password: `admin`
4. **Important**: Change password via User Settings!

---

## 📚 Documentation

### Complete Guides
For detailed developer documentation, visit the [docs_wiki folder](https://github.com/GabrielTenma/bp-go-def/tree/main/docs_wiki):
- **Integration Guide** - How to use Redis, PostgreSQL, Kafka, Cron
- **Architecture Diagrams** - System flow and architecture diagrams
- **API Response Structure** - Complete response pattern documentation
- **Request/Response Structure** - Request validation guide

### Main Features

#### 🎯 Modular Service System
Add new services easily:
```yaml
# config.yaml
services:
  enable_service_a: true
  enable_service_b: false
```

Services are automatically registered and appear in the monitoring dashboard!

#### 📊 Monitoring Tools
- **Config Editor**: Edit YAML, backup, restore
- **Redis Browser**: Scan keys, view values
- **Postgres Monitor**: Active sessions, query statistics
- **Kafka Debugger**: Topic inspection
- **Cron Monitor**: Job scheduling and execution history

#### 📝 Standardized Response
```go
// Success
return response.Success(c, data, "User retrieved")

// Pagination
meta := response.CalculateMeta(page, perPage, total)
return response.SuccessWithMeta(c, users, meta)

// Error
return response.NotFound(c, "User not found")
```

---

## 🏗️ Project Structure

```
.
├── cmd/app/              # Application entry point
├── config/               # Configuration logic
├── internal/
│   ├── middleware/       # Auth & permission middleware
│   ├── monitoring/       # Monitoring dashboard
│   ├── server/           # Main server logic
│   └── services/         # Service implementations
├── pkg/
│   ├── infrastructure/   # Redis, Postgres, Kafka, Cron
│   ├── logger/           # Rich console logger
│   ├── request/          # Request validation
│   ├── response/         # Response helpers
│   └── utils/            # System utilities
├── web/monitoring/       # Monitoring UI
└── docs/                 # Developer documentation
```

---

## 🔌 API Endpoints

### Main Application
- `GET /health` - Health check
- `GET /api/service-a` - Service A endpoints
- `GET /api/service-c` - Service C endpoints
- `DELETE /api/*` - Blocked by permission middleware

### Monitoring APIs
- `GET /api/status` - System status & metrics
- `GET /api/endpoints` - List all services
- `GET /api/config` - Get/Update configuration
- `GET /api/user/settings` - User profile management

---

## ⚙️ Configuration

Edit `config.yaml` to customize the application:

```yaml
app:
  name: "My Fancy Go App"
  debug: true
  env: "development"

server:
  port: "8080"

services:
  enable_service_a: true
  enable_service_b: false

monitoring:
  enabled: true
  port: "9090"
  password: "admin"

redis:
  enabled: true
  address: "localhost:6379"

postgres:
  enabled: true
  host: "localhost"
  port: 5432

kafka:
  enabled: true
  brokers: ["localhost:9092"]
```

---

## 🎯 Use Cases

### Perfect For:
- ✅ REST API development with standardized patterns
- ✅ Microservices with integrated monitoring
- ✅ Event-driven applications with Kafka
- ✅ Applications requiring scheduled jobs
- ✅ Projects with multiple infrastructure dependencies

### Production Ready:
- ✅ Comprehensive error handling
- ✅ Structured logging
- ✅ Health checks
- ✅ Graceful shutdown
- ✅ Security best practices
- ✅ Configuration management

---

## 📸 Screenshots

![Monitoring Dashboard](https://s6.imgcdn.dev/YToc5C.png)

*Monitoring dashboard with shadcn-admin design and dark mode support*

---

## 🛠️ Tech Stack

**Backend:**
- [Echo](https://echo.labstack.com/) - High performance web framework
- [Zerolog](https://github.com/rs/zerolog) - Zero allocation JSON logger
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Validator](https://github.com/go-playground/validator) - Request validation

**Frontend (Monitoring):**
- [Alpine.js](https://alpinejs.dev/) - Lightweight JavaScript framework
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS
- [Chart.js](https://www.chartjs.org/) - Beautiful charts
- [CodeMirror](https://codemirror.net/) - Code editor

**Infrastructure:**
- [Redis](https://redis.io/) - In-memory data store
- [PostgreSQL](https://www.postgresql.org/) - SQL database
- [Kafka](https://kafka.apache.org/) - Event streaming
- [MinIO](https://min.io/) - S3-compatible storage

---

## 📝 License

MIT License - feel free to use for commercial or personal projects!

---

## 🤝 Contributing

Contributions are very welcome! Please create an issue or pull request.

---

## 🔗 Links

- **GitHub Repository**: [https://github.com/GabrielTenma/bp-go-def](https://github.com/GabrielTenma/bp-go-def)
- **Documentation**: [https://gabrieltenma.github.io/bp-go-def/](https://gabrieltenma.github.io/bp-go-def/)
- **Developer Docs**: [/docs_wiki](https://github.com/GabrielTenma/bp-go-def/tree/main/docs_wiki)

---

**Built with ❤️ using Go, Echo, Alpine.js, and Tailwind CSS**
