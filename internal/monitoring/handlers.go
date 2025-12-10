package monitoring

import (
	"fmt"
	"net/http"
	"test-go/config"
	"test-go/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	config         *config.Config
	statusProvider StatusProvider
	broadcaster    *LogBroadcaster
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.serveUI)
	e.GET("/api/status", h.getStatus)
	e.GET("/api/config", h.getConfig)
	e.GET("/api/logs", h.streamLogs)
	e.GET("/api/cpu", h.streamCPU)
}

func (h *Handler) serveUI(c echo.Context) error {
	return c.HTML(http.StatusOK, IndexHTML)
}

func (h *Handler) getStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, h.statusProvider.GetStatus())
}

func (h *Handler) getConfig(c echo.Context) error {
	// Return config but redact sensitive info ideally
	// For now we assume the admin is trusted
	return c.JSON(http.StatusOK, h.config)
}

func (h *Handler) streamLogs(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	logs := h.broadcaster.Subscribe()
	defer h.broadcaster.Unsubscribe(logs)

	for {
		select {
		case msg := <-logs:
			fmt.Fprintf(c.Response(), "data: %s\n\n", msg) // Send raw log line
			c.Response().Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}
}

func (h *Handler) streamCPU(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats, _ := utils.GetSystemStats()
			fmt.Fprintf(c.Response(), "data: %.2f\n\n", stats["cpu_percent"])
			c.Response().Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}
}
