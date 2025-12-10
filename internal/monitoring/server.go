package monitoring

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"test-go/config"
	"test-go/pkg/infrastructure"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type StatusProvider interface {
	GetStatus() map[string]interface{}
}

type ServiceInfo struct {
	Name       string `json:"name"`
	StructName string `json:"struct_name"`
	Active     bool   `json:"active"`
	Endpoint   string `json:"endpoint"`
}

func Start(
	cfg config.MonitoringConfig,
	appConfig *config.Config,
	statusProvider StatusProvider,
	broadcaster *LogBroadcaster,
	redis *infrastructure.RedisManager,
	postgres *infrastructure.PostgresManager,
	kafka *infrastructure.KafkaManager,
	cron *infrastructure.CronManager,
	services []ServiceInfo,
) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // Enable CORS for development convenience

	// Basic Auth Middleware
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(password), []byte(cfg.Password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	// Serve Static Files
	e.Static("/", "web/monitoring")
	// If index.html is not served by default on /, explicitly handle it or rely on e.Static quirks
	e.GET("/", func(c echo.Context) error {
		return c.File("web/monitoring/index.html")
	})

	// Register Handlers
	h := &Handler{
		config:         appConfig,
		statusProvider: statusProvider,
		broadcaster:    broadcaster,
		redis:          redis,
		postgres:       postgres,
		kafka:          kafka,
		cron:           cron,
		services:       services,
	}
	h.RegisterRoutes(e)

	fmt.Printf("📊 Monitoring UI running on http://localhost:%s\n", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to start monitoring server: %v\n", err)
	}
}
