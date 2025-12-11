package modules

import (
	"fmt"
	"net/http"
	"strconv"
	"test-go/pkg/infrastructure"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type ServiceD struct {
	db      *infrastructure.PostgresManager
	enabled bool
}

func NewServiceD(db *infrastructure.PostgresManager, enabled bool) *ServiceD {
	if enabled && db != nil && db.ORM != nil {
		// Auto-migrate the schema
		if err := db.ORM.AutoMigrate(&Task{}); err != nil {
			fmt.Printf("Error migrating Task model: %v\n", err)
		}
	}
	return &ServiceD{
		db:      db,
		enabled: enabled,
	}
}

func (s *ServiceD) Name() string { return "Service D (Tasks - GORM)" }

func (s *ServiceD) Enabled() bool {
	// Service is enabled only if configured AND DB is available
	return s.enabled && s.db != nil && s.db.ORM != nil
}

func (s *ServiceD) RegisterRoutes(g *echo.Group) {
	sub := g.Group("/tasks")
	sub.GET("", s.listTasks)
	sub.POST("", s.createTask)
	sub.PUT("/:id", s.updateTask)
	sub.DELETE("/:id", s.deleteTask)
}

func (s *ServiceD) listTasks(c echo.Context) error {
	var tasks []Task
	result := s.db.ORM.Find(&tasks)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, tasks)
}

func (s *ServiceD) createTask(c echo.Context) error {
	task := new(Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if result := s.db.ORM.Create(task); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	return c.JSON(http.StatusCreated, task)
}

func (s *ServiceD) updateTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task
	if result := s.db.ORM.First(&task, id); result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	s.db.ORM.Save(&task)
	return c.JSON(http.StatusOK, task)
}

func (s *ServiceD) deleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if result := s.db.ORM.Delete(&Task{}, id); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Task deleted"})
}
