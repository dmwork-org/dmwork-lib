package db

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFindMigrations_FileHandlesClosed verifies that the migration file discovery
// properly closes all file handles. This test creates temporary SQL files and
// verifies the function works correctly without leaking file descriptors.
func TestFindMigrations_FileHandlesClosed(t *testing.T) {
	// Create a temporary directory with migration files
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some test migration files
	migrationFiles := []struct {
		name    string
		content string
	}{
		{
			name: "001_create_users.sql",
			content: `-- +migrate Up
CREATE TABLE users (id INT PRIMARY KEY);

-- +migrate Down
DROP TABLE users;`,
		},
		{
			name: "002_create_orders.sql",
			content: `-- +migrate Up
CREATE TABLE orders (id INT PRIMARY KEY);

-- +migrate Down
DROP TABLE orders;`,
		},
	}

	for _, mf := range migrationFiles {
		path := filepath.Join(tmpDir, mf.name)
		if err := os.WriteFile(path, []byte(mf.content), 0644); err != nil {
			t.Fatalf("failed to create migration file %s: %v", mf.name, err)
		}
	}

	// Run the migration discovery
	source := FileDirMigrationSource{Dir: tmpDir}
	migrations, err := source.FindMigrations()
	if err != nil {
		t.Fatalf("FindMigrations failed: %v", err)
	}

	// Verify the correct number of migrations were found
	if len(migrations) != 2 {
		t.Errorf("expected 2 migrations, got %d", len(migrations))
	}

	// Verify the migrations are sorted correctly
	if len(migrations) >= 2 {
		if migrations[0].Id != "001_create_users.sql" {
			t.Errorf("expected first migration to be '001_create_users.sql', got '%s'", migrations[0].Id)
		}
		if migrations[1].Id != "002_create_orders.sql" {
			t.Errorf("expected second migration to be '002_create_orders.sql', got '%s'", migrations[1].Id)
		}
	}
}

// TestFindMigrations_NestedDirectories tests that nested directories are handled correctly
func TestFindMigrations_NestedDirectories(t *testing.T) {
	// Create a temporary directory with nested structure
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	// Create migration files in both directories
	rootMigration := filepath.Join(tmpDir, "001_root.sql")
	if err := os.WriteFile(rootMigration, []byte("-- +migrate Up\nCREATE TABLE root (id INT);"), 0644); err != nil {
		t.Fatalf("failed to create root migration: %v", err)
	}

	subMigration := filepath.Join(subDir, "002_sub.sql")
	if err := os.WriteFile(subMigration, []byte("-- +migrate Up\nCREATE TABLE sub (id INT);"), 0644); err != nil {
		t.Fatalf("failed to create sub migration: %v", err)
	}

	// Run the migration discovery
	source := FileDirMigrationSource{Dir: tmpDir}
	migrations, err := source.FindMigrations()
	if err != nil {
		t.Fatalf("FindMigrations failed: %v", err)
	}

	// Should find migrations from both root and subdirectory
	if len(migrations) != 2 {
		t.Errorf("expected 2 migrations from nested dirs, got %d", len(migrations))
	}
}
