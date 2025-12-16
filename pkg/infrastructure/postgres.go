package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"test-go/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresManager struct {
	DB  *sql.DB
	ORM *gorm.DB
}

type PostgresConnectionManager struct {
	connections map[string]*PostgresManager
	mu          sync.RWMutex
}

func NewPostgresDB(cfg config.PostgresConfig) (*PostgresManager, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Open raw SQL connection
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Initialize GORM with the existing SQL connection
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %w", err)
	}

	return &PostgresManager{
		DB:  sqlDB,
		ORM: gormDB,
	}, nil
}

func NewPostgresConnectionManager(cfg config.PostgresMultiConfig) (*PostgresConnectionManager, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	manager := &PostgresConnectionManager{
		connections: make(map[string]*PostgresManager),
	}

	for _, connCfg := range cfg.Connections {
		if !connCfg.Enabled {
			continue
		}

		// Convert connection config to single config for backward compatibility
		singleCfg := config.PostgresConfig{
			Enabled:  connCfg.Enabled,
			Host:     connCfg.Host,
			Port:     connCfg.Port,
			User:     connCfg.User,
			Password: connCfg.Password,
			DBName:   connCfg.DBName,
			SSLMode:  connCfg.SSLMode,
		}

		db, err := NewPostgresDB(singleCfg)
		if err != nil {
			// Log error but continue with other connections
			// Don't fail the entire manager initialization
			continue
		}

		if db != nil {
			manager.connections[connCfg.Name] = db
		}
	}

	return manager, nil
}

// GetConnection returns a specific named connection
func (m *PostgresConnectionManager) GetConnection(name string) (*PostgresManager, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, exists := m.connections[name]
	return conn, exists
}

// GetDefaultConnection returns the first connection or nil if none exist
func (m *PostgresConnectionManager) GetDefaultConnection() (*PostgresManager, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, conn := range m.connections {
		return conn, true
	}
	return nil, false
}

// GetAllConnections returns all connections
func (m *PostgresConnectionManager) GetAllConnections() map[string]*PostgresManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Create a copy to avoid race conditions
	copy := make(map[string]*PostgresManager, len(m.connections))
	for k, v := range m.connections {
		copy[k] = v
	}
	return copy
}

// GetStatus returns status for all connections
func (m *PostgresConnectionManager) GetStatus() map[string]map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	status := make(map[string]map[string]interface{})

	for name, conn := range m.connections {
		status[name] = conn.GetStatus()
	}

	return status
}

// CloseAll closes all connections
func (m *PostgresConnectionManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error
	for name, conn := range m.connections {
		if err := conn.DB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection '%s': %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing connections: %v", errors)
	}
	return nil
}

func (p *PostgresManager) GetStatus() map[string]interface{} {
	stats := make(map[string]interface{})
	if p == nil || p.DB == nil {
		stats["connected"] = false
		return stats
	}

	err := p.DB.Ping()
	stats["connected"] = err == nil

	// DB Stats
	dbStats := p.DB.Stats()
	stats["open_connections"] = dbStats.OpenConnections
	stats["in_use"] = dbStats.InUse
	stats["idle"] = dbStats.Idle
	stats["wait_count"] = dbStats.WaitCount
	stats["wait_duration_ms"] = dbStats.WaitDuration.Milliseconds()

	return stats
}

// Query executes a query that returns rows, typically a SELECT.
func (p *PostgresManager) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (p *PostgresManager) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning any rows.
func (p *PostgresManager) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.DB.ExecContext(ctx, query, args...)
}

// Select is a semantic alias for Query.
func (p *PostgresManager) Select(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.Query(ctx, query, args...)
}

// Insert executes an INSERT statement and returns the number of rows affected.
func (p *PostgresManager) Insert(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := p.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// ExecuteRawQuery executes a raw SQL query and returns the results as a slice of maps
func (p *PostgresManager) ExecuteRawQuery(ctx context.Context, query string) ([]map[string]interface{}, error) {
	if p.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Initialize with make to ensure empty slice [] instead of nil
	results := make([]map[string]interface{}, 0)

	for rows.Next() {
		// Create a slice of interface{} to hold values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map for the current row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// Handle byte arrays (common for strings in some drivers)
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		results = append(results, rowMap)
	}

	return results, nil
}

// Update executes an UPDATE statement and returns the number of rows affected.
func (p *PostgresManager) Update(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := p.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete executes a DELETE statement and returns the number of rows affected.
func (p *PostgresManager) Delete(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := p.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Monitoring Helpers

type PGQuery struct {
	Pid      int    `json:"pid"`
	User     string `json:"user"`
	DB       string `json:"db"`
	State    string `json:"state"`
	Duration string `json:"duration"`
	Query    string `json:"query"`
}

func (p *PostgresManager) GetRunningQueries(ctx context.Context) ([]PGQuery, error) {
	rows, err := p.DB.QueryContext(ctx, `
		SELECT pid, usename, datname, state, (now() - query_start) as duration, query 
		FROM pg_stat_activity 
		WHERE state != 'idle' AND pid <> pg_backend_pid()
		ORDER BY duration DESC LIMIT 50;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queries []PGQuery
	for rows.Next() {
		var q PGQuery
		var user, db, state, query sql.NullString
		var duration sql.NullString
		if err := rows.Scan(&q.Pid, &user, &db, &state, &duration, &query); err != nil {
			continue
		}
		q.User = user.String
		q.DB = db.String
		q.State = state.String
		q.Duration = duration.String
		q.Query = query.String
		queries = append(queries, q)
	}
	return queries, nil
}

func (p *PostgresManager) GetSessionCount(ctx context.Context) (int, error) {
	var count int
	err := p.DB.QueryRowContext(ctx, "SELECT count(*) FROM pg_stat_activity").Scan(&count)
	return count, err
}

func (p *PostgresManager) GetDBInfo(ctx context.Context) (map[string]interface{}, error) {
	var version, dbName, user, sslMode string

	// Fetch Version
	if err := p.DB.QueryRowContext(ctx, "SELECT version()").Scan(&version); err != nil {
		return nil, err
	}

	// Fetch DB Size (formatted)
	var size string
	if err := p.DB.QueryRowContext(ctx, "SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&size); err != nil {
		return nil, err
	}

	// Fetch DB Name
	if err := p.DB.QueryRowContext(ctx, "SELECT current_database()").Scan(&dbName); err != nil {
		return nil, err
	}

	// Fetch Current User
	if err := p.DB.QueryRowContext(ctx, "SELECT current_user").Scan(&user); err != nil {
		return nil, err
	}

	// Fetch SSL Status
	// Note: checks if usage of SSL is active for this backend
	err := p.DB.QueryRowContext(ctx, "SELECT COALESCE((SELECT 'enable' FROM pg_stat_ssl WHERE pid = pg_backend_pid() AND ssl = true), 'disable')").Scan(&sslMode)
	if err != nil {
		sslMode = "unknown"
	}

	return map[string]interface{}{
		"version":  version,
		"size":     size,
		"db_name":  dbName,
		"user":     user,
		"ssl_mode": sslMode,
	}, nil
}
