package server

import (
	"os"
	"test-go/config"
	"test-go/internal/middleware"
	"test-go/internal/monitoring"
	"test-go/internal/services"
	"test-go/internal/services/modules"
	"test-go/pkg/infrastructure"
	"test-go/pkg/logger"
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
		e.DefaultHTTPErrorHandler(err, c)
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
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Restart Endpoint (Maintenance)
	s.echo.POST("/restart", func(c echo.Context) error {
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(1)
		}()
		return c.JSON(200, map[string]string{"status": "restarting", "message": "Service is restarting..."})
	})

	// Add Services here
	registry.Register(modules.NewServiceA(s.config.Services.EnableServiceA))
	registry.Register(modules.NewServiceB(s.config.Services.EnableServiceB))
	registry.Register(modules.NewServiceC(s.config.Services.EnableServiceC))
	registry.Register(modules.NewServiceD(s.postgresManager, s.config.Services.EnableServiceD))

	registry.Boot(s.echo)

	// 4. Start Monitoring (if enabled)
	if s.config.Monitoring.Enabled {
		servicesList := []monitoring.ServiceInfo{
			{Name: "Service A", StructName: "modules.ServiceA", Active: s.config.Services.EnableServiceA, Endpoint: "/api/v1/service-a"},
			{Name: "Service B", StructName: "modules.ServiceB", Active: s.config.Services.EnableServiceB, Endpoint: "/api/v1/service-b"},
			{Name: "Service C", StructName: "modules.ServiceC", Active: s.config.Services.EnableServiceC, Endpoint: "/api/v1/service-c"},
			{Name: "Service D", StructName: "modules.ServiceD", Active: s.config.Services.EnableServiceD, Endpoint: "/tasks"},
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
		"version": "1.0.0",
		"services": map[string]bool{
			"service_a": s.config.Services.EnableServiceA,
			"service_b": s.config.Services.EnableServiceB,
			"service_c": s.config.Services.EnableServiceC,
			"service_d": s.config.Services.EnableServiceD,
		},
		"infrastructure": infra,
		"system": map[string]interface{}{
			"disk":    diskStats,
			"network": netStats,
		},
	}
}
