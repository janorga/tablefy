package view

import (
	"tablefy/internal/layout"
	"tablefy/internal/model"
)

// Render renders the UI based on the current view mode
func Render(m model.Model) string {
	if m.ViewMode == model.FilterView {
		return RenderFilterView(m)
	}
	if m.ViewMode == model.ZoomView {
		return RenderZoomView(m)
	}
	return RenderNormalView(m)
}

// applyScrollOffset applies scroll offset to rows
func applyScrollOffset(rows [][]string, scrollOffset, visibleRows int) [][]string {
	startRow := 1 + scrollOffset // +1 to skip header
	endRow := min(startRow+visibleRows, len(rows))

	// Build rows to display (header + visible data rows)
	var displayRows [][]string
	displayRows = append(displayRows, rows[0]) // Always include header
	if startRow < len(rows) {
		displayRows = append(displayRows, rows[startRow:endRow]...)
	}

	return displayRows
}

// GetFilteredRows extracts rows at specified indices (wrapper for model function)
func GetFilteredRows(rows [][]string, filteredIndices []int) [][]string {
	return model.GetFilteredRows(rows, filteredIndices)
}

// sortSelectedColumns sorts the selected column indices
func sortSelectedColumns(selectedColumns map[int]bool) []int {
	var selectedIndices []int
	for colIdx := range selectedColumns {
		selectedIndices = append(selectedIndices, colIdx)
	}
	// Sort them to maintain column order
	for i := 0; i < len(selectedIndices); i++ {
		for j := i + 1; j < len(selectedIndices); j++ {
			if selectedIndices[i] > selectedIndices[j] {
				selectedIndices[i], selectedIndices[j] = selectedIndices[j], selectedIndices[i]
			}
		}
	}
	return selectedIndices
}

// extractSelectedColumns extracts only the selected columns from rows
func extractSelectedColumns(rows [][]string, selectedIndices []int) [][]string {
	numSelectedCols := len(selectedIndices)
	zoomedRows := make([][]string, len(rows))
	for i, row := range rows {
		zoomedRows[i] = make([]string, numSelectedCols)
		for j, colIdx := range selectedIndices {
			if colIdx < len(row) {
				zoomedRows[i][j] = row[colIdx]
			} else {
				zoomedRows[i][j] = ""
			}
		}
	}
	return zoomedRows
}

// calculateZoomWidths calculates optimal widths for zoomed table
func calculateZoomWidths(zoomedRows [][]string, termWidth int) []int {
	// First, try to use full widths without truncation
	widths := layout.CalculateFullColumnWidths(zoomedRows)

	// Check if the table fits in the terminal
	overhead := len(widths)*3 + 1 // borders and padding
	totalWidth := overhead
	for _, w := range widths {
		totalWidth += w
	}

	// If it doesn't fit, recalculate with terminal width constraint
	if totalWidth > termWidth {
		widths = layout.CalculateColumnWidths(zoomedRows, termWidth)
	}

	return widths
}
