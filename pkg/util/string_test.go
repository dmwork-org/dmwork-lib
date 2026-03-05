package util

import "testing"

func TestCamelName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple underscore",
			input:    "hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "multiple underscores",
			input:    "foo_bar_baz",
			expected: "FooBarBaz",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "already camel case",
			input:    "HelloWorld",
			expected: "Helloworld",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "leading underscore",
			input:    "_hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "trailing underscore",
			input:    "hello_world_",
			expected: "HelloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CamelName(tt.input)
			if result != tt.expected {
				t.Errorf("CamelName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnderscoreName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple camel case",
			input:    "HelloWorld",
			expected: "hello_world",
		},
		{
			name:     "multiple words",
			input:    "FooBarBaz",
			expected: "foo_bar_baz",
		},
		{
			name:     "single word lowercase",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "single word uppercase start",
			input:    "Hello",
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnderscoreName(tt.input)
			if result != tt.expected {
				t.Errorf("UnderscoreName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
