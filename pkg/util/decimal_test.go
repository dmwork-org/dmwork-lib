package util

import "testing"

func TestRoundCash15(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic cases - interval 15 means 5 gets rounded down
		{"3.45", "3.4"},
		{"3.35", "3.3"},
		{"3.25", "3.2"},
		{"3.15", "3.1"},
		{"3.05", "3"},
		// Non-multiples of 5 should round normally to nearest 10 cents
		{"3.43", "3.4"},
		{"3.47", "3.5"},
		{"3.44", "3.4"},
		{"3.46", "3.5"},
		// Edge cases
		{"0.05", "0"},
		{"0.15", "0.1"},
		{"10.25", "10.2"},
		{"100.45", "100.4"},
	}

	for _, tt := range tests {
		d := RequireFromString(tt.input)
		result := d.RoundCash(15)
		if result.String() != tt.expected {
			t.Errorf("RoundCash(15) for %s: got %s, want %s", tt.input, result.String(), tt.expected)
		}
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

	for _, tt := range tests {
		d := RequireFromString(tt.input)
		result := d.RoundCash(5)
		if result.String() != tt.expected {
			t.Errorf("RoundCash(5) for %s: got %s, want %s", tt.input, result.String(), tt.expected)
		}
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

	for _, tt := range tests {
		d := RequireFromString(tt.input)
		result := d.RoundCash(10)
		if result.String() != tt.expected {
			t.Errorf("RoundCash(10) for %s: got %s, want %s", tt.input, result.String(), tt.expected)
		}
	}
}
