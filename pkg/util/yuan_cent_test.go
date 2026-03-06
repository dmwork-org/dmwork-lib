package util

import "testing"

func TestYuanToCent(t *testing.T) {
	tests := []struct {
		name string
		yuan float64
		want int64
	}{
		{"zero", 0.0, 0},
		{"one yuan", 1.0, 100},
		{"decimal", 1.23, 123},
		{"large", 99999.99, 9999999},
		{"small", 0.01, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := YuanToCent(tt.yuan)
			if got != tt.want {
				t.Errorf("YuanToCent(%v) = %v, want %v", tt.yuan, got, tt.want)
			}
		})
	}
}

func TestCentToYuan(t *testing.T) {
	tests := []struct {
		name string
		cent int64
		want float64
	}{
		{"zero", 0, 0.0},
		{"one cent", 1, 0.01},
		{"one yuan", 100, 1.0},
		{"decimal", 123, 1.23},
		{"large", 9999999, 99999.99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CentToYuan(tt.cent)
			if got != tt.want {
				t.Errorf("CentToYuan(%v) = %v, want %v", tt.cent, got, tt.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that CentToYuan(YuanToCent(x)) == x for well-formed inputs
	tests := []float64{0.0, 1.0, 1.23, 99.99, 100.00}
	for _, yuan := range tests {
		cent := YuanToCent(yuan)
		result := CentToYuan(cent)
		if result != yuan {
			t.Errorf("Round trip failed: %v -> %v -> %v", yuan, cent, result)
		}
	}
}
