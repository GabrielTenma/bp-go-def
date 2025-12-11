package modules

import (
	"github.com/labstack/echo/v4"
)

type ServiceA struct {
	enabled bool
}

func NewServiceA(enabled bool) *ServiceA {
	return &ServiceA{enabled: enabled}
}

func (s *ServiceA) Name() string  { return "Service A (Users)" }
func (s *ServiceA) Enabled() bool { return s.enabled }

func (s *ServiceA) RegisterRoutes(g *echo.Group) {
	sub := g.Group("/users")
	sub.GET("", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Hello from Service A - Users"})
	})
	sub.DELETE("/:id", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Deleting user..."})
	})
}
