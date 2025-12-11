package monitoring

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"test-go/config"
	"test-go/pkg/infrastructure"
	"test-go/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
)

const bannerFilePath = "data/banner.txt" // Or wherever you want to store the banner file

type Handler struct {
	config         *config.Config
	statusProvider StatusProvider
	broadcaster    *LogBroadcaster
	redis          *infrastructure.RedisManager
	postgres       *infrastructure.PostgresManager
	kafka          *infrastructure.KafkaManager
	cron           *infrastructure.CronManager
	minio          *infrastructure.MinIOManager
	system         *infrastructure.SystemManager
	http           *infrastructure.HttpManager
	services       []ServiceInfo

	// Dummy Logs
	dummyMu     sync.Mutex
	dummyActive bool
	dummyStop   chan struct{}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	// e.GET("/", h.serveUI) // handled by static now
	g.GET("/api/status", h.getStatus)
	g.GET("/api/monitoring/config", h.getMonitoringConfig) // New
	g.GET("/api/config", h.getConfig)
	g.GET("/api/config/raw", h.getRawConfig)     // New
	g.POST("/api/config", h.saveConfig)          // New
	g.POST("/api/config/backup", h.backupConfig) // New
	g.GET("/api/logs", h.streamLogs)
	g.GET("/api/cpu", h.streamCPU)
	g.GET("/api/endpoints", h.getEndpoints)
	g.GET("/api/cron", h.getCronJobs)

	// Utils
	g.GET("/api/logs/dummy/status", h.getDummyStatus)
	g.POST("/api/logs/dummy", h.toggleDummyLogs)

	// Banner
	g.GET("/api/banner", h.getBanner)
	g.POST("/api/banner", h.saveBanner)

	// User Settings
	g.GET("/api/user/settings", h.getUserSettings)
	g.POST("/api/user/settings", h.updateUserSettings)
	g.POST("/api/user/password", h.changePassword)
	g.POST("/api/user/photo", h.uploadPhoto)
	g.DELETE("/api/user/photo", h.deleteUserPhoto)
	// Note: Static route for photos is registered in server.go

	// New Endpoints
	g.GET("/api/redis/keys", h.getRedisKeys)
	g.GET("/api/redis/key/:key", h.getRedisValue)
	g.GET("/api/postgres/queries", h.getPostgresQueries)
	g.GET("/api/postgres/info", h.getPostgresInfo)
	g.GET("/api/kafka/topics", h.getKafkaTopics)
	g.POST("/api/logs/dummy", h.toggleDummyLogs)
}

func (h *Handler) getDummyStatus(c echo.Context) error {
	h.dummyMu.Lock()
	defer h.dummyMu.Unlock()
	return c.JSON(http.StatusOK, map[string]bool{"active": h.dummyActive})
}

