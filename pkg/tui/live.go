package tui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LiveConfig contains configuration for the live TUI
type LiveConfig struct {
	AppName    string
	AppVersion string
	Banner     string
	Port       string
	Env        string
}

// LogEntry represents a log entry
type LogEntry struct {
	Time    time.Time
	Level   string
	Message string
}

// LiveModel is the Bubble Tea model for the live running dashboard
type LiveModel struct {
	spinner   spinner.Model
	config    LiveConfig
	logs      []LogEntry
	logsMutex sync.RWMutex
	startTime time.Time
	width     int
	height    int
	frame     int
	quitting  bool
	maxLogs   int
	program   *tea.Program
}

// Live TUI styles
var (
	liveBannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#BD93F9"))

	liveTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF79C6"))

	liveInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	liveStatusStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#50FA7B"))

	liveDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#44475A"))

	liveLogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6272A4")).
			Padding(0, 1)

	// Single cyan color for progress bar
	liveProgressColor = "#8BE9FD"
)

// Looping progress bar frames
var loopingProgressFrames = []string{
	"█░░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
	"██░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
	"███░░░░░░░░░░░░░░░░░░░░░░░░░░░",
	"████░░░░░░░░░░░░░░░░░░░░░░░░░░",
	"█████░░░░░░░░░░░░░░░░░░░░░░░░░",
	"██████░░░░░░░░░░░░░░░░░░░░░░░░",
	"███████░░░░░░░░░░░░░░░░░░░░░░░",
	"████████░░░░░░░░░░░░░░░░░░░░░░",
	"█████████░░░░░░░░░░░░░░░░░░░░░",
	"██████████░░░░░░░░░░░░░░░░░░░░",
	"░██████████░░░░░░░░░░░░░░░░░░░",
	"░░██████████░░░░░░░░░░░░░░░░░░",
	"░░░██████████░░░░░░░░░░░░░░░░░",
	"░░░░██████████░░░░░░░░░░░░░░░░",
	"░░░░░██████████░░░░░░░░░░░░░░░",
	"░░░░░░██████████░░░░░░░░░░░░░░",
	"░░░░░░░██████████░░░░░░░░░░░░░",
	"░░░░░░░░██████████░░░░░░░░░░░░",
	"░░░░░░░░░██████████░░░░░░░░░░░",
	"░░░░░░░░░░██████████░░░░░░░░░░",
	"░░░░░░░░░░░██████████░░░░░░░░░",
	"░░░░░░░░░░░░██████████░░░░░░░░",
	"░░░░░░░░░░░░░██████████░░░░░░░",
	"░░░░░░░░░░░░░░██████████░░░░░░",
	"░░░░░░░░░░░░░░░██████████░░░░░",
	"░░░░░░░░░░░░░░░░██████████░░░░",
	"░░░░░░░░░░░░░░░░░██████████░░░",
	"░░░░░░░░░░░░░░░░░░██████████░░",
	"░░░░░░░░░░░░░░░░░░░██████████░",
	"░░░░░░░░░░░░░░░░░░░░██████████",
	"░░░░░░░░░░░░░░░░░░░░░█████████",
	"░░░░░░░░░░░░░░░░░░░░░░████████",
	"░░░░░░░░░░░░░░░░░░░░░░░███████",
	"░░░░░░░░░░░░░░░░░░░░░░░░██████",
	"░░░░░░░░░░░░░░░░░░░░░░░░░█████",
	"░░░░░░░░░░░░░░░░░░░░░░░░░░████",
	"░░░░░░░░░░░░░░░░░░░░░░░░░░░███",
	"░░░░░░░░░░░░░░░░░░░░░░░░░░░░██",
	"░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█",
	"░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
}

// NewLiveModel creates a new live TUI model
func NewLiveModel(cfg LiveConfig) *LiveModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))

	return &LiveModel{
		spinner:   s,
		config:    cfg,
		logs:      make([]LogEntry, 0),
		startTime: time.Now(),
		width:     80,
		height:    24,
		maxLogs:   15,
	}
}

type liveTickMsg time.Time
type logMsg LogEntry

func liveTickCmd() tea.Cmd {
	return tea.Every(time.Millisecond*100, func(t time.Time) tea.Msg {
		return liveTickMsg(t)
	})
}

func (m *LiveModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		liveTickCmd(),
	)
}

func (m *LiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case liveTickMsg:
		m.frame = (m.frame + 1) % len(loopingProgressFrames)
		return m, tea.Batch(m.spinner.Tick, liveTickCmd())

	case logMsg:
		m.logsMutex.Lock()
		m.logs = append(m.logs, LogEntry(msg))
		// Keep only the last maxLogs entries
		if len(m.logs) > m.maxLogs {
			m.logs = m.logs[len(m.logs)-m.maxLogs:]
		}
		m.logsMutex.Unlock()
		return m, nil
	}

	return m, nil
}

func (m *LiveModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Sticky Banner
	if m.config.Banner != "" {
		b.WriteString(liveBannerStyle.Render(m.config.Banner))
		b.WriteString("\n")
	}

	// Header with app info (cyan accent)
	cyanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(liveProgressColor)).Bold(true)

	header := fmt.Sprintf("%s %s v%s %s",
		cyanStyle.Render("⚡"),
		liveTitleStyle.Render(m.config.AppName),
		m.config.AppVersion,
		cyanStyle.Render("⚡"),
	)
	b.WriteString(header)
	b.WriteString("\n")

	// Status line
	uptime := time.Since(m.startTime).Round(time.Second)
	statusLine := fmt.Sprintf("  %s %s  │  Port: %s  │  Env: %s  │  Uptime: %s",
		m.spinner.View(),
		liveStatusStyle.Render("RUNNING"),
		liveInfoStyle.Render(m.config.Port),
		liveInfoStyle.Render(m.config.Env),
		liveInfoStyle.Render(uptime.String()),
	)
	b.WriteString(statusLine)
	b.WriteString("\n\n")

	// Looping progress bar - single cyan color
	progressBar := loopingProgressFrames[m.frame]
	progressStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(liveProgressColor))
	b.WriteString("  ")
	b.WriteString(progressStyle.Render(progressBar))
	b.WriteString("\n\n")

	// Logs section
	logsContent := m.renderLogs()
	b.WriteString(liveLogBoxStyle.Render(logsContent))
	b.WriteString("\n")

	// Footer
	footer := liveDimStyle.Render(fmt.Sprintf("  Press 'q' to quit  │  %s", time.Now().Format("15:04:05")))
	b.WriteString(footer)

	return b.String()
}

