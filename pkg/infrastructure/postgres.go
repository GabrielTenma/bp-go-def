package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"test-go/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresManager struct {
	DB  *sql.DB
	ORM *gorm.DB
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
	var version string
	p.DB.QueryRowContext(ctx, "SELECT version()").Scan(&version)

	var size string
	p.DB.QueryRowContext(ctx, "SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&size)

	return map[string]interface{}{
		"version": version,
		"size":    size,
	}, nil
}
