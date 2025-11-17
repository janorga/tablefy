package layout

import (
	"testing"
)

func TestCalculateColumnWidthsWithAutoExpand_NoTruncation(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS", "AGE"},
		{"argocd", "deployed", "30d"},
		{"longhorn", "running", "15d"},
	}

	currentWidths := []int{10, 10, 10}
	focusedColumn := 0
	termWidth := 50

	result := CalculateColumnWidthsWithAutoExpand(rows, termWidth, focusedColumn, currentWidths)

	// Should not change widths if there's no truncation
	if result[0] != 10 {
		t.Errorf("Width[0] should remain 10, got %d", result[0])
	}
}

func TestCalculateColumnWidthsWithAutoExpand_WithTruncation(t *testing.T) {
	rows := [][]string{
		{"NAME", "DESCRIPTION", "STATUS"},
		{"app1", "Very long description that will be truncated", "deployed"},
		{"app2", "Another long one", "running"},
	}

	currentWidths := []int{10, 15, 10}
	focusedColumn := 1
	termWidth := 80

	result := CalculateColumnWidthsWithAutoExpand(rows, termWidth, focusedColumn, currentWidths)

	// The focused column should expand to show full content
	requiredWidth := GetRequiredWidthForColumn(rows, 1)
	if result[1] < requiredWidth {
		t.Errorf("Focused column width should expand to at least %d, got %d", requiredWidth, result[1])
	}
}

func TestCalculateColumnWidthsWithAutoExpand_InvalidColumn(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS"},
		{"app1", "running"},
	}

	currentWidths := []int{10, 10}
	focusedColumn := -1
	termWidth := 50

	result := CalculateColumnWidthsWithAutoExpand(rows, termWidth, focusedColumn, currentWidths)

	// Should return unchanged widths for invalid column
	if result[0] != 10 || result[1] != 10 {
		t.Errorf("Widths should remain unchanged for invalid column")
	}
}

func TestCalculateColumnWidthsWithAutoExpand_ColumnOutOfBounds(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS"},
		{"app1", "running"},
	}

	currentWidths := []int{10, 10}
	focusedColumn := 5
	termWidth := 50

	result := CalculateColumnWidthsWithAutoExpand(rows, termWidth, focusedColumn, currentWidths)

	// Should return unchanged widths for out of bounds column
	if result[0] != 10 || result[1] != 10 {
		t.Errorf("Widths should remain unchanged for out of bounds column")
	}
}

func TestCalculateFullColumnWidths(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS", "VERSION"},
		{"argocd", "deployed", "v3.0.12"},
		{"longhorn", "running", "v1.9.1"},
	}

	widths := CalculateFullColumnWidths(rows)

	expected := []int{len("longhorn"), len("deployed"), len("v3.0.12")}
	for i, w := range widths {
		if w != expected[i] {
			t.Errorf("Width[%d] = %d, want %d", i, w, expected[i])
		}
	}
}
