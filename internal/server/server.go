package server

import (
	"test-go/config"
	"test-go/internal/middleware"
	"test-go/internal/monitoring"
	"test-go/internal/services"
	"test-go/internal/services/modules"
	"test-go/pkg/infrastructure"
	"test-go/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo            *echo.Echo
	config          *config.Config
	logger          *logger.Logger
	redisManager    *infrastructure.RedisManager
	kafkaManager    *infrastructure.KafkaManager
	postgresManager *infrastructure.PostgresManager
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

	// 2. Init Middleware
	s.logger.Info("Initializing Middleware...")
	middleware.InitMiddlewares(s.echo, middleware.Config{
		AuthType: s.config.Auth.Type,
		Logger:   s.logger,
	})

	// 3. Init Services
	s.logger.Info("Booting Services...")
	registry := services.NewRegistry(s.logger)

	// Add Services here
	registry.Register(modules.NewServiceA(s.config.Services.EnableServiceA))
	registry.Register(modules.NewServiceB(s.config.Services.EnableServiceB))
	registry.Register(modules.NewServiceC(s.config.Services.EnableServiceC))
	registry.Register(modules.NewServiceD(s.postgresManager, s.config.Services.EnableServiceD))

	registry.Boot(s.echo)

	// 4. Start Monitoring (if enabled)
	if s.config.Monitoring.Enabled {
		go monitoring.Start(s.config.Monitoring, s.config, s, s.broadcaster)
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
	return map[string]interface{}{
		"version": "1.0.0",
		"services": map[string]bool{
			"service_a": s.config.Services.EnableServiceA,
			"service_b": s.config.Services.EnableServiceB,
			"service_c": s.config.Services.EnableServiceC,
			"service_d": s.config.Services.EnableServiceD,
		},
		"infrastructure": map[string]bool{
			"redis":    s.config.Redis.Enabled && s.redisManager != nil,
			"kafka":    s.config.Kafka.Enabled && s.kafkaManager != nil,
			"postgres": s.config.Postgres.Enabled && s.postgresManager != nil,
		},
	}
}
