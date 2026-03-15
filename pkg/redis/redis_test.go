package redis

import (
	"os"
	"testing"
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

// TestMSetMultipleKeyValues tests that MSet properly sets multiple key-value pairs.
// Before the fix, the []string slice was passed as a single interface{} instead of
// being expanded, causing the operation to fail or produce unexpected results.
func TestMSetMultipleKeyValues(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		t.Skip("REDIS_ADDR not set, skipping integration test")
	}

	password := os.Getenv("REDIS_PASSWORD")
	conn := New(addr, password)

	// Clean up test keys
	testKey1 := "test:mset:key1"
	testKey2 := "test:mset:key2"
	testKey3 := "test:mset:key3"
	defer conn.Del(testKey1)
	defer conn.Del(testKey2)
	defer conn.Del(testKey3)

	// Set multiple key-value pairs
	err := conn.MSet(testKey1, "value1", testKey2, "value2", testKey3, "value3")
	if err != nil {
		t.Fatalf("MSet failed: %v", err)
	}

	// Verify all values were set correctly
	val1, err := conn.GetString(testKey1)
	if err != nil {
		t.Fatalf("GetString for key1 failed: %v", err)
	}
	if val1 != "value1" {
		t.Errorf("expected 'value1' for key1, got '%s'", val1)
	}

	val2, err := conn.GetString(testKey2)
	if err != nil {
		t.Fatalf("GetString for key2 failed: %v", err)
	}
	if val2 != "value2" {
		t.Errorf("expected 'value2' for key2, got '%s'", val2)
	}

	val3, err := conn.GetString(testKey3)
	if err != nil {
		t.Fatalf("GetString for key3 failed: %v", err)
	}
	if val3 != "value3" {
		t.Errorf("expected 'value3' for key3, got '%s'", val3)
	}
}
