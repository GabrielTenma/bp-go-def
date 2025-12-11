package monitoring

import (
	"net/http"
	"strings"
	"test-go/internal/monitoring/database"
	"test-go/internal/monitoring/session"
	"time"

	"github.com/labstack/echo/v4"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// handleLogin handles user login
func handleLogin(sessionManager *session.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, LoginResponse{
				Success: false,
				Message: "Invalid request",
			})
		}

		// Get user settings from database
		settings, err := database.GetUserSettings()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, LoginResponse{
				Success: false,
				Message: "Internal server error",
			})
		}

		if settings == nil {
			return c.JSON(http.StatusUnauthorized, LoginResponse{
				Success: false,
				Message: "Invalid username or password",
			})
		}

		// Validate username matches database (case-insensitive)
		if !strings.EqualFold(req.Username, settings.Username) {
			return c.JSON(http.StatusUnauthorized, LoginResponse{
				Success: false,
				Message: "Invalid username or password",
			})
		}

		// Validate password against database
		err = database.VerifyPassword(req.Password)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, LoginResponse{
				Success: false,
				Message: "Invalid username or password",
			})
		}

		// Create session using the actual username from database
		sess, err := sessionManager.Create(settings.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, LoginResponse{
				Success: false,
				Message: "Failed to create session",
			})
		}

		// Set session cookie (24 hours)
		session.SetCookie(c, sess.ID, int(24*time.Hour.Seconds()))

		return c.JSON(http.StatusOK, LoginResponse{
			Success: true,
			Message: "Login successful",
		})
	}
}

// handleLogout handles user logout
func handleLogout(sessionManager *session.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get session cookie
		cookie, err := c.Cookie(session.SessionCookieName)
		if err == nil {
			// Delete session from manager
			sessionManager.Delete(cookie.Value)
		}

		// Clear cookie
		session.ClearCookie(c)

		// Prevent caching of logout response
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Logged out successfully",
		})
	}
}
