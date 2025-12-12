# Go Echo Modular App Walkthrough

This application is a modular web server built with `labstack/echo`, featuring a configurable architecture, fancy logger, and permission enforcement.

## Features Implemented

1.  **Modular Service Architecture**:
    -   Services can be enabled/disabled via `config.yaml`.
    -   **Refactored**: Service implementations (`Service A`, `Service B`) are located in `internal/services/modules` for better organization.
    -   Currently: `Service A` (Users) is ENABLED, `Service B` (Products) is DISABLED.

2.  **Configurable Middleware**:
    -   **Permission Check**: Strictly blocks any `DELETE` request with `403 Forbidden` ("allow accept permission kecuali delete data").
    -   **Request Logging**: Fancy console output with latency and status codes.

3.  **Fancy Logger**:
    -   Uses `zerolog` with custom console writer for emojis and colorful logs.
    
4.  **Customizable Banner**:
    -   Loads ASCII banner from `banner.txt` on startup.
    -   Default: Rabbit with Rocket.
    -   Configurable via `app.banner_path`.

## Verification Results

### Console Output
The application starts with a custom rabbit banner:
```
 (\_/)
 ( •_•)
 / >  Go Echo App Ready!

13:17:00 INFO Starting Application name="My Fancy Go App"
13:17:00 INFO Initializing Middleware...
...
```

### Endpoints Tested

1.  **GET /api/v1/users**
    -   Result: `200 OK`
    -   Response: `{"message":"Hello from Service A - Users"}`

2.  **GET /api/v1/products**
    -   Result: `404 Not Found` (Correct, as Service B is disabled)

3.  **DELETE /api/v1/users/1**
    -   Result: `403 Forbidden`
    -   Response: `{"error":"Permission Denied: DELETE actions are restricted."}`

## How to Run

1.  **Start the server**:
    ```bash
    go run cmd/app/main.go
    ```

2.  **Configure**:
    -   Edit `config.yaml` to toggle `enable_service_a` or `enable_service_b`.
    -   Edit `banner.txt` to change the startup artwork.
