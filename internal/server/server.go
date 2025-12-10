package server

import (
	"test-go/config"
	"test-go/internal/middleware"
	"test-go/internal/services"
	"test-go/internal/services/modules"
	"test-go/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	logger *logger.Logger
}

func New(cfg *config.Config, l *logger.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Custom HTTP Error Handler for JSON responses
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		l.Error("HTTP Error", err)
		e.DefaultHTTPErrorHandler(err, c)
	}

	return &Server{
		echo:   e,
		config: cfg,
		logger: l,
	}
}

func (s *Server) Start() error {
	// 1. Init Middleware
	s.logger.Info("Initializing Middleware...")
	middleware.InitMiddlewares(s.echo, middleware.Config{
		AuthType: s.config.Auth.Type,
		Logger:   s.logger,
	})

	// 2. Init Services
	s.logger.Info("Booting Services...")
	registry := services.NewRegistry(s.logger)

	// Add Services here
	registry.Register(modules.NewServiceA(s.config.Services.EnableServiceA))
	registry.Register(modules.NewServiceB(s.config.Services.EnableServiceB))
	registry.Register(modules.NewServiceC(s.config.Services.EnableServiceC))

	registry.Boot(s.echo)

	// 3. Start Server
	port := s.config.Server.Port
	s.logger.Info("Server is ready to handle requests", "port", port, "env", s.config.App.Env)

	// Create a channel to listen for OS signals (graceful shutdown could be added here)

	return s.echo.Start(":" + port)
}
