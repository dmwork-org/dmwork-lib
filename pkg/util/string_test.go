package util

import (
	"regexp"
	"testing"
)

func TestGetRandomSalt(t *testing.T) {
	salt := GetRandomSalt()

	// Check length is 8
	if len(salt) != 8 {
		t.Errorf("GetRandomSalt() returned length %d, expected 8", len(salt))
	}

	// Check charset is alphanumeric
	matched, _ := regexp.MatchString("^[0-9a-zA-Z]+$", salt)
	if !matched {
		t.Errorf("GetRandomSalt() returned invalid characters: %s", salt)
	}
}

func TestGetRandomSaltUniqueness(t *testing.T) {
	// Generate multiple salts and verify they are unique
	// This tests that crypto/rand is producing different outputs
	salts := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		salt := GetRandomSalt()
		if salts[salt] {
			t.Errorf("GetRandomSalt() produced duplicate salt: %s", salt)
		}
		salts[salt] = true
	}
}

func TestGetRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"zero length", 0},
		{"single char", 1},
		{"salt length", 8},
		{"medium length", 16},
		{"long length", 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRandomString(tt.length)
			if len(result) != tt.length {
				t.Errorf("GetRandomString(%d) returned length %d", tt.length, len(result))
			}

			if tt.length > 0 {
				matched, _ := regexp.MatchString("^[0-9a-zA-Z]+$", result)
				if !matched {
					t.Errorf("GetRandomString(%d) returned invalid characters: %s", tt.length, result)
				}
			}
		})
	}
}

func TestGetRandomStringDistribution(t *testing.T) {
	// Test that the distribution covers various character types
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charCount := make(map[byte]int)

	// Generate enough random strings to test distribution
	for i := 0; i < 1000; i++ {
		s := GetRandomString(62)
		for j := 0; j < len(s); j++ {
			charCount[s[j]]++
		}
	}

	// Verify that most characters from charset appear at least once
	missingChars := 0
	for i := 0; i < len(charset); i++ {
		if charCount[charset[i]] == 0 {
			missingChars++
		}
	}

	// Allow for some variance, but most characters should appear
	if missingChars > 5 {
		t.Errorf("GetRandomString() distribution issue: %d characters never appeared", missingChars)
	}
}

func TestGetRandomName(t *testing.T) {
	// Test that all names can be selected (including the last one)
	// Run multiple iterations to statistically verify the last element is reachable
	lastElement := names[len(names)-1]
	foundLast := false

	// Run enough iterations to have high probability of hitting the last element
	for i := 0; i < 10000; i++ {
		name := GetRandomName()
		if name == lastElement {
			foundLast = true
			break
		}
	}

	if !foundLast {
		t.Errorf("GetRandomName() never returned the last element %q after 10000 iterations", lastElement)
	}
}

func TestGetRandomNameReturnsValidName(t *testing.T) {
	for i := 0; i < 100; i++ {
		name := GetRandomName()
		found := false
		for _, n := range names {
			if n == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetRandomName() returned invalid name: %q", name)
		}
	}
}
