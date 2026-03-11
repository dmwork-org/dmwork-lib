package util

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestStructAttrToUnderscore(t *testing.T) {

	names := AttrToUnderscore(&struct {
		MessageID uint64
		Name      string
		UserAge   int
	}{})
	assert.Equal(t, "message_id", names[0])
	assert.Equal(t, "name", names[1])
	assert.Equal(t, "user_age", names[2])
}

func TestAttrToUnderscore_NilInput(t *testing.T) {
	// Test nil input - should not panic
	result := AttrToUnderscore(nil)
	if result != nil {
		t.Errorf("expected nil for nil input, got %v", result)
	}
}

func TestAttrToUnderscore_NonPointerInput(t *testing.T) {
	// Test non-pointer input - should not panic
	result := AttrToUnderscore(struct{ Name string }{})
	if result != nil {
		t.Errorf("expected nil for non-pointer input, got %v", result)
	}
}

func TestAttrToUnderscore_NilPointer(t *testing.T) {
	// Test nil pointer - should not panic
	var ptr *struct{ Name string }
	result := AttrToUnderscore(ptr)
	if result != nil {
		t.Errorf("expected nil for nil pointer, got %v", result)
	}
}

func TestAttrToUnderscore_NonStructPointer(t *testing.T) {
	// Test pointer to non-struct - should not panic
	s := "hello"
	result := AttrToUnderscore(&s)
	if result != nil {
		t.Errorf("expected nil for non-struct pointer, got %v", result)
	}
}
