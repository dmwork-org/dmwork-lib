package db

import (
	"testing"
)

func TestFileDirMigrationSourceFindMigrations(t *testing.T) {
	// Test that findMigrations properly closes file handles
	// This is a basic existence test - the actual functionality
	// would require setting up a directory with .sql files

	source := FileDirMigrationSource{
		Dir: "testdata", // non-existent, will return error
	}

	// Should not panic even with non-existent directory
	_, err := source.FindMigrations()
	if err == nil {
		t.Log("FindMigrations returned nil error for non-existent directory")
	}
}
