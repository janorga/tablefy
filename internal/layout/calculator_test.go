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

// TestCalculateColumnWidths_FillsAvailableSpace verifies that CalculateColumnWidths
// distributes available space proportionally even when content is shorter than terminal width
func TestCalculateColumnWidths_FillsAvailableSpace(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS", "AGE"},
		{"app1", "running", "1d"},
		{"app2", "deployed", "5d"},
	}

	termWidth := 80
	numCols := 3
	overhead := numCols*3 + 1 // borders and padding
	availableWidth := termWidth - overhead

	result := CalculateColumnWidths(rows, termWidth)

	// Calculate total width
	totalWidth := 0
	for _, w := range result {
		totalWidth += w
	}

	// Should use all available space (or very close, accounting for rounding)
	if totalWidth != availableWidth {
		t.Errorf("Total width should be %d, got %d (difference: %d)", availableWidth, totalWidth, availableWidth-totalWidth)
	}
}

// TestCalculateColumnWidths_ProportionalDistribution verifies that space is distributed
// proportionally based on natural column widths
func TestCalculateColumnWidths_ProportionalDistribution(t *testing.T) {
	rows := [][]string{
		{"SHORT", "MEDIUM", "VERYLONGNAME"},
		{"a", "bb", "cccccccccccc"},
	}

	termWidth := 100
	numCols := 3
	overhead := numCols*3 + 1
	availableWidth := termWidth - overhead

	result := CalculateColumnWidths(rows, termWidth)

	// The longest column should get more space
	if result[2] <= result[0] || result[2] <= result[1] {
		t.Errorf("Longest column should get more space: %v", result)
	}

	// Total should fill available space
	totalWidth := 0
	for _, w := range result {
		totalWidth += w
	}
	if totalWidth != availableWidth {
		t.Errorf("Should fill available space: %d != %d", totalWidth, availableWidth)
	}
}
