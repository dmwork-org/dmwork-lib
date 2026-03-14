package config

import (
	"encoding/json"
	"testing"
)

// simulateAPIResponse simulates receiving a raw JSON message response from the
// WuKongIM API and unmarshaling it into a MessageResp, exactly as the real
// application would. This is the shared E2E fixture — each test only varies the
// payload content.
func simulateAPIResponse(t *testing.T, payloadJSON string) (result int, panicked bool) {
	t.Helper()

	// Build a full MessageResp as it would arrive from the API, with the
	// Payload field set to the raw JSON bytes the server would send.
	msg := &MessageResp{
		Payload: []byte(payloadJSON),
	}

	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()

	result = msg.GetContentType()
	return result, false
}

// TestE2E_MessageProcessing_MissingTypeField simulates a real API consumer
// receiving a message whose payload has no "type" field (e.g., a system
// notification or a custom payload). GetContentType should return 0, not panic.
func TestE2E_MessageProcessing_MissingTypeField(t *testing.T) {
	result, panicked := simulateAPIResponse(t, `{"content": "hello", "from": "user123"}`)
	if panicked {
		t.Fatal("PANIC: GetContentType crashed when processing a message without a 'type' field — this would crash the application in production")
	}
	if result != 0 {
		t.Errorf("expected content type 0 for message without 'type' field, got %d", result)
	}
}

// TestE2E_MessageProcessing_StringTypeField simulates a real API consumer
// receiving a message where "type" is a string (e.g., "text" instead of a
// numeric content type). GetContentType should return 0, not panic.
func TestE2E_MessageProcessing_StringTypeField(t *testing.T) {
	result, panicked := simulateAPIResponse(t, `{"type": "text", "content": "hello"}`)
	if panicked {
		t.Fatal("PANIC: GetContentType crashed when 'type' is a string — this would crash the application when processing messages from a non-conforming sender")
	}
	if result != 0 {
		t.Errorf("expected content type 0 for string 'type', got %d", result)
	}
}

// TestE2E_MessageProcessing_ValidNumericType simulates the happy path: a
// well-formed message with a numeric "type" field. This is a regression guard
// that should pass both before and after the fix.
func TestE2E_MessageProcessing_ValidNumericType(t *testing.T) {
	result, panicked := simulateAPIResponse(t, `{"type": 5, "content": "hello"}`)
	if panicked {
		t.Fatal("PANIC: GetContentType crashed even with a valid numeric 'type' field")
	}
	if result != 5 {
		t.Errorf("expected content type 5, got %d", result)
	}
}

// TestE2E_MessageProcessing_EmptyPayload simulates receiving a message with an
// empty JSON payload (e.g., a heartbeat or presence update). GetContentType
// should return 0, not panic.
func TestE2E_MessageProcessing_EmptyPayload(t *testing.T) {
	result, panicked := simulateAPIResponse(t, `{}`)
	if panicked {
		t.Fatal("PANIC: GetContentType crashed on empty payload — this would crash the application when processing heartbeat/presence messages")
	}
	if result != 0 {
		t.Errorf("expected content type 0 for empty payload, got %d", result)
	}
}

// TestE2E_MessageBatchProcessing simulates a real-world scenario where the
// application processes a batch of messages with mixed payload formats. A single
// malformed message should not crash the entire batch processor.
func TestE2E_MessageBatchProcessing(t *testing.T) {
	// Simulate a batch of messages as they might arrive from a channel history API
	payloads := []struct {
		name    string
		json    string
		wantType int
	}{
		{"valid_text", `{"type": 1, "content": "hello"}`, 1},
		{"valid_image", `{"type": 2, "url": "img.png"}`, 2},
		{"missing_type", `{"content": "system notification"}`, 0},
		{"string_type", `{"type": "custom"}`, 0},
		{"null_type", `{"type": null}`, 0},
		{"valid_video", `{"type": 3, "url": "vid.mp4"}`, 3},
		{"empty", `{}`, 0},
	}

	// Ensure json.Number is used (matching the real decoder path)
	_ = json.Number("1")

	for _, p := range payloads {
		t.Run(p.name, func(t *testing.T) {
			result, panicked := simulateAPIResponse(t, p.json)
			if panicked {
				t.Fatalf("PANIC processing message %q — in production this would crash the batch processor and lose all subsequent messages", p.name)
			}
			if result != p.wantType {
				t.Errorf("message %q: expected content type %d, got %d", p.name, p.wantType, result)
			}
		})
	}
}