func (m *LiveModel) renderLogs() string {
	var lines []string

	// Calculate available width for logs
	logWidth := m.width - 6 // Account for box borders and padding
	if logWidth < 60 {
		logWidth = 60
	}
	if logWidth > 120 {
		logWidth = 120
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#8BE9FD")).
		Render("◆ Live Logs")
	lines = append(lines, header)
	lines = append(lines, liveDimStyle.Render(strings.Repeat("─", logWidth)))

	m.logsMutex.RLock()
	defer m.logsMutex.RUnlock()

	if len(m.logs) == 0 {
		lines = append(lines, liveDimStyle.Render("  Waiting for logs..."))
	} else {
		for _, log := range m.logs {
			levelStyle := m.getLevelStyle(log.Level)
			timeStr := log.Time.Format("15:04:05")
			levelStr := fmt.Sprintf("[%-5s]", strings.ToUpper(log.Level))

			// Build the line with proper formatting - NO truncation
			line := fmt.Sprintf("  %s %s %s",
				liveDimStyle.Render(timeStr),
				levelStyle.Render(levelStr),
				lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2")).Render(log.Message),
			)
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *LiveModel) getLevelStyle(level string) lipgloss.Style {
	switch strings.ToLower(level) {
	case "debug":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	case "info":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	case "warn", "warning":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	case "error":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	case "fatal":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	}
}

// AddLog adds a log entry to the TUI
func (m *LiveModel) AddLog(level, message string) {
	if m.program != nil {
		m.program.Send(logMsg{
			Time:    time.Now(),
			Level:   level,
			Message: message,
		})
	}
}

// SetProgram sets the tea.Program reference for sending messages
func (m *LiveModel) SetProgram(p *tea.Program) {
	m.program = p
}

// LiveTUI manages the live TUI instance
type LiveTUI struct {
	model   *LiveModel
	program *tea.Program
}

// NewLiveTUI creates a new live TUI instance
func NewLiveTUI(cfg LiveConfig) *LiveTUI {
	model := NewLiveModel(cfg)
	return &LiveTUI{
		model: model,
	}
}

// Start starts the live TUI in a goroutine
func (t *LiveTUI) Start() {
	t.program = tea.NewProgram(t.model, tea.WithAltScreen())
	t.model.SetProgram(t.program)
	go func() {
		t.program.Run()
	}()
}

// Stop stops the live TUI
func (t *LiveTUI) Stop() {
	if t.program != nil {
		t.program.Quit()
	}
}

// AddLog adds a log to the live TUI
func (t *LiveTUI) AddLog(level, message string) {
	t.model.AddLog(level, message)
}

// Write implements io.Writer for use as a log broadcaster
func (t *LiveTUI) Write(p []byte) (n int, err error) {
	// Parse the log line and add it
	line := strings.TrimSpace(string(p))
	if line != "" {
		level, message := parseLogLine(line)
		if message != "" {
			t.AddLog(level, message)
		}
	}
	return len(p), nil
}

// parseLogLine extracts the level and clean message from a zerolog console output line
// Example input: "15:00:51 INF Scheduled Cron Job job=health_check schedule="*/10 * * * * *""
// Returns: level="info", message="Scheduled Cron Job job=health_check schedule="*/10 * * * * *""
func parseLogLine(line string) (level, message string) {
	level = "info" // default

	// Split by space to find components
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 2 {
		return level, line
	}

	// Check if first part is a timestamp (HH:MM:SS format)
	if len(parts[0]) == 8 && strings.Count(parts[0], ":") == 2 {
		// Second part should be the level abbreviation
		switch strings.ToUpper(parts[1]) {
		case "DBG", "DEBUG":
			level = "debug"
		case "INF", "INFO":
			level = "info"
		case "WRN", "WARN", "WARNING":
			level = "warn"
		case "ERR", "ERROR":
			level = "error"
		case "FTL", "FATAL":
			level = "fatal"
		case "PNC", "PANIC":
			level = "fatal"
		}

		// Message is everything after timestamp and level
		if len(parts) >= 3 {
			message = parts[2]
		} else {
			message = ""
		}
	} else {
		// No timestamp, try to detect level from content
		upperLine := strings.ToUpper(line)
		if strings.Contains(upperLine, "DEBUG") || strings.Contains(upperLine, "DBG") {
			level = "debug"
		} else if strings.Contains(upperLine, "WARN") || strings.Contains(upperLine, "WRN") {
			level = "warn"
		} else if strings.Contains(upperLine, "ERROR") || strings.Contains(upperLine, "ERR") {
			level = "error"
		} else if strings.Contains(upperLine, "FATAL") || strings.Contains(upperLine, "FTL") {
			level = "fatal"
		}
		message = line
	}

	return level, message
}

// RunLiveTUI runs the live TUI and blocks until quit
func RunLiveTUI(cfg LiveConfig) error {
	model := NewLiveModel(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())
	model.SetProgram(p)
	_, err := p.Run()
	return err
}
