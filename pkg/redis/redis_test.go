package redis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestRedisConn(t *testing.T) *Conn {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")

	conn := New(addr, password)
	_, err := conn.Ping()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	return conn
}

func TestHmget_NilValues(t *testing.T) {
	conn := getTestRedisConn(t)

	testKey := "test:hmget:nil:values"
	defer conn.Del(testKey)

	// Set only field1, field2 will not exist
	err := conn.Hset(testKey, "field1", "value1")
	assert.NoError(t, err)

	// Request both fields - field2 should return empty string instead of panic
	results, err := conn.Hmget(testKey, "field1", "field2")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "value1", results[0])
	assert.Equal(t, "", results[1]) // nil values should be converted to empty string
}

func TestHmget_AllNilValues(t *testing.T) {
	conn := getTestRedisConn(t)

	testKey := "test:hmget:all:nil"
	defer conn.Del(testKey)

	// Key doesn't exist, all values should be nil -> empty strings
	results, err := conn.Hmget(testKey, "field1", "field2")
	assert.NoError(t, err)
	// HMGet returns array with nil values for non-existent key
	if results != nil {
		assert.Len(t, results, 2)
		assert.Equal(t, "", results[0])
		assert.Equal(t, "", results[1])
	}
}

func TestHmget_AllValidValues(t *testing.T) {
	conn := getTestRedisConn(t)

	testKey := "test:hmget:valid"
	defer conn.Del(testKey)

	err := conn.Hset(testKey, "field1", "value1")
	assert.NoError(t, err)
	err = conn.Hset(testKey, "field2", "value2")
	assert.NoError(t, err)

	results, err := conn.Hmget(testKey, "field1", "field2")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "value1", results[0])
	assert.Equal(t, "value2", results[1])
}
