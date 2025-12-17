# bp-go-def

A robust, production-ready Go application boilerplate built with [Echo](https://echo.labstack.com/). Features modular architecture, comprehensive monitoring, and built-in infrastructure support.

## Quick Start

### Prerequisites
- Go 1.21+

### Installation & Run

```bash
# Clone the repository
git clone https://github.com/GabrielTenma/bp-go-def.git
cd bp-go-def

# Install dependencies
go mod download

# Run the application
go run cmd/app/main.go
```

**First Access:**
1. Open `http://localhost:9090` (monitoring dashboard)
2. Login with password: `admin`
3. **Important**: Change the default password immediately!

## Screenshots

### Backend Console
![Backend Console](.assets/Recording%202025-12-14%20223856.gif)

### Monitoring Dashboard
![Monitoring Dashboard](.assets/Recording%202025-12-14%20230230.gif)

## Key Features

- **Modular Services**: Enable/disable services via configuration
- **Monitoring Dashboard**: Real-time metrics, logs, and system monitoring
- **Terminal UI**: Interactive boot sequence and live CLI dashboard
- **Infrastructure Support**: Redis, PostgreSQL (multi-tenant), Kafka, MinIO
- **Security**: API encryption, authentication, and access controls
- **Build Tools**: Automated build scripts with backup and archiving

## Documentation

📖 **[Full Documentation](docs_wiki/)** - Comprehensive guides and references

### Core Documentation
- **[Configuration Guide](docs_wiki/CONFIGURATION_GUIDE.md)** - Complete configuration reference
- **[API Response Structure](docs_wiki/API_RESPONSE_STRUCTURE.md)** - Standard response formats
- **[Architecture Diagrams](docs_wiki/ARCHITECTURE_DIAGRAMS.md)** - System design and flow diagrams
- **[Service Implementation](docs_wiki/SERVICE_IMPLEMENTATION.md)** - How to add new services

### Infrastructure & Integration
- **[Integration Guide](docs_wiki/INTEGRATION_GUIDE.md)** - Redis, PostgreSQL, Kafka, MinIO setup
- **[Build Scripts](docs_wiki/BUILD_SCRIPTS.md)** - Production deployment automation
- **[Package Management](docs_wiki/CHANGE_PACKAGE_SCRIPTS.md)** - Module renaming tools

### Security & Features
- **[API Encryption](docs_wiki/ENCRYPTION_API.md)** - End-to-end encryption
- **[API Obfuscation](docs_wiki/API_OBFUSCATION.md)** - Data obfuscation mechanisms
- **[TUI Implementation](docs_wiki/TUI_IMPLEMENTATION.md)** - Terminal interface details

## Project Structure

```
.
├── cmd/app/              # Application entry point
├── config/               # Configuration logic
├── docs_wiki/            # 📖 Documentation
├── internal/
│   ├── middleware/       # Auth & security middleware
│   ├── monitoring/       # Web monitoring dashboard
│   ├── server/           # Main server logic
│   └── services/         # Modular service implementations
├── pkg/
│   ├── infrastructure/   # Redis, Postgres, Kafka, etc.
│   ├── logger/           # Rich console logging
│   ├── tui/              # Terminal User Interface
│   └── utils/            # System utilities
├── web/monitoring/       # Monitoring web UI
└── config.yaml           # Main configuration
```

## License

MIT

---

**Built with ❤️ using Go, Echo, Alpine.js, Tailwind CSS**
