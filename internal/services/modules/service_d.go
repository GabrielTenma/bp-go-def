package modules

import (
	"fmt"
	"strconv"
	"test-go/pkg/infrastructure"
	"test-go/pkg/response"

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

func (s *ServiceD) Endpoints() []string { return []string{"/tasks"} }

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
		return response.InternalServerError(c, result.Error.Error())
	}
	return response.Success(c, tasks)
}

func (s *ServiceD) createTask(c echo.Context) error {
	task := new(Task)
	if err := c.Bind(task); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	if result := s.db.ORM.Create(task); result.Error != nil {
		return response.InternalServerError(c, result.Error.Error())
	}

	return response.Created(c, task)
}

func (s *ServiceD) updateTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task
	if result := s.db.ORM.First(&task, id); result.Error != nil {
		return response.NotFound(c, "Task not found")
	}

	if err := c.Bind(&task); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	s.db.ORM.Save(&task)
	return response.Success(c, task)
}

func (s *ServiceD) deleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if result := s.db.ORM.Delete(&Task{}, id); result.Error != nil {
		return response.InternalServerError(c, result.Error.Error())
	}
	return response.Success(c, nil, "Task deleted")
}
