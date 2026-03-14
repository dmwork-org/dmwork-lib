package util

import (
	"testing"
)

func TestToPageNumOrDefault(t *testing.T) {
	tests := []struct {
		name          string
		pageIndex     string
		pageSize      string
		wantPageIndex uint64
		wantPageSize  uint64
	}{
		{
			name:          "empty strings use defaults",
			pageIndex:     "",
			pageSize:      "",
			wantPageIndex: 1,
			wantPageSize:  10,
		},
		{
			name:          "valid numbers",
			pageIndex:     "5",
			pageSize:      "20",
			wantPageIndex: 5,
			wantPageSize:  20,
		},
		{
			name:          "invalid pageIndex uses default",
			pageIndex:     "abc",
			pageSize:      "20",
			wantPageIndex: 1,
			wantPageSize:  20,
		},
		{
			name:          "invalid pageSize uses default",
			pageIndex:     "5",
			pageSize:      "xyz",
			wantPageIndex: 5,
			wantPageSize:  10,
		},
		{
			name:          "both invalid use defaults",
			pageIndex:     "foo",
			pageSize:      "bar",
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
			name:          "negative numbers use defaults",
			pageIndex:     "-1",
			pageSize:      "-10",
			wantPageIndex: 1,
			wantPageSize:  10,
		},
		{
			name:          "whitespace uses defaults",
			pageIndex:     " ",
			pageSize:      " ",
			wantPageIndex: 1,
			wantPageSize:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotSize := ToPageNumOrDefault(tt.pageIndex, tt.pageSize)
			if gotIndex != tt.wantPageIndex {
				t.Errorf("ToPageNumOrDefault() pageIndex = %v, want %v", gotIndex, tt.wantPageIndex)
			}
			if gotSize != tt.wantPageSize {
				t.Errorf("ToPageNumOrDefault() pageSize = %v, want %v", gotSize, tt.wantPageSize)
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	data := []string{"a", "b", "c"}
	page := NewPage(1, 10, 100, data)

	if page.PageIndex != 1 {
		t.Errorf("NewPage() PageIndex = %v, want %v", page.PageIndex, 1)
	}
	if page.PageSize != 10 {
		t.Errorf("NewPage() PageSize = %v, want %v", page.PageSize, 10)
	}
	if page.Total != 100 {
		t.Errorf("NewPage() Total = %v, want %v", page.Total, 100)
	}
	if page.Data == nil {
		t.Error("NewPage() Data should not be nil")
	}
}
