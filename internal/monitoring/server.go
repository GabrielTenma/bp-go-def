package monitoring

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"test-go/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type StatusProvider interface {
	GetStatus() map[string]interface{}
}

func Start(cfg config.MonitoringConfig, appConfig *config.Config, statusProvider StatusProvider, broadcaster *LogBroadcaster) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())

	// Basic Auth Middleware (Password only, effectively)
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// We ignore username, only check password
		if subtle.ConstantTimeCompare([]byte(password), []byte(cfg.Password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	// Register Handlers
	h := &Handler{
		config:         appConfig,
		statusProvider: statusProvider,
		broadcaster:    broadcaster,
	}
	h.RegisterRoutes(e)

	fmt.Printf("📊 Monitoring UI running on http://localhost:%s\n", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to start monitoring server: %v\n", err)
	}
}
