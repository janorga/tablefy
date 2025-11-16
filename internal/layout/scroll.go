package layout

// GetMaxScroll calculates the maximum scroll offset for a table
func GetMaxScroll(numRows, termHeight int) int {
	// Account for header, borders, and help text (approximately 6 lines)
	visibleRows := termHeight - 6
	if visibleRows < 1 {
		visibleRows = 1
	}

	// Number of data rows (excluding header)
	dataRows := numRows - 1

	maxScroll := dataRows - visibleRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

// GetVisibleRows calculates how many rows are visible in the terminal
func GetVisibleRows(termHeight int) int {
	visibleRows := termHeight - 6 // Account for header, borders, help text
	if visibleRows < 1 {
		visibleRows = 1
	}
	return visibleRows
}

// GetVisibleRowsForZoom calculates how many rows are visible in zoom mode
func GetVisibleRowsForZoom(termHeight int) int {
	visibleRows := termHeight - 8 // Account for title, header, borders, help text
	if visibleRows < 1 {
		visibleRows = 1
	}
	return visibleRows
}
