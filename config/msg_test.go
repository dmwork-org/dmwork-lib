package config

import (
	"fmt"
	"testing"
)

func TestSetting(t *testing.T) {
	setting := SettingFromUint8(160)
	fmt.Println(setting.Signal)
}

func TestMessageResp_GetContentType(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		expected int
	}{
		{
			name:     "valid type field",
			payload:  []byte(`{"type": 1}`),
			expected: 1,
		},
		{
			name:     "missing type field",
			payload:  []byte(`{"content": "hello"}`),
			expected: 0,
		},
		{
			name:     "type is string instead of number",
			payload:  []byte(`{"type": "text"}`),
			expected: 0,
		},
		{
			name:     "type is null",
			payload:  []byte(`{"type": null}`),
			expected: 0,
		},
		{
			name:     "empty payload",
			payload:  []byte(`{}`),
			expected: 0,
		},
		{
			name:     "invalid json payload",
			payload:  []byte(`invalid`),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageResp{
				Payload: tt.payload,
			}
			result := m.GetContentType()
			if result != tt.expected {
				t.Errorf("GetContentType() = %d, expected %d", result, tt.expected)
			}
		})
	}
}
