package app

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// MigrateFileBasedDiariesToDatabase migrates existing file-based diaries to encrypted database
func MigrateFileBasedDiariesToDatabase(userID int, encryptionKey []byte) error {
	// Check if diaries directory exists
	if _, err := os.Stat(diariesDir); os.IsNotExist(err) {
		return nil // No migration needed
	}

	// Get all existing diaries
	existingDiaries, err := GetDiariesList()
	if err != nil {
		return fmt.Errorf("failed to get existing diaries: %v", err)
	}

	if len(existingDiaries) == 0 {
		return nil // No diaries to migrate
	}

	migratedCount := 0
	failedCount := 0

	// Migrate each diary
	for _, diary := range existingDiaries {
		if err := SaveEncryptedDiary(&diary, uint(userID), encryptionKey); err != nil {
			fmt.Printf("Failed to migrate diary %s: %v\n", diary.ID, err)
			failedCount++
			continue
		}
		migratedCount++
	}

	// Create backup of original files before cleanup
	backupDir := "diaries_backup"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Move original files to backup
	files, err := os.ReadDir(diariesDir)
	if err != nil {
		return fmt.Errorf("failed to read diaries directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			sourcePath := filepath.Join(diariesDir, file.Name())
			backupPath := filepath.Join(backupDir, file.Name())

			if err := os.Rename(sourcePath, backupPath); err != nil {
				fmt.Printf("Failed to backup file %s: %v\n", file.Name(), err)
			}
		}
	}

	fmt.Printf("Migration completed: %d diaries migrated, %d failed\n", migratedCount, failedCount)

	if failedCount > 0 {
		return fmt.Errorf("migration completed with %d failures", failedCount)
	}

	return nil
}

// CheckMigrationNeeded checks if migration is needed
func CheckMigrationNeeded() (bool, int, error) {
	// Check if diaries directory exists and has files
	if _, err := os.Stat(diariesDir); os.IsNotExist(err) {
		return false, 0, nil
	}

	existingDiaries, err := GetDiariesList()
	if err != nil {
		return false, 0, err
	}

	return len(existingDiaries) > 0, len(existingDiaries), nil
}

// GetMigrationStatus returns migration status information
func GetMigrationStatus() (*MigrationStatus, error) {
	needed, count, err := CheckMigrationNeeded()
	if err != nil {
		return nil, err
	}

	return &MigrationStatus{
		MigrationNeeded: needed,
		DiaryCount:      count,
	}, nil
}

// MigrationStatus contains information about migration status
type MigrationStatus struct {
	MigrationNeeded bool `json:"migrationNeeded"`
	DiaryCount      int  `json:"diaryCount"`
}

// RunMigrations runs database migrations
func RunMigrations() error {
	// Add any database schema migrations here
	if err := addEncryptionModeColumns(); err != nil {
		return fmt.Errorf("failed to add encryption mode columns: %v", err)
	}

	return nil
}

// addEncryptionModeColumns adds encryption_mode and encryption_salt columns to encrypted_diaries table
func addEncryptionModeColumns() error {
	// Check if table exists first
	checkTableQuery := `SELECT name FROM sqlite_master WHERE type='table' AND name='encrypted_diaries';`
	var tableName string
	err := db.QueryRow(checkTableQuery).Scan(&tableName)
	if err == sql.ErrNoRows {
		// Table doesn't exist yet, skip migration
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to check table existence: %v", err)
	}

	// Check if columns already exist
	checkQuery := `PRAGMA table_info(encrypted_diaries);`
	rows, err := db.Query(checkQuery)
	if err != nil {
		return fmt.Errorf("failed to get table info: %v", err)
	}
	defer rows.Close()

	hasEncryptionMode := false
	hasEncryptionSalt := false

	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull int
		var defaultValue sql.NullString
		var pk int

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			continue
		}

		if name == "encryption_mode" {
			hasEncryptionMode = true
		}
		if name == "encryption_salt" {
			hasEncryptionSalt = true
		}
	}

	// Add missing columns
	if !hasEncryptionMode {
		addModeQuery := `ALTER TABLE encrypted_diaries ADD COLUMN encryption_mode TEXT NOT NULL DEFAULT 'unified';`
		if _, err := db.Exec(addModeQuery); err != nil {
			return fmt.Errorf("failed to add encryption_mode column: %v", err)
		}
		fmt.Println("Added encryption_mode column to encrypted_diaries table")
	}

	if !hasEncryptionSalt {
		addSaltQuery := `ALTER TABLE encrypted_diaries ADD COLUMN encryption_salt TEXT;`
		if _, err := db.Exec(addSaltQuery); err != nil {
			return fmt.Errorf("failed to add encryption_salt column: %v", err)
		}
		fmt.Println("Added encryption_salt column to encrypted_diaries table")
	}

	return nil
}
