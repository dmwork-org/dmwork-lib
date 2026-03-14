package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindMigrationsClosesFiles(t *testing.T) {
	// Create a temporary directory with SQL migration files
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test SQL files
	sqlContent := `-- +migrate Up
CREATE TABLE test (id INT);
-- +migrate Down
DROP TABLE test;
`
	files := []string{"001_init.sql", "002_add_column.sql"}
	for _, f := range files {
		err := os.WriteFile(filepath.Join(tmpDir, f), []byte(sqlContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", f, err)
		}
	}

	// Create a subdirectory with another SQL file
	subDir := filepath.Join(tmpDir, "sub")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir, "003_sub.sql"), []byte(sqlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sub SQL file: %v", err)
	}

	// Test FindMigrations
	source := FileDirMigrationSource{Dir: tmpDir}
	migrations, err := source.FindMigrations()
	if err != nil {
		t.Fatalf("FindMigrations failed: %v", err)
	}

	// Verify all migrations were found
	if len(migrations) != 3 {
		t.Errorf("Expected 3 migrations, got %d", len(migrations))
	}
}

func TestFindMigrationsNonExistentDir(t *testing.T) {
	source := FileDirMigrationSource{Dir: "/nonexistent/path"}
	_, err := source.FindMigrations()
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestFindMigrationsEmptyDir(t *testing.T) {
	// Create an empty temporary directory
	tmpDir, err := os.MkdirTemp("", "empty_migrations")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	source := FileDirMigrationSource{Dir: tmpDir}
	migrations, err := source.FindMigrations()
	if err != nil {
		t.Fatalf("FindMigrations failed on empty dir: %v", err)
	}

	if len(migrations) != 0 {
		t.Errorf("Expected 0 migrations, got %d", len(migrations))
	}
}
