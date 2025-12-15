package server

import (
	"os"
	"reflect"
	"test-go/config"
	"test-go/internal/middleware"
	"test-go/internal/monitoring"
	"test-go/internal/services"
	"test-go/internal/services/modules"
	"test-go/pkg/infrastructure"
	"test-go/pkg/logger"
	"test-go/pkg/response"
	"test-go/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo            *echo.Echo
	config          *config.Config
	logger          *logger.Logger
	redisManager    *infrastructure.RedisManager
	kafkaManager    *infrastructure.KafkaManager
	postgresManager *infrastructure.PostgresManager
	cronManager     *infrastructure.CronManager
	broadcaster     *monitoring.LogBroadcaster
}

func New(cfg *config.Config, l *logger.Logger, b *monitoring.LogBroadcaster) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Custom HTTP Error Handler for JSON responses
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		l.Error("HTTP Error", err)

		// Handle HTTP errors with JSON response
		if he, ok := err.(*echo.HTTPError); ok {
			var message string
			code := he.Code

			// Custom message for 404 Not Found
			if code == 404 {
				message = "Endpoint not found. This incident will be reported."
				response.Error(c, code, "ENDPOINT_NOT_FOUND", message, map[string]interface{}{
					"path":   c.Request().URL.Path,
					"method": c.Request().Method,
				})
				return
			}

			// For other HTTP errors, use the original message if it's a string
			if msg, ok := he.Message.(string); ok {
				message = msg
			} else {
				message = "An unexpected error occurred"
			}
			response.Error(c, code, "HTTP_ERROR", message)
			return
		}

		// For non-HTTP errors, return internal server error
		response.InternalServerError(c, "An unexpected error occurred")
	}

	return &Server{
		echo:        e,
		config:      cfg,
		logger:      l,
		broadcaster: b,
	}
}

func (s *Server) Start() error {
	// 1. Init Infrastructure
	s.logger.Info("Initializing Infrastructure...")

	// Redis
	if s.config.Redis.Enabled {
		rdb, err := infrastructure.NewRedisClient(s.config.Redis)
		if err != nil {
			s.logger.Error("Failed to initialize Redis", err)
		} else {
			s.redisManager = rdb
			s.logger.Info("Redis initialized")
		}
	}

	// Kafka
	if s.config.Kafka.Enabled {
		// Note: NewKafkaManager replaces NewKafkaProducer
		km, err := infrastructure.NewKafkaManager(s.config.Kafka)
		if err != nil {
			s.logger.Error("Failed to initialize Kafka", err)
		} else {
			s.kafkaManager = km
			s.logger.Info("Kafka initialized")
		}
	}

	// Postgres
	if s.config.Postgres.Enabled {
		db, err := infrastructure.NewPostgresDB(s.config.Postgres)
		if err != nil {
			s.logger.Error("Failed to initialize Postgres", err)
		} else {
			s.postgresManager = db
			s.logger.Info("Postgres initialized")
		}
	}

	// Cron Jobs
	if s.config.Cron.Enabled {
		s.cronManager = infrastructure.NewCronManager()
		s.logger.Info("Initializing Cron Jobs...")

		for name, schedule := range s.config.Cron.Jobs {
			// Capture variables
			jobName := name

			// Dummy action
			_, err := s.cronManager.AddJob(jobName, schedule, func() {
				s.logger.Info("Executing Cron Job", "job", jobName)
			})
			if err != nil {
				s.logger.Error("Failed to schedule cron job", err, "job", jobName)
			} else {
				s.logger.Info("Scheduled Cron Job", "job", jobName, "schedule", schedule)
			}
		}
		s.cronManager.Start()
	}

	// 2. Init Middleware
	s.logger.Info("Initializing Middleware...")
	middleware.InitMiddlewares(s.echo, middleware.Config{
		AuthType: s.config.Auth.Type,
		Logger:   s.logger,
	})

	// 3. Init Services
	s.logger.Info("Booting Services...")
	registry := services.NewRegistry(s.logger)

	// Health Check Endpoint
	s.echo.GET("/health", func(c echo.Context) error {
		return response.Success(c, map[string]string{"status": "ok"})
	})

	// Restart Endpoint (Maintenance)
	s.echo.POST("/restart", func(c echo.Context) error {
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(1)
		}()
		return response.Success(c, map[string]string{"status": "restarting", "message": "Service is restarting..."})
	})

	// Add Services here - use IsEnabled() for dynamic config lookup
	registry.Register(modules.NewServiceA(s.config.Services.IsEnabled("service_a")))
	registry.Register(modules.NewServiceB(s.config.Services.IsEnabled("service_b")))
	registry.Register(modules.NewServiceC(s.config.Services.IsEnabled("service_c")))
	registry.Register(modules.NewServiceD(s.postgresManager, s.config.Services.IsEnabled("service_d")))

	registry.Boot(s.echo)

	// 4. Start Monitoring (if enabled)
	if s.config.Monitoring.Enabled {
		// Dynamic Service List Generation
		var servicesList []monitoring.ServiceInfo
		for _, srv := range registry.GetServices() {
			// Prepend /api/v1 to endpoints
			var fullEndpoints []string
			for _, endp := range srv.Endpoints() {
				fullEndpoints = append(fullEndpoints, "/api/v1"+endp)
			}

			servicesList = append(servicesList, monitoring.ServiceInfo{
				Name:       srv.Name(),
				StructName: reflect.TypeOf(srv).Elem().String(),
				Active:     srv.Enabled(),
				Endpoints:  fullEndpoints,
			})
		}
		go monitoring.Start(s.config.Monitoring, s.config, s, s.broadcaster, s.redisManager, s.postgresManager, s.kafkaManager, s.cronManager, servicesList)
		s.logger.Info("Monitoring interface started", "port", s.config.Monitoring.Port)
	}

	// 5. Start Server
	port := s.config.Server.Port
	s.logger.Info("Server is ready to handle requests", "port", port, "env", s.config.App.Env)

	// Create a channel to listen for OS signals (graceful shutdown could be added here)

	return s.echo.Start(":" + port)
}

// GetStatus satisfies monitoring.StatusProvider
func (s *Server) GetStatus() map[string]interface{} {
	diskStats, _ := utils.GetDiskUsage()
	netStats, _ := utils.GetNetworkInfo()

	infra := map[string]bool{
		"redis":    s.config.Redis.Enabled && s.redisManager != nil,
		"kafka":    s.config.Kafka.Enabled && s.kafkaManager != nil,
		"postgres": s.config.Postgres.Enabled && s.postgresManager != nil,
		"cron":     s.config.Cron.Enabled && s.cronManager != nil,
	}

	return map[string]interface{}{
		"version":        "1.0.0",
		"services":       s.config.Services, // Dynamic map from config
		"infrastructure": infra,
		"system": map[string]interface{}{
			"disk":    diskStats,
			"network": netStats,
		},
	}
}
