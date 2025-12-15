# Welcome to the Project Documentation

This wiki serves as the central knowledge base for the project, covering architecture, feature implementations, and integration guides.

## Core Documentation

### Architecture & Design

*   **[Architecture Diagrams](ARCHITECTURE_DIAGRAMS.md)**
    *   Visual guide to the system's request/response flow.
    *   Package organization and dependency graphs.
    *   Sequence diagrams for validation and error handling.

*   **[API Response Structure](API_RESPONSE_STRUCTURE.md)**
    *   Standard format for all API responses (`success`, `data`, `meta`).
    *   Built-in helper functions for success and error responses.
    *   Pagination and validation standards.

*   **[Request Response Structure](REQUEST_RESPONSE_STRUCTURE.md)**
    *   Detailed overview of the Echo service structure.
    *   Request validation patterns and custom validators.
    *   Dependencies and best practices for creating new endpoints.

### Security & Privacy

*   **[API Obfuscation](API_OBFUSCATION.md)**
    *   Mechanism for obscuring JSON data in transit using Base64.
    *   Configuration guide for enabling/disabling obfuscation.
    *   Frontend and Backend implementation details.

### User Interface

*   **[TUI Implementation](TUI_IMPLEMENTATION.md)**
    *   Documentation for the Terminal User Interface (Bubble Tea).
    *   **Boot Sequence**: Visual feedback during service initialization.
    *   **Live Dashboard**: Real-time monitoring of system resources and service health.

### Integration & Infrastructure

*   **[Service Implementation Guide](SERVICE_IMPLEMENTATION.md)**
    *   Creating and implementing new services.
    *   Service interface and registration.
    *   Dynamic configuration setup.

*   **[Integration Guide](INTEGRATION_GUIDE.md)**
    *   **Redis**: Configuration and usage of the Redis manager.
    *   **Postgres**: Database connection and Helper methods.
    *   **Kafka**: Message producing and configuration.
    *   **MinIO**: Object storage integration for file uploads.
    *   **Cron Jobs**: Dynamic job scheduling and management.

---

## Getting Started

If you are new to the project, we recommend starting with the **[Integration Guide](INTEGRATION_GUIDE.md)** to understand the available infrastructure components, followed by the **[API Response Structure](API_RESPONSE_STRUCTURE.md)** to learn how to build consistent API endpoints.