func (h *Handler) toggleDummyLogs(c echo.Context) error {
	h.dummyMu.Lock()
	defer h.dummyMu.Unlock()

	type Req struct {
		Enable bool `json:"enable"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Enable {
		if h.dummyActive {
			return c.JSON(http.StatusOK, map[string]string{"message": "Already active"})
		}
		h.dummyActive = true
		h.dummyStop = make(chan struct{})
		go h.runDummyLogs(h.dummyStop)
		return c.JSON(http.StatusOK, map[string]string{"message": "Dummy logs enabled"})
	} else {
		if !h.dummyActive {
			return c.JSON(http.StatusOK, map[string]string{"message": "Already inactive"})
		}
		h.dummyActive = false
		close(h.dummyStop)
		return c.JSON(http.StatusOK, map[string]string{"message": "Dummy logs disabled"})
	}
}

func (h *Handler) getMonitoringConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"title":    h.config.Monitoring.Title,
		"subtitle": h.config.Monitoring.Subtitle,
	})
}

func (h *Handler) runDummyLogs(stop chan struct{}) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	messages := []string{
		"User login successful",
		"Cache miss for key user:123",
		"Database query took 150ms",
		"Background job processing started",
		"Kafka message consumed from topic: orders",
		"Payment gateway timeout",
		"Service health check passed",
		"Redis connection pool refreshing",
	}

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			level := levels[time.Now().UnixNano()%int64(len(levels))]
			msg := messages[time.Now().UnixNano()%int64(len(messages))]

			// Format as zerolog JSON-like output (or whatever format frontend expects)
			// Frontend expects raw text.
			// But broadcaster Write method expects []byte.
			// We can format it nicely.

			timestamp := time.Now().Format(time.RFC3339)
			logLine := fmt.Sprintf(`{"time":"%s","level":"%s","message":"[DUMMY] %s"}`+"\n", timestamp, level, msg)

			// If frontend expects raw string from `data:`, and `h.streamLogs` writes `msg` directly...
			// The broadcaster receives `[]byte` and sends it to channel.
			// `streamLogs` reads `msg` and writes `fmt.Fprintf(c.Response(), "data: %s\n\n", msg)`
			// So `msg` should be the full string line.

			h.broadcaster.Write([]byte(logLine))
		}
	}
}

func (h *Handler) getStatus(c echo.Context) error {
	// Collect status from all sources
	status := h.statusProvider.GetStatus()
	status["redis"] = h.redis.GetStatus()
	status["postgres"] = h.postgres.GetStatus()
	status["kafka"] = h.kafka.GetStatus()
	status["cron"] = h.cron.GetStatus()

	// New Infrastructure
	status["storage"] = h.minio.GetStatus()
	status["system"] = h.system.GetStats()
	status["system_info"] = h.system.GetHostInfo()
	status["external"] = h.http.GetStatus()

	status["services"] = h.services
	return c.JSON(http.StatusOK, status)
}

func (h *Handler) getConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, h.config)
}

func (h *Handler) getEndpoints(c echo.Context) error {
	return c.JSON(http.StatusOK, h.services)
}

// ... existing streamLogs and streamCPU ...

func (h *Handler) getRedisKeys(c echo.Context) error {
	if h.redis == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Redis not enabled"})
	}
	pattern := c.QueryParam("pattern")
	if pattern == "" {
		pattern = "*"
	}
	keys, err := h.redis.ScanKeys(c.Request().Context(), pattern)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, keys)
}

func (h *Handler) getRedisValue(c echo.Context) error {
	if h.redis == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Redis not enabled"})
	}
	key := c.Param("key")
	val, err := h.redis.GetValue(c.Request().Context(), key)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"key": key, "value": val})
}

func (h *Handler) getPostgresQueries(c echo.Context) error {
	if h.postgres == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Postgres not enabled"})
	}
	queries, err := h.postgres.GetRunningQueries(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, queries)
}

func (h *Handler) getPostgresInfo(c echo.Context) error {
	if h.postgres == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Postgres not enabled"})
	}
	info, err := h.postgres.GetDBInfo(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	count, _ := h.postgres.GetSessionCount(c.Request().Context())
	info["sessions"] = count

	return c.JSON(http.StatusOK, info)
}

func (h *Handler) getKafkaTopics(c echo.Context) error {
	// Placeholder: To implement true Kafka monitoring, we need Admin client in KafkaManager.
	// For now return dummy or basic status.
	if h.kafka == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Kafka not enabled"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Kafka monitoring requires Admin API (not implemented yet)"})
}

func (h *Handler) getCronJobs(c echo.Context) error {
	if h.cron == nil {
		return c.JSON(http.StatusOK, []interface{}{}) // Return empty if disabled
	}
	return c.JSON(http.StatusOK, h.cron.GetJobs())
}

func (h *Handler) getBanner(c echo.Context) error {
	path := h.config.App.BannerPath
	if path == "" {
		path = "banner.txt"
	}

	content, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, return empty string or error?
		// User might want to create it.
		if os.IsNotExist(err) {
			return c.JSON(http.StatusOK, map[string]string{"content": ""})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"content": string(content)})
}

func (h *Handler) saveBanner(c echo.Context) error {
	type Req struct {
		Content string `json:"content"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	path := h.config.App.BannerPath
	if path == "" {
		path = "banner.txt"
	}

	// Write file (create if local, 0644)
	err := os.WriteFile(path, []byte(req.Content), 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save banner: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Banner saved successfully"})
}

// Config Handlers
func (h *Handler) getRawConfig(c echo.Context) error {
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"content": string(content)})
}

func (h *Handler) saveConfig(c echo.Context) error {
	type Req struct {
		Content string `json:"content"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err := os.WriteFile("config.yaml", []byte(req.Content), 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save config: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Config saved successfully. Restart required to apply changes."})
}

func (h *Handler) backupConfig(c echo.Context) error {
	input, err := os.ReadFile("config.yaml")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read config: " + err.Error()})
	}

	backupName := fmt.Sprintf("config.yaml.bak.%d", time.Now().Unix())
	err = os.WriteFile(backupName, input, 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create backup: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Backup created: " + backupName})
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
			fmt.Fprintf(c.Response(), "data: %s\n\n", msg)
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
