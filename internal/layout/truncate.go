package layout

// TruncateCell truncates a cell to the maximum width
func TruncateCell(cell string, maxWidth int) string {
	if len(cell) <= maxWidth {
		return cell
	}
	if maxWidth <= 3 {
		return cell[:maxWidth]
	}
	return cell[:maxWidth-3] + "..."
}

// TruncateRows truncates rows according to column widths
func TruncateRows(rows [][]string, widths []int) [][]string {
	truncated := make([][]string, len(rows))

	for i, row := range rows {
		truncated[i] = make([]string, len(row))
		for j, cell := range row {
			if j < len(widths) {
				truncated[i][j] = TruncateCell(cell, widths[j])
			} else {
				truncated[i][j] = cell
			}
		}
	}

	return truncated
}
