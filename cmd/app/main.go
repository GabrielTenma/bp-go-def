package main

import (
	"fmt"
	"os"
	"test-go/config"
	"test-go/internal/server"
	"test-go/pkg/logger"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		// We can't use our fancy logger yet if config fails, but we'll try to init a basic one or just panic
		panic("Failed to load config: " + err.Error())
	}

	// 2. Init Logger
	l := logger.New(cfg.App.Debug)

	// Print Banner
	if cfg.App.BannerPath != "" {
		banner, err := os.ReadFile(cfg.App.BannerPath)
		if err == nil {
			fmt.Println(string(banner))
		} else {
			l.Warn("Failed to load banner", "path", cfg.App.BannerPath, "error", err)
		}
	}

	l.Info("Starting Application", "name", cfg.App.Name)

	// 3. Init & Start Server
	srv := server.New(cfg, l)
	if err := srv.Start(); err != nil {
		l.Fatal("Server shutdown abruptly", err)
	}
}
