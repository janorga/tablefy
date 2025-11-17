package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"tablefy/internal/layout"
	"tablefy/internal/model"
)

// RenderZoomView renders selected columns in full width
func RenderZoomView(m model.Model) string {
	if len(m.Rows) == 0 || len(m.SelectedColumns) == 0 {
		return "No data to display"
	}

	// Determine which rows to use based on active filter
	rowsToUse := m.Rows
	if len(m.FilteredRowIndices) > 0 {
		rowsToUse = GetFilteredRows(m.Rows, m.FilteredRowIndices)
	}

	// Get sorted list of selected column indices
	selectedIndices := sortSelectedColumns(m.SelectedColumns)

	// Extract selected columns
	zoomedRows := extractSelectedColumns(rowsToUse, selectedIndices)

	// Calculate visible rows based on terminal height (account for title and help)
	visibleRows := layout.GetVisibleRowsForZoom(m.TermHeight)

	// Apply scroll offset to get visible subset of rows
	displayRows := applyScrollOffset(zoomedRows, m.ScrollOffset, visibleRows)

	// Calculate optimal widths for zoomed table
	widths := calculateZoomWidths(zoomedRows, m.TermWidth)

	// Truncate only the display rows according to widths
	truncatedRows := layout.TruncateRows(displayRows, widths)

	// Create table
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(0, 1)
		})

	// Add header and rows
	if len(truncatedRows) > 0 {
		t.Headers(truncatedRows[0]...)
		for i := 1; i < len(truncatedRows); i++ {
			// Debug: Ensure all rows have the same number of columns
			if len(truncatedRows[i]) != len(truncatedRows[0]) {
				// Pad or truncate to match header length
				row := make([]string, len(truncatedRows[0]))
				copy(row, truncatedRows[i])
				t.Row(row...)
			} else {
				t.Row(truncatedRows[i]...)
			}
		}
	}

	// Build title
	title := buildZoomTitle(m, selectedIndices)

	// Build filter indicator if active
	output := t.Render()
	if len(m.FilteredRowIndices) > 0 {
		filterIndicator := buildFilterIndicator(m)
		output = filterIndicator + "\n" + output
	}

	// Build help text with scroll indicator
	helpText := buildZoomViewHelp(m, zoomedRows, visibleRows)
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(helpText)

	return fmt.Sprintf("%s\n\n%s\n%s", title, output, help)
}

// buildZoomTitle builds the title for zoom view
func buildZoomTitle(m model.Model, selectedIndices []int) string {
	var columnNames []string
	for _, colIdx := range selectedIndices {
		if colIdx < len(m.Rows[0]) {
			columnNames = append(columnNames, m.Rows[0][colIdx])
		}
	}

	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9D4EDD")).
		Render(fmt.Sprintf("Zoomed: %s", strings.Join(columnNames, ", ")))
}

// buildZoomViewHelp builds the help text for zoom view
func buildZoomViewHelp(m model.Model, zoomedRows [][]string, visibleRows int) string {
	totalDataRows := len(zoomedRows) - 1
	scrollInfo := ""
	if totalDataRows > visibleRows {
		currentPos := m.ScrollOffset + 1
		maxPos := totalDataRows - visibleRows + 1
		scrollInfo = fmt.Sprintf(" | ↑↓/jk/PgUp/PgDn: Scroll (%d/%d)", currentPos, maxPos)
	}
	return fmt.Sprintf("q: Exit zoom%s", scrollInfo)
}
