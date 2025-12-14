package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"test-go/config"
	"test-go/internal/monitoring"
	"test-go/internal/server"
	"test-go/pkg/logger"
	"test-go/pkg/tui"
	"time"
)

// clearScreen clears the terminal screen (cross-platform)
func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows: use cmd /c cls
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		// Linux, macOS, and others: use clear command
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	// Clear the terminal screen for a fresh start
	clearScreen()

	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 2. Load Banner
	var bannerText string
	if cfg.App.BannerPath != "" {
		banner, err := os.ReadFile(cfg.App.BannerPath)
		if err == nil {
			bannerText = string(banner)
		}
	}

	// 3. Init Broadcaster for monitoring
	broadcaster := monitoring.NewLogBroadcaster()

	// Check if TUI mode is enabled
	if cfg.App.EnableTUI {
		// ===== TUI MODE =====
		runWithTUI(cfg, bannerText, broadcaster)
	} else {
		// ===== TRADITIONAL CONSOLE MODE =====
		runWithConsole(cfg, bannerText, broadcaster)
	}
}

// runWithTUI runs the application with fancy TUI interface
func runWithTUI(cfg *config.Config, bannerText string, broadcaster *monitoring.LogBroadcaster) {
	tuiConfig := tui.StartupConfig{
		AppName:     cfg.App.Name,
		AppVersion:  "1.0.0",
		Banner:      bannerText,
		Port:        cfg.Server.Port,
		Env:         cfg.App.Env,
		IdleSeconds: cfg.App.StartupDelay,
	}

	// Define boot sequence
	initQueue := []tui.ServiceInit{
		{Name: "Configuration", Enabled: true, InitFunc: nil},
		{Name: "Redis Cache", Enabled: cfg.Redis.Enabled, InitFunc: nil},
		{Name: "Kafka Messaging", Enabled: cfg.Kafka.Enabled, InitFunc: nil},
		{Name: "PostgreSQL", Enabled: cfg.Postgres.Enabled, InitFunc: nil},
		{Name: "Cron Scheduler", Enabled: cfg.Cron.Enabled, InitFunc: nil},
		{Name: "Middleware", Enabled: true, InitFunc: nil},
		{Name: "Service A", Enabled: cfg.Services.EnableServiceA, InitFunc: nil},
		{Name: "Service B", Enabled: cfg.Services.EnableServiceB, InitFunc: nil},
		{Name: "Service C", Enabled: cfg.Services.EnableServiceC, InitFunc: nil},
		{Name: "Service D", Enabled: cfg.Services.EnableServiceD, InitFunc: nil},
		{Name: "Monitoring", Enabled: cfg.Monitoring.Enabled, InitFunc: nil},
	}

	// Run the boot sequence TUI
	_, _ = tui.RunBootSequence(tuiConfig, initQueue)

	// Create Live TUI for continuous display
	liveTUI := tui.NewLiveTUI(tui.LiveConfig{
		AppName:    cfg.App.Name,
		AppVersion: "1.0.0",
		Banner:     bannerText,
		Port:       cfg.Server.Port,
		Env:        cfg.App.Env,
	})

	// Init Logger (quiet mode so logs go to TUI only)
	// We also broadcast to the monitoring system so the Web UI Live Logs work
	multiWriter := io.MultiWriter(liveTUI, broadcaster)
	l := logger.NewQuiet(cfg.App.Debug, multiWriter)

	// Start Live TUI in background
	liveTUI.Start()

	// Give TUI a moment to initialize
	time.Sleep(100 * time.Millisecond)

	// Add initial logs
	liveTUI.AddLog("info", "Server starting on port "+cfg.Server.Port)
	liveTUI.AddLog("info", "Environment: "+cfg.App.Env)

	// Start Server in background
	srv := server.New(cfg, l, broadcaster)
	go func() {
		liveTUI.AddLog("info", "HTTP server listening...")
		if err := srv.Start(); err != nil {
			liveTUI.AddLog("fatal", "Server error: "+err.Error())
		}
	}()

	// Give server a moment to start
	time.Sleep(500 * time.Millisecond)
	liveTUI.AddLog("info", "Server ready at http://localhost:"+cfg.Server.Port)
	if cfg.Monitoring.Enabled {
		liveTUI.AddLog("info", "Monitoring at http://localhost:"+cfg.Monitoring.Port)
	}

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal
	<-sigChan

	liveTUI.AddLog("warn", "Shutting down...")
	liveTUI.Stop()
}

// runWithConsole runs the application with traditional console logging
func runWithConsole(cfg *config.Config, bannerText string, broadcaster *monitoring.LogBroadcaster) {
	// Print banner to console
	if bannerText != "" {
		fmt.Print("\033[35m") // Purple color
		fmt.Println(bannerText)
		fmt.Print("\033[0m") // Reset color
	}

	// Init Logger (normal mode with console output)
	l := logger.New(cfg.App.Debug, broadcaster)

	// Log startup info
	l.Info("Starting Application", "name", cfg.App.Name, "env", cfg.App.Env)
	l.Info("TUI mode disabled, using traditional console logging")

	// Log enabled services
	l.Info("Initializing services...")
	logServiceStatus(l, "Redis Cache", cfg.Redis.Enabled)
	logServiceStatus(l, "Kafka Messaging", cfg.Kafka.Enabled)
	logServiceStatus(l, "PostgreSQL", cfg.Postgres.Enabled)
	logServiceStatus(l, "Cron Scheduler", cfg.Cron.Enabled)
	logServiceStatus(l, "Service A", cfg.Services.EnableServiceA)
	logServiceStatus(l, "Service B", cfg.Services.EnableServiceB)
	logServiceStatus(l, "Service C", cfg.Services.EnableServiceC)
	logServiceStatus(l, "Service D", cfg.Services.EnableServiceD)
	logServiceStatus(l, "Monitoring", cfg.Monitoring.Enabled)

	// Start Server
	srv := server.New(cfg, l, broadcaster)
	go func() {
		l.Info("HTTP server listening", "port", cfg.Server.Port)
		if err := srv.Start(); err != nil {
			l.Fatal("Server error", err)
		}
	}()

	// Give server a moment to start
	time.Sleep(500 * time.Millisecond)
	l.Info("Server ready", "url", "http://localhost:"+cfg.Server.Port)
	if cfg.Monitoring.Enabled {
		l.Info("Monitoring dashboard", "url", "http://localhost:"+cfg.Monitoring.Port)
	}

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal
	<-sigChan

	l.Warn("Shutting down...")
}

// logServiceStatus logs whether a service is enabled or skipped
func logServiceStatus(l *logger.Logger, name string, enabled bool) {
	if enabled {
		l.Info("Service initialized", "service", name, "status", "enabled")
	} else {
		l.Debug("Service skipped", "service", name, "status", "disabled")
	}
}
