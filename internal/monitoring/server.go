package monitoring

import (
	"fmt"
	"net/http"
	"test-go/config"
	"test-go/internal/monitoring/database"
	"test-go/internal/monitoring/session"
	"test-go/pkg/infrastructure"
	"time"

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

	// Initialize Infrastructure Managers
	minioMgr, err := infrastructure.NewMinIOManager(appConfig.Monitoring.MinIO)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to connect to MinIO: %v\n", err)
	} else {
		fmt.Println("✅ MinIO Manager initialized")
	}

	systemMgr := infrastructure.NewSystemManager()
	httpMgr := infrastructure.NewHttpManager(appConfig.Monitoring.External)

	// Initialize session manager
	sessionManager := session.NewManager(24 * time.Hour)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // Enable CORS for development convenience

	// Public routes (no auth required)
	e.GET("/", func(c echo.Context) error {
		return c.File("web/monitoring/login.html")
	})
	e.Static("/assets", "web/monitoring/assets")

	// Auth endpoints
	e.POST("/login", handleLogin(sessionManager))
	e.POST("/logout", handleLogout(sessionManager))

	// Protected routes group (require session)
	protected := e.Group("")
	protected.Use(session.Middleware(sessionManager))

	// Dashboard and API routes (protected)
	protected.GET("/dashboard", func(c echo.Context) error {
		return c.File("web/monitoring/index.html")
	})
	protected.Static("/api/user/photos", appConfig.Monitoring.UploadDir+"/profiles")

	// Register API Handlers
	h := &Handler{
		config:         appConfig,
		statusProvider: statusProvider,
		broadcaster:    broadcaster,
		redis:          redis,
		postgres:       postgres,
		kafka:          kafka,
		cron:           cron,
		services:       services,
		minio:          minioMgr,
		system:         systemMgr,
		http:           httpMgr,
	}
	h.RegisterRoutes(protected)

	fmt.Printf("📊 Monitoring UI running on http://localhost:%s\n", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to start monitoring server: %v\n", err)
	}
}
