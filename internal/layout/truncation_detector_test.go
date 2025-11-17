package layout

import (
	"testing"
)

func TestIsTruncated(t *testing.T) {
	tests := []struct {
		cell     string
		maxWidth int
		expected bool
	}{
		{"hello", 10, false},
		{"hello", 5, false},
		{"hello", 4, true},
		{"hello world", 5, true},
		{"a", 1, false},
		{"ab", 1, true},
	}

	for _, tt := range tests {
		result := IsTruncated(tt.cell, tt.maxWidth)
		if result != tt.expected {
			t.Errorf("IsTruncated(%q, %d) = %v, want %v", tt.cell, tt.maxWidth, result, tt.expected)
		}
	}
}

func TestColumnHasTruncatedCells(t *testing.T) {
	rows := [][]string{
		{"NAME", "NAMESPACE", "STATUS"},
		{"argocd", "argocd", "deployed"},
		{"very-long-name-that-will-be-truncated", "kube-system", "pending"},
	}

	// Column 0 with width 10 - should have truncated cells
	if !ColumnHasTruncatedCells(rows, 0, 10) {
		t.Error("ColumnHasTruncatedCells should return true for column 0 with width 10")
	}

	// Column 1 with width 20 - should not have truncated cells
	if ColumnHasTruncatedCells(rows, 1, 20) {
		t.Error("ColumnHasTruncatedCells should return false for column 1 with width 20")
	}

	// Column 0 with width 50 - should not have truncated cells
	if ColumnHasTruncatedCells(rows, 0, 50) {
		t.Error("ColumnHasTruncatedCells should return false for column 0 with width 50")
	}
}

func TestGetRequiredWidthForColumn(t *testing.T) {
	rows := [][]string{
		{"NAME", "NAMESPACE", "STATUS"},
		{"argocd", "argocd", "deployed"},
		{"longhorn", "longhorn-system", "pending"},
	}

	// Test column 0
	width := GetRequiredWidthForColumn(rows, 0)
	if width != len("longhorn") {
		t.Errorf("GetRequiredWidthForColumn for column 0 = %d, want %d", width, len("longhorn"))
	}

	// Test column 1
	width = GetRequiredWidthForColumn(rows, 1)
	if width != len("longhorn-system") {
		t.Errorf("GetRequiredWidthForColumn for column 1 = %d, want %d", width, len("longhorn-system"))
	}

	// Test column 2
	width = GetRequiredWidthForColumn(rows, 2)
	if width != len("deployed") {
		t.Errorf("GetRequiredWidthForColumn for column 2 = %d, want %d", width, len("deployed"))
	}
}
