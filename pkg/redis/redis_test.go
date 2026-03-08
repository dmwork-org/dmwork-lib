package redis

import (
	"os"
	"strings"
	"testing"

	rd "github.com/go-redis/redis"
)

// TestHmgetNilHandling tests that Hmget correctly handles nil values
// for non-existent fields instead of panicking.
func TestHmgetNilHandling(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		t.Skip("REDIS_ADDR not set, skipping integration test")
	}

	password := os.Getenv("REDIS_PASSWORD")
	conn := New(addr, password)

	// Clean up test key
	testKey := "test:hmget:nilhandling"
	defer conn.Del(testKey)

	// Set only one field
	err := conn.Hset(testKey, "field1", "value1")
	if err != nil {
		t.Fatalf("Hset failed: %v", err)
	}

	// Request two fields - field2 does not exist (will be nil)
	// Before the fix, this would panic with: interface conversion: interface {} is nil, not string
	results, err := conn.Hmget(testKey, "field1", "field2")
	if err != nil {
		t.Fatalf("Hmget failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if results[0] != "value1" {
		t.Errorf("expected 'value1' for field1, got '%s'", results[0])
	}

	// Non-existent field should return empty string, not panic
	if results[1] != "" {
		t.Errorf("expected empty string for non-existent field2, got '%s'", results[1])
	}
}

// TestZAddOddNumberOfArgs tests that ZAdd returns an error when given odd number of arguments
// instead of panicking with array index out of bounds.
func TestZAddOddNumberOfArgs(t *testing.T) {
	// Create a connection with dummy address - we're testing parameter validation
	// which happens before any Redis call
	conn := &Conn{}

	// Test with odd number of arguments (missing member for last score)
	err := conn.ZAdd("testkey", 1.0, "member1", 2.0)
	if err == nil {
		t.Error("expected error for odd number of arguments, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "even number") {
		t.Errorf("expected error about even number of elements, got: %v", err)
	}
}

// TestZAddInvalidScoreType tests that ZAdd returns an error when score is not float64
// instead of panicking with type assertion failure.
func TestZAddInvalidScoreType(t *testing.T) {
	// Create a connection with dummy address - we're testing parameter validation
	// which happens before any Redis call
	conn := &Conn{}

	// Test with invalid score type (string instead of float64)
	err := conn.ZAdd("testkey", "not_a_float", "member1")
	if err == nil {
		t.Error("expected error for invalid score type, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "must be float64") {
		t.Errorf("expected error about float64 type, got: %v", err)
	}
}

// TestZAddValidInput tests that ZAdd accepts valid input without error (parameter validation only).
func TestZAddValidInput(t *testing.T) {
	// Create a connection with dummy address - we're testing parameter validation only
	conn := &Conn{}

	// Valid input should pass parameter validation
	// Note: This will panic at the Redis call since client is nil, but we're only
	// testing that parameter validation passes
	defer func() {
		if r := recover(); r != nil {
			// Expected: nil client panic means validation passed
			// This is acceptable for a unit test without real Redis
		}
	}()

	err := conn.ZAdd("testkey", 1.0, "member1", 2.0, "member2")

	// If we get a parameter validation error, that's a test failure
	if err != nil && (strings.Contains(err.Error(), "even number") ||
		strings.Contains(err.Error(), "must be float64")) {
		t.Errorf("unexpected parameter validation error: %v", err)
	}
}

// TestZAddIntegration tests ZAdd with a real Redis connection.
func TestZAddIntegration(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		t.Skip("REDIS_ADDR not set, skipping integration test")
	}

	password := os.Getenv("REDIS_PASSWORD")
	conn := New(addr, password)
	defer conn.Close()

	testKey := "test:zadd:integration"
	defer conn.Del(testKey)

	// Test valid ZAdd
	err := conn.ZAdd(testKey, 1.0, "member1", 2.5, "member2")
	if err != nil {
		t.Fatalf("ZAdd failed: %v", err)
	}

	// Verify the data was added correctly
	results, err := conn.ZRangeByScore(testKey, rd.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	})
	if err != nil {
		t.Fatalf("ZRangeByScore failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 members, got %d", len(results))
	}
}
