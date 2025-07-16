package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/imraushankr/bervity/server/src/configs"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps gorm.DB with additional methods and configuration
type DB struct {
	*gorm.DB
	config *configs.DatabaseConfig
}

// ConnectDB establishes a database connection and returns a DB instance
func ConnectDB(cfg *configs.DatabaseConfig) (*DB, error) {
	// Ensure database directory exists
	if err := os.MkdirAll(filepath.Dir(cfg.SQLite.Path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Create DSN with SQLite connection parameters
	dsn := fmt.Sprintf("%s?_busy_timeout=%d&_foreign_keys=%t&_journal_mode=%s&_cache_size=%d",
		cfg.SQLite.Path,
		cfg.SQLite.BusyTimeout,
		cfg.SQLite.ForeignKeys,
		cfg.SQLite.JournalMode,
		cfg.SQLite.CacheSize,
	)

	// Initialize GORM with SQLite
	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable GORM's built-in logger
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool settings
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Log successful connection
	zap.L().Info("Database connection established",
		zap.String("path", cfg.SQLite.Path),
		zap.String("journal_mode", cfg.SQLite.JournalMode),
	)

	return &DB{DB: gormDB, config: cfg}, nil
}

// Ping verifies the database connection is alive
func (db *DB) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	zap.L().Info("Closing database connection")
	return sqlDB.Close()
}

// WithTx executes a function within a transaction
func (db *DB) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return db.DB.WithContext(ctx).Transaction(fn)
}
