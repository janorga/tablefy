package view

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"tablefy/internal/layout"
	"tablefy/internal/model"
)

// RenderNormalView renders the table with all columns
func RenderNormalView(m model.Model) string {
	if len(m.Rows) == 0 {
		return "No data to display"
	}

	// Determine which rows to display based on active filter
	rowsToDisplay := m.Rows
	if len(m.FilteredRowIndices) > 0 {
		rowsToDisplay = GetFilteredRows(m.Rows, m.FilteredRowIndices)
	}

	// Calculate visible rows based on terminal height
	visibleRows := layout.GetVisibleRows(m.TermHeight)

	// Apply scroll offset to get visible subset of rows
	displayRows := applyScrollOffset(rowsToDisplay, m.ScrollOffset, visibleRows)

	// Calculate optimal widths
	// Base widths on the rows that will be displayed (filtered or unfiltered)
	widths := layout.CalculateColumnWidths(rowsToDisplay, m.TermWidth)

	// Apply auto-expand if enabled
	if m.AutoExpand {
		widths = layout.CalculateColumnWidthsWithAutoExpand(rowsToDisplay, m.TermWidth, m.CurrentColumn, widths)
	}

	// Truncate rows according to widths
	truncatedRows := layout.TruncateRows(displayRows, widths)

	// Create a new table
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)

			// Highlight current column
			if col == m.CurrentColumn {
				style = style.Background(lipgloss.Color("#3D3D3D"))
			}

			// Mark selected columns with different background
			if m.SelectedColumns[col] {
				style = style.Background(lipgloss.Color("#5A4E8C"))
			}

			return style.Foreground(lipgloss.Color("252"))
		})

	// Add all rows
	t.Headers(truncatedRows[0]...)
	for i := 1; i < len(truncatedRows); i++ {
		t.Row(truncatedRows[i]...)
	}

	// Build filter indicator if active
	output := t.Render()
	if len(m.FilteredRowIndices) > 0 {
		filterIndicator := buildFilterIndicator(m)
		output = filterIndicator + "\n" + output
	}

	// Build help text with scroll indicator
	helpText := buildNormalViewHelp(m, rowsToDisplay, visibleRows)
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(helpText)

	return output + "\n" + help
}

// buildFilterIndicator builds the filter status indicator
func buildFilterIndicator(m model.Model) string {
	columnName := ""
	if m.FilterColumnIndex < len(m.Rows[0]) {
		columnName = m.Rows[0][m.FilterColumnIndex]
	}

	filterText := fmt.Sprintf("🔍 Filter active: [%s] = \"%s\" (%d results)", columnName, m.FilterInput, len(m.FilteredRowIndices))

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")). // Bright yellow
		Bold(true).
		Padding(0, 1).
		Render(filterText)
}

// buildNormalViewHelp builds the help text for normal view
func buildNormalViewHelp(m model.Model, rowsToDisplay [][]string, visibleRows int) string {
	selectedCount := len(m.SelectedColumns)
	totalDataRows := len(rowsToDisplay) - 1 // Exclude header
	scrollInfo := ""
	if totalDataRows > visibleRows {
		currentPos := m.ScrollOffset + 1
		maxPos := totalDataRows - visibleRows + 1
		scrollInfo = fmt.Sprintf(" | ↑↓/jk/PgUp/PgDn: Scroll (%d/%d)", currentPos, maxPos)
	}

	autoExpandInfo := ""
	if m.AutoExpand {
		autoExpandInfo = " | [AUTO-EXPAND ON]"
	}

	filterInfo := ""
	if len(m.FilteredRowIndices) > 0 {
		filterInfo = fmt.Sprintf(" | [FILTERED: %d/%d rows]", totalDataRows, len(m.Rows)-1)
	}

	return fmt.Sprintf("← → / h l: Navigate | s: Toggle select (%d selected) | Enter: Zoom | f: Filter%s%s%s | q: Quit", selectedCount, scrollInfo, autoExpandInfo, filterInfo)
}
