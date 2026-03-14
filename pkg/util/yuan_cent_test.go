package util

import (
	"testing"
)

func TestCentToYuan(t *testing.T) {
	tests := []struct {
		name     string
		cent     int64
		expected float64
	}{
		{"zero", 0, 0},
		{"one cent", 1, 0.01},
		{"ten cents", 10, 0.1},
		{"one yuan", 100, 1},
		{"mixed", 123, 1.23},
		{"large value", 123456, 1234.56},
		{"negative", -100, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CentToYuan(tt.cent)
			if result != tt.expected {
				t.Errorf("CentToYuan(%d) = %v, want %v", tt.cent, result, tt.expected)
			}
		})
	}
}

func TestYuanToCent(t *testing.T) {
	tests := []struct {
		name     string
		yuan     float64
		expected int64
	}{
		{"zero", 0, 0},
		{"one cent", 0.01, 1},
		{"ten cents", 0.1, 10},
		{"one yuan", 1, 100},
		{"mixed", 1.23, 123},
		{"large value", 1234.56, 123456},
		{"negative", -1, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := YuanToCent(tt.yuan)
			if result != tt.expected {
				t.Errorf("YuanToCent(%v) = %d, want %d", tt.yuan, result, tt.expected)
			}
		})
	}
}
