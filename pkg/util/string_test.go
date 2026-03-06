package util

import (
	"testing"
)

func TestGetRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 8", 8},
		{"length 16", 16},
		{"length 32", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRandomString(tt.length)
			if len(result) != tt.length {
				t.Errorf("GetRandomString(%d) = %q, want length %d, got %d", tt.length, result, tt.length, len(result))
			}

			// Verify only valid characters are used
			validChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			for _, c := range result {
				found := false
				for _, vc := range validChars {
					if c == vc {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetRandomString(%d) contains invalid character: %c", tt.length, c)
				}
			}
		})
	}
}

func TestGetRandomStringUniqueness(t *testing.T) {
	// Generate multiple random strings and ensure they are unique
	// (with high probability for crypto-secure random)
	results := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		result := GetRandomString(16)
		if results[result] {
			t.Errorf("GetRandomString(16) produced duplicate value: %s", result)
		}
		results[result] = true
	}
}

func TestGetRandomSalt(t *testing.T) {
	salt := GetRandomSalt()
	if len(salt) != 8 {
		t.Errorf("GetRandomSalt() = %q, want length 8, got %d", salt, len(salt))
	}
}

func TestGetRandomSaltUniqueness(t *testing.T) {
	// Generate multiple salts and ensure they are unique
	results := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		result := GetRandomSalt()
		if results[result] {
			t.Errorf("GetRandomSalt() produced duplicate value: %s", result)
		}
		results[result] = true
	}
}
