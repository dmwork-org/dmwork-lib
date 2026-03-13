package redis

import (
	"os"
	"testing"
)

func TestZAddOddLengthError(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		// Test without redis connection - just validate input checking
		conn := &Conn{}
		err := conn.ZAdd("testkey", 1.0, "member1", 2.0)
		if err == nil {
			t.Error("expected error for odd-length scoremember")
		}
		if err.Error() != "scoremember requires pairs of score and member" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
		return
	}

	password := os.Getenv("REDIS_PASSWORD")
	conn := New(addr, password)
	defer conn.Close()

	err := conn.ZAdd("test:zadd:oddlength", 1.0, "member1", 2.0)
	if err == nil {
		t.Error("expected error for odd-length scoremember")
	}
	if err.Error() != "scoremember requires pairs of score and member" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestZAddWrongTypeError(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		// Test without redis connection - just validate input checking
		conn := &Conn{}
		err := conn.ZAdd("testkey", "not-a-float", "member1")
		if err == nil {
			t.Error("expected error for wrong score type")
		}
		if err.Error() != "score must be a float64" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
		return
	}

	password := os.Getenv("REDIS_PASSWORD")
	conn := New(addr, password)
	defer conn.Close()

	err := conn.ZAdd("test:zadd:wrongtype", "not-a-float", "member1")
	if err == nil {
		t.Error("expected error for wrong score type")
	}
	if err.Error() != "score must be a float64" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

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
