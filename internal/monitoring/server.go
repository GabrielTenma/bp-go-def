package monitoring

import (
	"fmt"
	"net/http"
	"test-go/config"
	"test-go/internal/monitoring/database"
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
	// Initialize database
	if err := database.InitDB(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to initialize user settings database: %v\n", err)
	} else {
		fmt.Println("✅ User settings database initialized")

		// Ensure upload directory exists
		uploadDir := cfg.UploadDir
		if uploadDir == "" {
			uploadDir = "web/monitoring/uploads"
		}
		if err := database.EnsureUploadDirectory(uploadDir); err != nil {
			fmt.Printf("⚠️  Warning: Failed to create upload directory: %v\n", err)
		}

		// Create default user if not exists
		settings, _ := database.GetUserSettings()
		if settings == nil {
			if err := database.CreateDefaultUser(cfg.Password); err != nil {
				fmt.Printf("⚠️  Warning: Failed to create default user: %v\n", err)
			} else {
				fmt.Println("✅ Default user created")
			}
		}
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // Enable CORS for development convenience

	// Serve login page (no auth required)
	e.GET("/", func(c echo.Context) error {
		return c.File("web/monitoring/login.html")
	})

	// Logout endpoint - returns 401 to clear cached credentials
	e.GET("/logout", func(c echo.Context) error {
		c.Response().Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		return c.String(http.StatusUnauthorized, "Logged out")
	})

	// Static assets (no auth required)
	e.Static("/assets", "web/monitoring/assets")

	// Protected routes group
	protected := e.Group("")
	protected.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Check database password
		err := database.VerifyPassword(password)
		if err == nil {
			return true, nil
		}
		return false, nil
	}))

	// Dashboard and API routes (protected)
	protected.GET("/dashboard", func(c echo.Context) error {
		return c.File("web/monitoring/index.html")
	})
	protected.Static("/api/user/photos", appConfig.Monitoring.UploadDir+"/profiles")

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
	h.RegisterRoutes(protected)

	fmt.Printf("📊 Monitoring UI running on http://localhost:%s\n", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to start monitoring server: %v\n", err)
	}
}
