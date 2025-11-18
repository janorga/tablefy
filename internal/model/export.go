package model

import (
	"fmt"
	"strings"

	"tablefy/internal/layout"
)

// GetExportData returns the currently visible table data (excluding header)
// Data is formatted with aligned columns but without borders
// In NormalView: exports all visible columns
// In FilterView: exports filtered rows of all columns
// In ZoomView: exports selected columns
func (m *Model) GetExportData() string {
	if len(m.Rows) == 0 {
		return ""
	}

	// Determine which rows to use based on filter status
	var rowIndices []int
	if len(m.FilteredRowIndices) > 0 {
		// Use filtered rows
		rowIndices = m.FilteredRowIndices
	} else {
		// Use all rows except header (row 0)
		for i := 1; i < len(m.Rows); i++ {
			rowIndices = append(rowIndices, i)
		}
	}

	if len(rowIndices) == 0 {
		return ""
	}

	// Determine which columns to include based on view mode
	var colIndices []int
	if m.ViewMode == ZoomView && len(m.SelectedColumns) > 0 {
		// In zoom view, only include selected columns (sorted)
		for col := 0; col < len(m.Rows[0]); col++ {
			if m.SelectedColumns[col] {
				colIndices = append(colIndices, col)
			}
		}
	} else {
		// Include all columns
		for col := 0; col < len(m.Rows[0]); col++ {
			colIndices = append(colIndices, col)
		}
	}

	if len(colIndices) == 0 {
		return ""
	}

	// Build rows to export with only selected columns
	var rowsToExport [][]string
	for _, rowIdx := range rowIndices {
		if rowIdx >= 0 && rowIdx < len(m.Rows) {
			var rowCols []string
			for _, colIdx := range colIndices {
				if colIdx < len(m.Rows[rowIdx]) {
					rowCols = append(rowCols, m.Rows[rowIdx][colIdx])
				} else {
					rowCols = append(rowCols, "")
				}
			}
			rowsToExport = append(rowsToExport, rowCols)
		}
	}

	// Calculate optimal column widths based on the data
	widths := layout.CalculateFullColumnWidths(rowsToExport)

	// Format rows with proper alignment (no borders)
	var lines []string
	for _, row := range rowsToExport {
		var paddedCols []string
		for col := 0; col < len(row); col++ {
			value := row[col]

			// Get the width for this column
			width := 0
			if col < len(widths) {
				width = widths[col]
			}

			// Pad the value to the column width
			paddedValue := fmt.Sprintf("%-*s", width, value)
			paddedCols = append(paddedCols, paddedValue)
		}

		// Join columns with spaces for alignment
		lines = append(lines, strings.Join(paddedCols, "  "))
	}

	return strings.Join(lines, "\n")
}
