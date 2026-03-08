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
	names := AttrToUnderscore(nil)
	assert.Equal(t, 0, len(names))
}

func TestAttrToUnderscore_NonPointer(t *testing.T) {
	names := AttrToUnderscore(struct {
		Name string
	}{})
	assert.Equal(t, 0, len(names))
}

func TestAttrToUnderscore_NilPointer(t *testing.T) {
	var s *struct{ Name string }
	names := AttrToUnderscore(s)
	assert.Equal(t, 0, len(names))
}

func TestAttrToUnderscore_PointerToNonStruct(t *testing.T) {
	str := "hello"
	names := AttrToUnderscore(&str)
	assert.Equal(t, 0, len(names))
}
