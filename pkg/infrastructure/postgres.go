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
