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
