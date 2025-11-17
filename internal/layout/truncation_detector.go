package layout

// IsTruncated checks if a cell has been truncated
func IsTruncated(cell string, maxWidth int) bool {
	if maxWidth <= 3 {
		return len(cell) > maxWidth
	}
	return len(cell) > maxWidth
}

// ColumnHasTruncatedCells checks if any cell in a column is truncated
func ColumnHasTruncatedCells(rows [][]string, columnIndex int, width int) bool {
	for _, row := range rows {
		if columnIndex < len(row) {
			if IsTruncated(row[columnIndex], width) {
				return true
			}
		}
	}
	return false
}

// ColumnFullWidths returns the natural widths for all columns without truncation
func ColumnFullWidths(rows [][]string) []int {
	if len(rows) == 0 {
		return nil
	}

	numCols := len(rows[0])
	widths := make([]int, numCols)

	// Calculate the maximum width of each column
	for _, row := range rows {
		for i := 0; i < len(row) && i < numCols; i++ {
			if len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	return widths
}

// GetRequiredWidthForColumn returns the minimum width needed to show all content in a column without truncation
func GetRequiredWidthForColumn(rows [][]string, columnIndex int) int {
	maxWidth := 0
	for _, row := range rows {
		if columnIndex < len(row) && len(row[columnIndex]) > maxWidth {
			maxWidth = len(row[columnIndex])
		}
	}
	return maxWidth
}
