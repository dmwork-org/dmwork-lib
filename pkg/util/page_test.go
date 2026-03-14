package util

import "testing"

func TestToPageNumOrDefault(t *testing.T) {
	tests := []struct {
		name          string
		pageIndex     string
		pageSize      string
		wantPageIndex uint64
		wantPageSize  uint64
	}{
		{
			name:          "valid values",
			pageIndex:     "5",
			pageSize:      "20",
			wantPageIndex: 5,
			wantPageSize:  20,
		},
		{
			name:          "empty strings use defaults",
			pageIndex:     "",
			pageSize:      "",
			wantPageIndex: 1,
			wantPageSize:  10,
		},
		{
			name:          "invalid pageIndex uses default",
			pageIndex:     "invalid",
			pageSize:      "20",
			wantPageIndex: 1,
			wantPageSize:  20,
		},
		{
			name:          "invalid pageSize uses default",
			pageIndex:     "5",
			pageSize:      "invalid",
			wantPageIndex: 5,
			wantPageSize:  10,
		},
		{
			name:          "both invalid use defaults",
			pageIndex:     "abc",
			pageSize:      "xyz",
			wantPageIndex: 1,
			wantPageSize:  10,
		},
		{
			name:          "zero pageIndex uses default",
			pageIndex:     "0",
			pageSize:      "20",
			wantPageIndex: 1,
			wantPageSize:  20,
		},
		{
			name:          "zero pageSize uses default",
			pageIndex:     "5",
			pageSize:      "0",
			wantPageIndex: 5,
			wantPageSize:  10,
		},
		{
			name:          "negative values use defaults",
			pageIndex:     "-1",
			pageSize:      "-5",
			wantPageIndex: 1,
			wantPageSize:  10,
		},
		{
			name:          "float values use defaults",
			pageIndex:     "1.5",
			pageSize:      "10.5",
			wantPageIndex: 1,
			wantPageSize:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPageIndex, gotPageSize := ToPageNumOrDefault(tt.pageIndex, tt.pageSize)
			if gotPageIndex != tt.wantPageIndex {
				t.Errorf("ToPageNumOrDefault() pageIndex = %v, want %v", gotPageIndex, tt.wantPageIndex)
			}
			if gotPageSize != tt.wantPageSize {
				t.Errorf("ToPageNumOrDefault() pageSize = %v, want %v", gotPageSize, tt.wantPageSize)
			}
		})
	}
}
