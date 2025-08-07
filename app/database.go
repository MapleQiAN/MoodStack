package app

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *sql.DB
var gormDB *gorm.DB

// InitDatabase initializes the SQLite database
func InitDatabase() error {
	dbPath := filepath.Join("data", "moodstack.db")

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Initialize GORM
	var err error
	gormDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to open database with GORM: %v", err)
	}

	// Get the underlying sql.DB for legacy code compatibility
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}
	db = sqlDB

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Auto-migrate GORM models
	if err := AutoMigrateModels(); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %v", err)
	}

	return nil
}

// AutoMigrateModels automatically migrates all GORM models
func AutoMigrateModels() error {
	// Auto-migrate all models
	err := gormDB.AutoMigrate(
		&User{},
		&EncryptedDiary{},
		&Session{},
		&AppSetting{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate models: %v", err)
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// GetGormDB returns the GORM database connection
func GetGormDB() *gorm.DB {
	return gormDB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
