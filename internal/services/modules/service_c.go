package modules

import (
	"net/http"
	"time"

	"test-go/pkg/cache"

	"github.com/labstack/echo/v4"
)

type ServiceC struct {
	enabled bool
	store   *cache.Cache[string]
}

func NewServiceC(enabled bool) *ServiceC {
	return &ServiceC{
		enabled: enabled,
		store:   cache.New[string](),
	}
}

func (s *ServiceC) Name() string  { return "Service C (Cache Demo)" }
func (s *ServiceC) Enabled() bool { return s.enabled }

type CacheRequest struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl_seconds"` // Optional
}

func (s *ServiceC) RegisterRoutes(g *echo.Group) {
	sub := g.Group("/cache")

	// GET /cache/:key
	sub.GET("/:key", func(c echo.Context) error {
		key := c.Param("key")
		val, found := s.store.Get(key)
		if !found {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Key not found or expired"})
		}
		return c.JSON(http.StatusOK, map[string]string{"key": key, "value": val})
	})

	// POST /cache/:key
	sub.POST("/:key", func(c echo.Context) error {
		key := c.Param("key")
		var req CacheRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		}

		ttl := time.Duration(req.TTL) * time.Second
		s.store.Set(key, req.Value, ttl)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Cached successfully",
			"key":     key,
			"ttl":     ttl.String(),
		})
	})
}
