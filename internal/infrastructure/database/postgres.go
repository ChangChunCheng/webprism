package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ChangChunCheng/webprism/internal/infrastructure/config"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// PostgresDB wraps a Bun DB instance with additional functionality.
type PostgresDB struct {
	*bun.DB
	logger *logger.Logger
}

// NewPostgresDB creates a new PostgreSQL connection using Bun ORM.
func NewPostgresDB(cfg config.DatabaseConfig, log *logger.Logger) (*PostgresDB, error) {
	// Build DSN string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	// Create pgdriver connector using DSN
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))

	// Create sql.DB
	sqlDB := sql.OpenDB(connector)

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Create Bun DB
	db := bun.NewDB(sqlDB, pgdialect.New())

	// Add query hook for debugging (only in debug mode)
	if log != nil {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(false), // Set to true for verbose query logging
		))
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if log != nil {
		log.Info("PostgreSQL connection established",
			logger.String("host", cfg.Host),
			logger.Int("port", cfg.Port),
			logger.String("database", cfg.Database),
		)
	}

	return &PostgresDB{
		DB:     db,
		logger: log,
	}, nil
}

// Close closes the database connection.
func (db *PostgresDB) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	if db.logger != nil {
		db.logger.Info("PostgreSQL connection closed")
	}

	return nil
}

// Ping checks if the database connection is alive.
func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

// BeginTx starts a new transaction with the given options.
func (db *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}

// RunInTransaction runs the given function in a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (db *PostgresDB) RunInTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return db.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(ctx, tx)
	})
}
