package config

import (
	"os"
	"testing"
)

func TestGetEnvInt64(t *testing.T) {
	const testKey = "TEST_ENV_INT64"
	const defaultValue int64 = 100

	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		expected int64
	}{
		{
			name:     "empty env returns default",
			envValue: "",
			setEnv:   false,
			expected: defaultValue,
		},
		{
			name:     "valid int64 value",
			envValue: "42",
			setEnv:   true,
			expected: 42,
		},
		{
			name:     "invalid value returns default",
			envValue: "abc",
			setEnv:   true,
			expected: defaultValue,
		},
		{
			name:     "float value returns default",
			envValue: "3.14",
			setEnv:   true,
			expected: defaultValue,
		},
		{
			name:     "whitespace only returns default",
			envValue: "   ",
			setEnv:   true,
			expected: defaultValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(testKey)
			if tt.setEnv {
				os.Setenv(testKey, tt.envValue)
			}
			defer os.Unsetenv(testKey)

			result := GetEnvInt64(testKey, defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvInt64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	const testKey = "TEST_ENV_INT"
	const defaultValue int = 100

	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		expected int
	}{
		{
			name:     "empty env returns default",
			envValue: "",
			setEnv:   false,
			expected: defaultValue,
		},
		{
			name:     "valid int value",
			envValue: "42",
			setEnv:   true,
			expected: 42,
		},
		{
			name:     "invalid value returns default",
			envValue: "abc",
			setEnv:   true,
			expected: defaultValue,
		},
		{
			name:     "float value returns default",
			envValue: "3.14",
			setEnv:   true,
			expected: defaultValue,
		},
		{
			name:     "whitespace only returns default",
			envValue: "   ",
			setEnv:   true,
			expected: defaultValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(testKey)
			if tt.setEnv {
				os.Setenv(testKey, tt.envValue)
			}
			defer os.Unsetenv(testKey)

			result := GetEnvInt(testKey, defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnvFloat64(t *testing.T) {
	const testKey = "TEST_ENV_FLOAT64"
	const defaultValue float64 = 1.5

	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		expected float64
	}{
		{
			name:     "empty env returns default",
			envValue: "",
			setEnv:   false,
			expected: defaultValue,
		},
		{
			name:     "valid float64 value",
			envValue: "3.14",
			setEnv:   true,
			expected: 3.14,
		},
		{
			name:     "valid int value as float",
			envValue: "42",
			setEnv:   true,
			expected: 42.0,
		},
		{
			name:     "invalid value returns default",
			envValue: "abc",
			setEnv:   true,
			expected: defaultValue,
		},
		{
			name:     "whitespace only returns default",
			envValue: "   ",
			setEnv:   true,
			expected: defaultValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(testKey)
			if tt.setEnv {
				os.Setenv(testKey, tt.envValue)
			}
			defer os.Unsetenv(testKey)

			result := GetEnvFloat64(testKey, defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvFloat64() = %v, want %v", result, tt.expected)
			}
		})
	}
}
