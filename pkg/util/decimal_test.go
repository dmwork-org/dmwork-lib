package util

import (
	"testing"
)

func TestRoundCash_Interval15(t *testing.T) {
	// interval=15: 10 cent rounding where 5 gets rounded down
	// e.g., 3.45 => 3.40 (5 is rounded down, not up)
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"1.25 rounds down to 1.20", "1.25", "1.2"},
		{"1.35 rounds down to 1.30", "1.35", "1.3"},
		{"1.20 stays 1.20", "1.20", "1.2"},
		{"1.30 stays 1.30", "1.30", "1.3"},
		{"2.55 rounds down to 2.50", "2.55", "2.5"},
		{"2.65 rounds down to 2.60", "2.65", "2.6"},
		{"0.05 rounds down to 0.00", "0.05", "0"},
		{"0.15 rounds down to 0.10", "0.15", "0.1"},
		{"10.25 rounds down to 10.20", "10.25", "10.2"},
		{"10.35 rounds down to 10.30", "10.35", "10.3"},
		{"3.45 rounds down to 3.40 (from docs)", "3.45", "3.4"},
		{"1.23 rounds to 1.20", "1.23", "1.2"},
		{"1.28 rounds to 1.30", "1.28", "1.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewFromString(tt.input)
			if err != nil {
				t.Fatalf("failed to parse decimal %s: %v", tt.input, err)
			}
			result := d.RoundCash(15)
			if result.String() != tt.expected {
				t.Errorf("RoundCash(15) for %s = %s, want %s", tt.input, result.String(), tt.expected)
			}
		})
	}
}

// TestRoundCash_Interval15_NegativeExp tests that RoundCash(15) works correctly
// with negative exponents. This was broken when using XOR (^) instead of proper
// exponentiation - the fix uses New(1, orgExp) to represent the minimum unit.
func TestRoundCash_Interval15_NegativeExp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// These cases would fail with XOR bug: 10^2 vs 10 XOR 2 = 8
		{"exp=-2: 1.05 rounds to 1.00", "1.05", "1"},
		{"exp=-2: 99.95 rounds to 99.90", "99.95", "99.9"},
		{"exp=-3: 1.005 rounds to 1.000", "1.005", "1"},
		{"exp=-1: 1.5 rounds to 1.4", "1.5", "1.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewFromString(tt.input)
			if err != nil {
				t.Fatalf("failed to parse decimal %s: %v", tt.input, err)
			}
			result := d.RoundCash(15)
			if result.String() != tt.expected {
				t.Errorf("RoundCash(15) for %s = %s, want %s", tt.input, result.String(), tt.expected)
			}
		})
	}
}

func TestRoundCash_AllIntervals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		interval uint8
		expected string
	}{
		{"interval 5: 1.23 -> 1.25", "1.23", 5, "1.25"},
		{"interval 10: 1.23 -> 1.20", "1.23", 10, "1.2"},
		{"interval 15: 1.25 -> 1.20", "1.25", 15, "1.2"},
		{"interval 25: 1.13 -> 1.25", "1.13", 25, "1.25"},
		{"interval 50: 1.23 -> 1.00", "1.23", 50, "1"},
		{"interval 100: 1.50 -> 2.00", "1.50", 100, "2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewFromString(tt.input)
			if err != nil {
				t.Fatalf("failed to parse decimal %s: %v", tt.input, err)
			}
			result := d.RoundCash(tt.interval)
			if result.String() != tt.expected {
				t.Errorf("RoundCash(%d) for %s = %s, want %s", tt.interval, tt.input, result.String(), tt.expected)
			}
		})
	}
}
