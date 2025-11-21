package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/akmadan/throome/pkg/adapters"
	"github.com/akmadan/throome/pkg/cluster"
)

// PostgresAdapter implements the DatabaseAdapter interface for PostgreSQL
type PostgresAdapter struct {
	*adapters.BaseAdapter
	config *cluster.ServiceConfig
	pool   *pgxpool.Pool
}

// NewPostgresAdapter creates a new PostgreSQL adapter
func NewPostgresAdapter(config *cluster.ServiceConfig) (adapters.Adapter, error) {
	adapter := &PostgresAdapter{
		BaseAdapter: adapters.NewBaseAdapter(config),
		config:      config,
	}
	return adapter, nil
}

// Connect establishes a connection pool to PostgreSQL
func (p *PostgresAdapter) Connect(ctx context.Context) error {
	// Build connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		p.config.Username,
		p.config.Password,
		p.config.Host,
		p.config.Port,
		p.config.Database,
	)

	// Parse config
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure pool
	if p.config.Pool.MaxConnections > 0 {
		poolConfig.MaxConns = int32(p.config.Pool.MaxConnections)
	}
	if p.config.Pool.MinConnections > 0 {
		poolConfig.MinConns = int32(p.config.Pool.MinConnections)
	}
	if p.config.Pool.MaxIdleTime > 0 {
		poolConfig.MaxConnIdleTime = time.Duration(p.config.Pool.MaxIdleTime) * time.Second
	}
	if p.config.Pool.MaxLifetime > 0 {
		poolConfig.MaxConnLifetime = time.Duration(p.config.Pool.MaxLifetime) * time.Second
	}

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	p.pool = pool

	// Test connection
	if err := p.Ping(ctx); err != nil {
		p.pool.Close()
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	p.SetConnected(true)
	return nil
}

// Disconnect closes the PostgreSQL connection pool
func (p *PostgresAdapter) Disconnect(ctx context.Context) error {
	if p.pool != nil {
		p.pool.Close()
		p.SetConnected(false)
	}
	return nil
}

// Ping checks if the PostgreSQL connection is alive
func (p *PostgresAdapter) Ping(ctx context.Context) error {
	start := time.Now()
	err := p.pool.Ping(ctx)
	p.RecordRequest(time.Since(start), err == nil)
	return err
}

// HealthCheck performs a health check
func (p *PostgresAdapter) HealthCheck(ctx context.Context) (*adapters.HealthStatus, error) {
	start := time.Now()
	err := p.Ping(ctx)
	responseTime := time.Since(start)

	status := &adapters.HealthStatus{
		Healthy:      err == nil,
		ResponseTime: responseTime,
		LastChecked:  time.Now(),
	}

	if err != nil {
		status.ErrorMessage = err.Error()
	}

	return status, nil
}

// Execute executes a query/command
func (p *PostgresAdapter) Execute(ctx context.Context, query string, args ...interface{}) (adapters.Result, error) {
	start := time.Now()
	tag, err := p.pool.Exec(ctx, query, args...)
	p.RecordRequest(time.Since(start), err == nil)

	if err != nil {
		return nil, err
	}

	return &postgresResult{tag: tag}, nil
}

// Query performs a query and returns rows
func (p *PostgresAdapter) Query(ctx context.Context, query string, args ...interface{}) (adapters.Rows, error) {
	start := time.Now()
	rows, err := p.pool.Query(ctx, query, args...)
	p.RecordRequest(time.Since(start), err == nil)

	if err != nil {
		return nil, err
	}

	return &postgresRows{rows: rows}, nil
}

// QueryRow performs a query that returns a single row
func (p *PostgresAdapter) QueryRow(ctx context.Context, query string, args ...interface{}) adapters.Row {
	start := time.Now()
	row := p.pool.QueryRow(ctx, query, args...)
	p.RecordRequest(time.Since(start), true) // Record as success since error is deferred

	return &postgresRow{row: row}
}

// Begin starts a transaction
func (p *PostgresAdapter) Begin(ctx context.Context) (adapters.Transaction, error) {
	start := time.Now()
	tx, err := p.pool.Begin(ctx)
	p.RecordRequest(time.Since(start), err == nil)

	if err != nil {
		return nil, err
	}

	return &postgresTransaction{tx: tx, adapter: p}, nil
}

// postgresResult implements adapters.Result
type postgresResult struct {
	tag pgconn.CommandTag
}

func (r *postgresResult) RowsAffected() int64 {
	return r.tag.RowsAffected()
}

func (r *postgresResult) LastInsertID() int64 {
	// PostgreSQL doesn't have a native LastInsertID concept
	// This would need to be handled differently (e.g., RETURNING clause)
	return 0
}

// postgresRows implements adapters.Rows
type postgresRows struct {
	rows pgx.Rows
}

func (r *postgresRows) Next() bool {
	return r.rows.Next()
}

func (r *postgresRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *postgresRows) Close() error {
	r.rows.Close()
	return nil
}

func (r *postgresRows) Err() error {
	return r.rows.Err()
}

// postgresRow implements adapters.Row
type postgresRow struct {
	row pgx.Row
}

func (r *postgresRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

// postgresTransaction implements adapters.Transaction
type postgresTransaction struct {
	tx      pgx.Tx
	adapter *PostgresAdapter
}

func (t *postgresTransaction) Commit() error {
	return t.tx.Commit(context.Background())
}

func (t *postgresTransaction) Rollback() error {
	return t.tx.Rollback(context.Background())
}

func (t *postgresTransaction) Execute(ctx context.Context, query string, args ...interface{}) (adapters.Result, error) {
	start := time.Now()
	tag, err := t.tx.Exec(ctx, query, args...)
	t.adapter.RecordRequest(time.Since(start), err == nil)

	if err != nil {
		return nil, err
	}

	return &postgresResult{tag: tag}, nil
}

func (t *postgresTransaction) Query(ctx context.Context, query string, args ...interface{}) (adapters.Rows, error) {
	start := time.Now()
	rows, err := t.tx.Query(ctx, query, args...)
	t.adapter.RecordRequest(time.Since(start), err == nil)

	if err != nil {
		return nil, err
	}

	return &postgresRows{rows: rows}, nil
}

// GetPoolStats returns connection pool statistics
func (p *PostgresAdapter) GetPoolStats() *pgxpool.Stat {
	if p.pool == nil {
		return nil
	}
	stat := p.pool.Stat()
	return stat
}

// Ensure PostgresAdapter implements DatabaseAdapter
var _ adapters.DatabaseAdapter = (*PostgresAdapter)(nil)
