package util

import (
	"testing"
)

func TestYuanToCent(t *testing.T) {
	tests := []struct {
		yuan float64
		want int64
	}{
		{1.00, 100},
		{0.01, 1},
		{10.50, 1050},
		{0.00, 0},
		{99.99, 9999},
		{123.45, 12345},
	}

	for _, tt := range tests {
		got := YuanToCent(tt.yuan)
		if got != tt.want {
			t.Errorf("YuanToCent(%v) = %v, want %v", tt.yuan, got, tt.want)
		}
	}
}

func TestCentToYuan(t *testing.T) {
	tests := []struct {
		cent int64
		want float64
	}{
		{100, 1.00},
		{1, 0.01},
		{1050, 10.50},
		{0, 0.00},
		{9999, 99.99},
		{12345, 123.45},
	}

	for _, tt := range tests {
		got := CentToYuan(tt.cent)
		if got != tt.want {
			t.Errorf("CentToYuan(%v) = %v, want %v", tt.cent, got, tt.want)
		}
	}
}

func TestYuanToCentAndBack(t *testing.T) {
	// Test round-trip conversion
	tests := []float64{1.00, 0.01, 10.50, 99.99, 123.45}

	for _, yuan := range tests {
		cent := YuanToCent(yuan)
		back := CentToYuan(cent)
		if back != yuan {
			t.Errorf("Round-trip failed: %v -> %v -> %v", yuan, cent, back)
		}
	}
}
