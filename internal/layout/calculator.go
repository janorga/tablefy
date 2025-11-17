package layout

// CalculateFullColumnWidths calculates the natural width for each column without truncation
func CalculateFullColumnWidths(rows [][]string) []int {
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

// CalculateColumnWidths calculates the optimal width for each column
func CalculateColumnWidths(rows [][]string, termWidth int) []int {
	if len(rows) == 0 {
		return nil
	}

	numCols := len(rows[0])
	widths := make([]int, numCols)
	minWidths := make([]int, numCols)

	// Calculate the maximum width of each column
	for _, row := range rows {
		for i := 0; i < len(row) && i < numCols; i++ {
			if len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	// Set minimum widths: smaller for short columns
	for i := range numCols {
		if i == numCols-1 {
			// Last column (typically COMMAND) - more space
			minWidths[i] = 20
		} else if widths[i] <= 5 {
			// Very small columns (PID, etc)
			minWidths[i] = widths[i]
		} else {
			// Medium columns
			minWidths[i] = 8
		}
	}

	// Calculate space needed for borders and padding
	overhead := numCols*3 + 1
	availableWidth := termWidth - overhead

	// Calculate total width of columns
	totalWidth := 0
	for _, w := range widths {
		totalWidth += w
	}

	// If everything fits, return original widths
	if totalWidth <= availableWidth {
		return widths
	}

	// If it doesn't fit, adjust with priority to the last column
	adjustedWidths := make([]int, numCols)

	// First pass: assign minimum widths
	remainingWidth := availableWidth
	for i := range adjustedWidths {
		adjustedWidths[i] = minWidths[i]
		remainingWidth -= minWidths[i]
	}

	if remainingWidth <= 0 {
		return adjustedWidths
	}

	// Second pass: give more space to the last column
	if numCols > 0 {
		// Give 60% of remaining space to the last column
		lastColExtra := int(float64(remainingWidth) * 0.6)
		adjustedWidths[numCols-1] += lastColExtra
		remainingWidth -= lastColExtra

		// Distribute the rest proportionally among the others
		if remainingWidth > 0 && numCols > 1 {
			totalOtherWidth := 0
			for i := 0; i < numCols-1; i++ {
				totalOtherWidth += widths[i]
			}

			for i := 0; i < numCols-1; i++ {
				if totalOtherWidth > 0 {
					proportion := float64(widths[i]) / float64(totalOtherWidth)
					additionalWidth := int(proportion * float64(remainingWidth))
					adjustedWidths[i] += additionalWidth
				}
			}
		}
	}

	return adjustedWidths
}

// CalculateColumnWidthsWithAutoExpand calculates column widths with auto-expand support
// When a column has focus and contains truncated cells, it gets expanded to full width
func CalculateColumnWidthsWithAutoExpand(rows [][]string, termWidth int, focusedColumn int, currentWidths []int) []int {
	if len(rows) == 0 || focusedColumn < 0 || focusedColumn >= len(currentWidths) {
		return currentWidths
	}

	// Check if the focused column has truncated cells
	if !ColumnHasTruncatedCells(rows, focusedColumn, currentWidths[focusedColumn]) {
		return currentWidths
	}

	// Get the full width needed for the focused column
	fullWidth := GetRequiredWidthForColumn(rows, focusedColumn)

	// If the current width is already sufficient, no need to expand
	if currentWidths[focusedColumn] >= fullWidth {
		return currentWidths
	}

	// Calculate how much extra space we need
	extraNeeded := fullWidth - currentWidths[focusedColumn]

	// Calculate space needed for borders and padding
	overhead := len(currentWidths)*3 + 1
	totalCurrentWidth := 0
	for _, w := range currentWidths {
		totalCurrentWidth += w
	}

	totalAvailableWidth := termWidth - overhead

	// Check if we have enough space to expand without shrinking
	if totalCurrentWidth+extraNeeded <= totalAvailableWidth {
		// We have space, just expand the focused column
		expandedWidths := make([]int, len(currentWidths))
		copy(expandedWidths, currentWidths)
		expandedWidths[focusedColumn] = fullWidth
		return expandedWidths
	}

	// We need to shrink other columns. Do it proportionally.
	expandedWidths := make([]int, len(currentWidths))
	copy(expandedWidths, currentWidths)

	// Calculate how much we need to shrink from other columns
	totalToShrink := totalCurrentWidth + extraNeeded - totalAvailableWidth

	if totalToShrink <= 0 {
		// Actually we have enough space (shouldn't get here, but just in case)
		expandedWidths[focusedColumn] = fullWidth
		return expandedWidths
	}

	// Shrink all non-focused columns proportionally
	totalNonFocusedWidth := 0
	for i, w := range expandedWidths {
		if i != focusedColumn {
			totalNonFocusedWidth += w
		}
	}

	if totalNonFocusedWidth > 0 {
		for i := range expandedWidths {
			if i != focusedColumn {
				proportion := float64(expandedWidths[i]) / float64(totalNonFocusedWidth)
				shrinkAmount := int(proportion * float64(totalToShrink))
				expandedWidths[i] -= shrinkAmount

				// Ensure minimum width for readability (at least 5 chars)
				if expandedWidths[i] < 5 {
					expandedWidths[i] = 5
				}
			}
		}
	}

	// Now expand the focused column with available space
	expandedWidths[focusedColumn] = fullWidth

	return expandedWidths
}
