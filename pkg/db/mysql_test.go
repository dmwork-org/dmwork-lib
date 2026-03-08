package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileDirMigrationSource_FindMigrations_ClosesFiles(t *testing.T) {
	// Create a temporary directory with test SQL files
	tempDir, err := os.MkdirTemp("", "migration_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test SQL migration files
	testFiles := []string{
		"001_create_users.sql",
		"002_create_orders.sql",
	}

	for _, name := range testFiles {
		content := `-- +migrate Up
CREATE TABLE test (id INT);
-- +migrate Down
DROP TABLE test;
`
		err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", name, err)
		}
	}

	// Run FindMigrations multiple times to ensure no file handle leaks
	source := FileDirMigrationSource{Dir: tempDir}

	for i := 0; i < 100; i++ {
		migrations, err := source.FindMigrations()
		if err != nil {
			t.Fatalf("iteration %d: FindMigrations failed: %v", i, err)
		}
		if len(migrations) != len(testFiles) {
			t.Errorf("iteration %d: expected %d migrations, got %d", i, len(testFiles), len(migrations))
		}
	}
}

func TestFileDirMigrationSource_FindMigrations_NestedDir(t *testing.T) {
	// Create a temporary directory with nested structure
	tempDir, err := os.MkdirTemp("", "migration_nested_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create nested directory
	nestedDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(nestedDir, 0755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}

	// Create test SQL migration files in both directories
	content := `-- +migrate Up
CREATE TABLE test (id INT);
-- +migrate Down
DROP TABLE test;
`
	if err := os.WriteFile(filepath.Join(tempDir, "001_root.sql"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to create root SQL file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "002_nested.sql"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to create nested SQL file: %v", err)
	}

	// Run FindMigrations
	source := FileDirMigrationSource{Dir: tempDir}
	migrations, err := source.FindMigrations()
	if err != nil {
		t.Fatalf("FindMigrations failed: %v", err)
	}

	// Should find both migrations
	if len(migrations) != 2 {
		t.Errorf("expected 2 migrations, got %d", len(migrations))
	}
}
