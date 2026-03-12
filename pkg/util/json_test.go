package util

import (
	"strings"
	"testing"
)

func TestToJson(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple map",
			input:    map[string]string{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "struct",
			input:    struct{ Name string }{Name: "test"},
			expected: `{"Name":"test"}`,
		},
		{
			name:     "nil",
			input:    nil,
			expected: "null",
		},
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToJson(tt.input)
			if result != tt.expected {
				t.Errorf("ToJson(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToJson_UnmarshalableTypes(t *testing.T) {
	ch := make(chan int)
	result := ToJson(ch)
	if result != "" {
		t.Errorf("ToJson(channel) = %q, want empty string", result)
	}
}

func TestToJsonSafe(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  string
		expectErr bool
	}{
		{
			name:      "simple map",
			input:     map[string]string{"key": "value"},
			expected:  `{"key":"value"}`,
			expectErr: false,
		},
		{
			name:      "struct",
			input:     struct{ Name string }{Name: "test"},
			expected:  `{"Name":"test"}`,
			expectErr: false,
		},
		{
			name:      "nil",
			input:     nil,
			expected:  "null",
			expectErr: false,
		},
		{
			name:      "nested struct",
			input:     struct{ Data map[string]int }{Data: map[string]int{"count": 42}},
			expected:  `{"Data":{"count":42}}`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToJsonSafe(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ToJsonSafe(%v) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToJsonSafe(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToJsonSafe(%v) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestToJsonSafe_UnmarshalableTypes(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "channel",
			input: make(chan int),
		},
		{
			name:  "function",
			input: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToJsonSafe(tt.input)
			if err == nil {
				t.Errorf("ToJsonSafe(%s) expected error, got nil", tt.name)
			}
			if result != "" {
				t.Errorf("ToJsonSafe(%s) = %q, want empty string on error", tt.name, result)
			}
			if !strings.Contains(err.Error(), "json marshal failed") {
				t.Errorf("ToJsonSafe(%s) error = %q, want error containing 'json marshal failed'", tt.name, err.Error())
			}
		})
	}
}

func TestJsonToMap(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "valid json",
			input:     `{"key": "value", "number": 123}`,
			expectErr: false,
		},
		{
			name:      "empty object",
			input:     `{}`,
			expectErr: false,
		},
		{
			name:      "invalid json",
			input:     `{invalid}`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JsonToMap(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("JsonToMap(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("JsonToMap(%q) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("JsonToMap(%q) returned nil map", tt.input)
				}
			}
		})
	}
}

func TestReadJsonByByte(t *testing.T) {
	var result struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	input := []byte(`{"name": "test", "count": 42}`)
	err := ReadJsonByByte(input, &result)
	if err != nil {
		t.Errorf("ReadJsonByByte unexpected error: %v", err)
	}
	if result.Name != "test" {
		t.Errorf("ReadJsonByByte Name = %q, want %q", result.Name, "test")
	}
	if result.Count != 42 {
		t.Errorf("ReadJsonByByte Count = %d, want %d", result.Count, 42)
	}
}
