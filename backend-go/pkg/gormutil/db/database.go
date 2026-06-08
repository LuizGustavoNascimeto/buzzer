package gormutil

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds all configuration needed to open a database connection.
type Config struct {
	// Connection pool settings
	SSLMode         string
	MaxOpenConns    int           // maximum number of open connections to the database
	MaxIdleConns    int           // maximum number of idle connections in the pool
	ConnMaxLifetime time.Duration // maximum amount of time a connection may be reused
	ConnMaxIdleTime time.Duration // maximum amount of time a connection may be idle
}

// DefaultConfig returns a Config with sensible pool defaults.
func DefaultConfig() Config {
	return Config{
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// DB is the singleton instance holder.
type DB struct {
	Gorm *gorm.DB
}

var (
	instance *DB
	once     sync.Once
	mu       sync.RWMutex
)

// NewDBConn initialises the singleton connection using the provided Config.
// It is safe to call concurrently; only the first call creates the instance.
// Subsequent calls return an error to signal that a connection is already open.
func NewDBConn(cfg Config, dsn string) (*DB, error) {
	var initErr error

	once.Do(func() {
		gormCfg := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		}
		db, err := gorm.Open(postgres.Open(dsn), gormCfg)
		if err != nil {
			initErr = fmt.Errorf("database: failed to open connection: %w", err)
			return
		}

		// Grab the underlying *sql.DB to configure the pool.
		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("database: failed to retrieve sql.DB: %w", err)
			return
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

		// Verify the connection is actually reachable.
		if err = sqlDB.Ping(); err != nil {
			initErr = fmt.Errorf("database: ping failed: %w", err)
			return
		}

		mu.Lock()
		instance = &DB{Gorm: db}
		mu.Unlock()

		log.Println("database: connection established successfully")
	})

	if initErr != nil {
		// Reset once so a future call may retry after a corrected config.
		once = sync.Once{}
		return nil, initErr
	}

	mu.RLock()
	defer mu.RUnlock()
	return instance, nil
}

// GetConn returns the existing singleton instance.
// If no connection has been initialised yet it calls NewDBConn with the
// supplied Config, so callers that hold a Config can use GetConn exclusively.
func GetConn(cfg ...Config) (*DB, error) {
	mu.RLock()
	inst := instance
	mu.RUnlock()

	if inst != nil {
		return inst, nil
	}

	if len(cfg) == 0 {
		return nil, fmt.Errorf("database: no active connection; provide a Config to initialise one")
	}

	return nil, errors.New("Connection no initialized")
}

// MustGetConn is like GetConn but panics on error — useful in main/init contexts.
func MustGetConn(cfg ...Config) *DB {
	db, err := GetConn(cfg...)
	if err != nil {
		panic(err)
	}
	return db
}

// Ping verifies the connection is still alive.
func (d *DB) Ping() error {
	sqlDB, err := d.Gorm.DB()
	if err != nil {
		return fmt.Errorf("database: failed to retrieve sql.DB: %w", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("database: ping failed: %w", err)
	}
	return nil
}

// Stats returns the current connection-pool statistics.
func (d *DB) Stats() (map[string]any, error) {
	sqlDB, err := d.Gorm.DB()
	if err != nil {
		return nil, fmt.Errorf("database: failed to retrieve sql.DB: %w", err)
	}

	s := sqlDB.Stats()
	return map[string]any{
		"max_open_connections": s.MaxOpenConnections,
		"open_connections":     s.OpenConnections,
		"in_use":               s.InUse,
		"idle":                 s.Idle,
		"wait_count":           s.WaitCount,
		"wait_duration":        s.WaitDuration.String(),
	}, nil
}

// Close gracefully shuts down the underlying connection pool and resets the
// singleton so NewDBConn / GetConn may be called again if needed.
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		return nil
	}

	sqlDB, err := instance.Gorm.DB()
	if err != nil {
		return fmt.Errorf("database: failed to retrieve sql.DB on close: %w", err)
	}

	if err = sqlDB.Close(); err != nil {
		return fmt.Errorf("database: failed to close connection: %w", err)
	}

	instance = nil
	once = sync.Once{}

	log.Println("database: connection closed")
	return nil
}
