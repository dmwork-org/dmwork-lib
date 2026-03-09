package util

import (
	"testing"
)

func TestRoundCash15(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// When the value ends in 5 (e.g., 3.45), interval 15 rounds down
		{"3.45", "3.4"},
		{"3.55", "3.5"},
		{"3.65", "3.6"},
		{"3.75", "3.7"},
		{"3.85", "3.8"},
		{"3.95", "3.9"},
		// Values not ending in 5 round normally
		{"3.41", "3.4"},
		{"3.42", "3.4"},
		{"3.43", "3.4"},
		{"3.44", "3.4"},
		{"3.46", "3.5"},
		{"3.47", "3.5"},
		{"3.48", "3.5"},
		{"3.49", "3.5"},
		// Edge cases
		{"0.05", "0"},
		{"0.15", "0.1"},
		{"1.00", "1"},
		{"1.05", "1"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			d, err := NewFromString(tc.input)
			if err != nil {
				t.Fatalf("NewFromString(%q) failed: %v", tc.input, err)
			}
			result := d.RoundCash(15)
			if result.String() != tc.expected {
				t.Errorf("RoundCash(15) = %q, want %q", result.String(), tc.expected)
			}
		})
	}
}

func TestRoundCash5(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3.43", "3.45"},
		{"3.42", "3.4"},
		{"3.47", "3.45"},
		{"3.48", "3.5"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			d, err := NewFromString(tc.input)
			if err != nil {
				t.Fatalf("NewFromString(%q) failed: %v", tc.input, err)
			}
			result := d.RoundCash(5)
			if result.String() != tc.expected {
				t.Errorf("RoundCash(5) = %q, want %q", result.String(), tc.expected)
			}
		})
	}
}

func TestRoundCash10(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3.45", "3.5"},
		{"3.44", "3.4"},
		{"3.46", "3.5"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			d, err := NewFromString(tc.input)
			if err != nil {
				t.Fatalf("NewFromString(%q) failed: %v", tc.input, err)
			}
			result := d.RoundCash(10)
			if result.String() != tc.expected {
				t.Errorf("RoundCash(10) = %q, want %q", result.String(), tc.expected)
			}
		})
	}
}
