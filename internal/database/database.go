package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/midgard/gateway/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Type {
	case "sqlite":
		// Use modernc.org/sqlite (pure Go, no CGO required)
		dsn := cfg.DSN
		if dsn == "" {
			dsn = "midgard.db"
		}
		// Add SQLite connection parameters for better concurrency
		// _busy_timeout: sets the busy timeout in milliseconds
		// _journal_mode: WAL mode improves concurrency
		// _foreign_keys: enable foreign key constraints
		if !strings.Contains(dsn, "?") {
			dsn += "?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=1"
		} else {
			dsn += "&_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=1"
		}
		// Open database using modernc.org/sqlite driver
		sqlDB, err := sql.Open("sqlite", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite database: %w", err)
		}
		// Configure connection pool for better concurrency
		sqlDB.SetMaxOpenConns(25)                 // Maximum number of open connections
		sqlDB.SetMaxIdleConns(5)                  // Maximum number of idle connections
		sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime
		// Create GORM DB from sql.DB
		db, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			sqlDB.Close()
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
		if cfg.DSN != "" {
			dsn = cfg.DSN
		}
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	DB = db

	// Auto migrate
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Collection{},
		&Endpoint{},
		&RequestLog{},
	)
}
