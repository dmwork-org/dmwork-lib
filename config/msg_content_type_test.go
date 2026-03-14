package config

import (
	"encoding/json"
	"testing"
)

// TestGetContentType_MissingTypeKey verifies that GetContentType does not panic
// when the payload JSON has no "type" key. Currently, the unsafe chained type
// assertion `payloadMap["type"].(json.Number).Int64()` panics because the map
// lookup returns nil and the assertion to json.Number fails.
func TestGetContentType_MissingTypeKey(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetContentType panicked on missing 'type' key: %v", r)
		}
	}()

	msg := &MessageResp{
		Payload: []byte(`{"content": "hello"}`),
	}
	result := msg.GetContentType()
	if result != 0 {
		t.Errorf("expected 0 for missing 'type' key, got %d", result)
	}
}

// TestGetContentType_TypeIsString verifies that GetContentType does not panic
// when the "type" field is a string instead of a json.Number. The current code
// performs an unsafe assertion to json.Number which panics on type mismatch.
func TestGetContentType_TypeIsString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetContentType panicked on string 'type' value: %v", r)
		}
	}()

	msg := &MessageResp{
		Payload: []byte(`{"type": "not_a_number"}`),
	}
	result := msg.GetContentType()
	if result != 0 {
		t.Errorf("expected 0 for non-numeric 'type', got %d", result)
	}
}

// TestGetContentType_TypeIsBoolean verifies that GetContentType does not panic
// when the "type" field is a boolean.
func TestGetContentType_TypeIsBoolean(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetContentType panicked on boolean 'type' value: %v", r)
		}
	}()

	msg := &MessageResp{
		Payload: []byte(`{"type": true}`),
	}
	result := msg.GetContentType()
	if result != 0 {
		t.Errorf("expected 0 for boolean 'type', got %d", result)
	}
}

// TestGetContentType_TypeIsNull verifies that GetContentType does not panic
// when the "type" field is null.
func TestGetContentType_TypeIsNull(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetContentType panicked on null 'type' value: %v", r)
		}
	}()

	msg := &MessageResp{
		Payload: []byte(`{"type": null}`),
	}
	result := msg.GetContentType()
	if result != 0 {
		t.Errorf("expected 0 for null 'type', got %d", result)
	}
}

// TestGetContentType_EmptyPayload verifies that GetContentType does not panic
// when the payload is an empty JSON object.
func TestGetContentType_EmptyPayload(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetContentType panicked on empty payload: %v", r)
		}
	}()

	msg := &MessageResp{
		Payload: []byte(`{}`),
	}
	result := msg.GetContentType()
	if result != 0 {
		t.Errorf("expected 0 for empty payload, got %d", result)
	}
}

// TestGetContentType_ValidNumericType verifies that GetContentType correctly
// returns the integer value when a valid numeric "type" field is present.
// This is a regression guard — it should pass both before and after the fix.
func TestGetContentType_ValidNumericType(t *testing.T) {
	msg := &MessageResp{
		Payload: []byte(`{"type": 5, "content": "hello"}`),
	}
	result := msg.GetContentType()
	// With UseNumber(), 5 is decoded as json.Number, so this should work
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

// TestGetContentType_InvalidJSON verifies that GetContentType returns 0
// when the payload is not valid JSON (uses existing error path).
func TestGetContentType_InvalidJSON(t *testing.T) {
	msg := &MessageResp{
		Payload: []byte(`not json at all`),
	}
	result := msg.GetContentType()
	if result != 0 {
		t.Errorf("expected 0 for invalid JSON, got %d", result)
	}
}

// TestGetContentType_FloatType verifies behavior with a floating-point "type"
// value. json.Number.Int64() returns an error for non-integer numbers.
func TestGetContentType_FloatType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetContentType panicked on float 'type' value: %v", r)
		}
	}()

	// Ensure json.Number is used by the decoder (UseNumber is set in ReadJsonByByte)
	_ = json.Number("3.14")

	msg := &MessageResp{
		Payload: []byte(`{"type": 3.14}`),
	}
	result := msg.GetContentType()
	// After fix, this should return 0 since Int64() fails for floats
	// Before fix with current code, the chained assertion still works for json.Number
	// but Int64() error is silently ignored via _, so result would be 0
	_ = result
}
