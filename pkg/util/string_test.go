package util

import (
	"regexp"
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
		{"length 0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRandomString(tt.length)
			if len(result) != tt.length {
				t.Errorf("GetRandomString(%d) returned length %d, want %d", tt.length, len(result), tt.length)
			}

			// Verify charset - only alphanumeric characters
			if tt.length > 0 {
				matched, _ := regexp.MatchString("^[0-9a-zA-Z]+$", result)
				if !matched {
					t.Errorf("GetRandomString(%d) returned invalid characters: %s", tt.length, result)
				}
			}
		})
	}
}

func TestGetRandomString_Uniqueness(t *testing.T) {
	// Generate multiple strings and verify they are unique (high probability)
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s := GetRandomString(16)
		if seen[s] {
			t.Errorf("GetRandomString produced duplicate: %s", s)
		}
		seen[s] = true
	}
}

func TestGetRandomSalt(t *testing.T) {
	salt := GetRandomSalt()
	if len(salt) != 8 {
		t.Errorf("GetRandomSalt() returned length %d, want 8", len(salt))
	}

	// Verify charset
	matched, _ := regexp.MatchString("^[0-9a-zA-Z]+$", salt)
	if !matched {
		t.Errorf("GetRandomSalt() returned invalid characters: %s", salt)
	}
}

func TestGetRandomSalt_Uniqueness(t *testing.T) {
	// Generate multiple salts and verify they are unique
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		salt := GetRandomSalt()
		if seen[salt] {
			t.Errorf("GetRandomSalt produced duplicate: %s", salt)
		}
		seen[salt] = true
	}
}

func TestGetRandomName(t *testing.T) {
	name := GetRandomName()
	if name == "" {
		t.Error("GetRandomName() returned empty string")
	}

	// Verify the name is from the names slice
	found := false
	for _, n := range names {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetRandomName() returned unknown name: %s", name)
	}
}

func TestUnderscoreName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CamelCase", "camel_case"},
		{"camelCase", "camel_case"},
		{"HTTPServer", "http_server"},
		{"ID", "id"},
		{"userID", "user_id"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnderscoreName(tt.input)
			if result != tt.expected {
				t.Errorf("UnderscoreName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveRepeatedElement(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"no duplicates", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"with duplicates", []string{"a", "b", "a", "c"}, []string{"b", "a", "c"}},
		{"all same", []string{"a", "a", "a"}, []string{"a"}},
		{"empty", []string{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveRepeatedElement(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("RemoveRepeatedElement() returned %v, want %v", result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("RemoveRepeatedElement() returned %v, want %v", result, tt.expected)
					return
				}
			}
		})
	}
}
