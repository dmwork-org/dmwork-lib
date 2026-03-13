package util

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestStructAttrToUnderscore(t *testing.T) {
	names, err := AttrToUnderscore(&struct {
		MessageID uint64
		Name      string
		UserAge   int
	}{})
	assert.Equal(t, nil, err)
	assert.Equal(t, "message_id", names[0])
	assert.Equal(t, "name", names[1])
	assert.Equal(t, "user_age", names[2])
}

func TestAttrToUnderscoreNilInput(t *testing.T) {
	_, err := AttrToUnderscore(nil)
	if err == nil {
		t.Error("expected error for nil input")
	}
	if err.Error() != "input cannot be nil" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestAttrToUnderscoreNonPointer(t *testing.T) {
	_, err := AttrToUnderscore(struct{ Name string }{})
	if err == nil {
		t.Error("expected error for non-pointer input")
	}
	if err.Error() != "input must be a pointer to struct" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestAttrToUnderscoreNilPointer(t *testing.T) {
	var s *struct{ Name string }
	_, err := AttrToUnderscore(s)
	if err == nil {
		t.Error("expected error for nil pointer input")
	}
	if err.Error() != "input cannot be nil pointer" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestAttrToUnderscorePointerToNonStruct(t *testing.T) {
	s := "not a struct"
	_, err := AttrToUnderscore(&s)
	if err == nil {
		t.Error("expected error for pointer to non-struct")
	}
	if err.Error() != "input must be a pointer to struct" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
