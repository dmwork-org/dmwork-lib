package util

import (
	"testing"
	"time"
)

func TestToyyyyMMddHHmm(t *testing.T) {
	tm := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	expected := "2024-03-15 14:30"
	result := ToyyyyMMddHHmm(tm)
	if result != expected {
		t.Errorf("ToyyyyMMddHHmm() = %v, want %v", result, expected)
	}
}

func TestToyyyyMMddHHmmss(t *testing.T) {
	tm := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	expected := "2024-03-15 14:30:45"
	result := ToyyyyMMddHHmmss(tm)
	if result != expected {
		t.Errorf("ToyyyyMMddHHmmss() = %v, want %v", result, expected)
	}
}

func TestToyyyyMM2(t *testing.T) {
	tm := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	expected := "202403"
	result := ToyyyyMM2(tm)
	if result != expected {
		t.Errorf("ToyyyyMM2() = %v, want %v", result, expected)
	}
}

func TestToyyyyMMdd(t *testing.T) {
	tm := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	expected := "20240315"
	result := ToyyyyMMdd(tm)
	if result != expected {
		t.Errorf("ToyyyyMMdd() = %v, want %v", result, expected)
	}
}

func TestParseTimeStrForYYYYMMdd(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid date",
			input:   "20240315",
			want:    time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "another valid date",
			input:   "20231225",
			want:    time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "2024-03-15",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimeStrForYYYYMMdd(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeStrForYYYYMMdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseTimeStrForYYYYMMdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToyyyy_MM_dd(t *testing.T) {
	tm := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	expected := "2024-03-15"
	result := Toyyyy_MM_dd(tm)
	if result != expected {
		t.Errorf("Toyyyy_MM_dd() = %v, want %v", result, expected)
	}
}

func TestToyyyy_MM(t *testing.T) {
	tm := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	expected := "2024-03"
	result := Toyyyy_MM(tm)
	if result != expected {
		t.Errorf("Toyyyy_MM() = %v, want %v", result, expected)
	}
}

func TestParseTimeStrForYYYY_mm_dd(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid date",
			input:   "2024-03-15",
			want:    time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "another valid date",
			input:   "2023-12-25",
			want:    time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "20240315",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimeStrForYYYY_mm_dd(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeStrForYYYY_mm_dd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseTimeStrForYYYY_mm_dd() = %v, want %v", got, tt.want)
			}
		})
	}
}
