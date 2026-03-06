package util

import (
	"testing"
)

func TestTenValue2Char_UniqueMapping(t *testing.T) {
	// Test that all values 0-61 map to unique characters
	seen := make(map[string]int64)

	for i := int64(0); i <= 61; i++ {
		char := tenValue2Char(i)
		if char == "" {
			t.Errorf("tenValue2Char(%d) returned empty string", i)
			continue
		}
		if prev, exists := seen[char]; exists {
			t.Errorf("Duplicate mapping: both %d and %d map to %q", prev, i, char)
		}
		seen[char] = i
	}

	// Should have exactly 62 unique characters
	if len(seen) != 62 {
		t.Errorf("Expected 62 unique characters, got %d", len(seen))
	}
}

func TestTenValue2Char_CorrectValues(t *testing.T) {
	testCases := []struct {
		input    int64
		expected string
	}{
		// 0-9 -> "0"-"9"
		{0, "0"},
		{1, "1"},
		{9, "9"},
		// 10-35 -> "a"-"z"
		{10, "a"},
		{11, "b"},
		{33, "x"}, // This was the bug: was returning "s"
		{34, "y"},
		{35, "z"},
		// 36-61 -> "A"-"Z"
		{36, "A"},
		{54, "S"},
		{59, "X"}, // This was the bug: was returning "S"
		{60, "Y"},
		{61, "Z"},
	}

	for _, tc := range testCases {
		result := tenValue2Char(tc.input)
		if result != tc.expected {
			t.Errorf("tenValue2Char(%d) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestTen2Hex(t *testing.T) {
	testCases := []struct {
		input    int64
		expected string
	}{
		{0, ""},
		{1, "1"},
		{10, "a"},
		{61, "Z"},
		{62, "10"},  // 1*62 + 0
		{63, "11"},  // 1*62 + 1
		{124, "20"}, // 2*62 + 0
	}

	for _, tc := range testCases {
		result := Ten2Hex(tc.input)
		if result != tc.expected {
			t.Errorf("Ten2Hex(%d) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestTen2Hex_RoundTrip(t *testing.T) {
	// Test that encoding produces unique results for different inputs
	results := make(map[string]int64)

	for i := int64(1); i <= 1000; i++ {
		hex := Ten2Hex(i)
		if prev, exists := results[hex]; exists {
			t.Errorf("Ten2Hex collision: %d and %d both produce %q", prev, i, hex)
		}
		results[hex] = i
	}
}
