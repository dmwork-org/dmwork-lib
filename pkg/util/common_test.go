package util

import (
	"errors"
	"testing"
)

func TestMustNoErr_NilError(t *testing.T) {
	// Should not panic when err is nil
	MustNoErr(nil)
}

func TestMustNoErr_NonNilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustNoErr should panic on non-nil error")
		}
	}()
	MustNoErr(errors.New("test error"))
}

func TestSubstr(t *testing.T) {
	tests := []struct {
		str    string
		start  int
		length int
		want   string
	}{
		{"hello", 0, 2, "he"},
		{"hello", 1, 3, "ell"},
		{"hello", 0, 0, ""},
		{"hello", -2, 2, "lo"},
		{"hello", 0, 10, "hello"},
	}

	for _, tt := range tests {
		got := Substr(tt.str, tt.start, tt.length)
		if got != tt.want {
			t.Errorf("Substr(%q, %d, %d) = %q, want %q", tt.str, tt.start, tt.length, got, tt.want)
		}
	}
}

func TestObjToStr(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"int", int(123), "123"},
		{"uint", uint(456), "456"},
		{"int64", int64(789), "789"},
		{"uint64", uint64(1000), "1000"},
		{"int8", int8(10), "10"},
		{"uint8", uint8(20), "20"},
		{"int16", int16(30), "30"},
		{"uint16", uint16(40), "40"},
		{"int32", int32(50), "50"},
		{"uint32", uint32(60), "60"},
		{"string", "hello", "hello"},
		{"float32", float32(3.14), "3.14"},
		{"float64", float64(2.718), "2.718"},
		{"float64_scientific", float64(0.0000001), "1e-07"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objToStr(tt.input)
			if result != tt.expected {
				t.Errorf("objToStr(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMapToQueryParamSort(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]interface{}
		expected string
	}{
		{
			name:     "empty map",
			params:   map[string]interface{}{},
			expected: "",
		},
		{
			name: "mixed types sorted",
			params: map[string]interface{}{
				"b": "value",
				"a": 123,
				"c": uint32(456),
			},
			expected: "a=123&b=value&c=456",
		},
		{
			name: "float values",
			params: map[string]interface{}{
				"price": float64(19.99),
			},
			expected: "price=19.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToQueryParamSort(tt.params)
			if result != tt.expected {
				t.Errorf("MapToQueryParamSort(%v) = %q, want %q", tt.params, result, tt.expected)
			}
		})
	}
}
