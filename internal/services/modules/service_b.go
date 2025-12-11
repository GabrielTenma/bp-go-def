package modules

import (
	"github.com/labstack/echo/v4"
)

type ServiceB struct {
	enabled bool
}

func NewServiceB(enabled bool) *ServiceB {
	return &ServiceB{enabled: enabled}
}

func (s *ServiceB) Name() string  { return "Service B (Products)" }
func (s *ServiceB) Enabled() bool { return s.enabled }

func (s *ServiceB) RegisterRoutes(g *echo.Group) {
	sub := g.Group("/products")
	sub.GET("", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Hello from Service B - Products"})
	})
}
