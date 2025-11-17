package view

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"tablefy/internal/layout"
	"tablefy/internal/model"
)

// RenderFilterView renders the table with a fuzzy filter input overlay
func RenderFilterView(m model.Model) string {
	if len(m.Rows) == 0 {
		return "No data to display"
	}

	// Get the header for display
	columnName := ""
	if m.FilterColumnIndex < len(m.Rows[0]) {
		columnName = m.Rows[0][m.FilterColumnIndex]
	}

	// Build filtered rows to display
	filteredRows := GetFilteredRows(m.Rows, m.FilteredRowIndices)

	// Calculate visible rows based on terminal height
	visibleRows := layout.GetVisibleRows(m.TermHeight)

	// Apply scroll offset to filtered rows
	displayRows := applyScrollOffset(filteredRows, m.FilterScrollOffset, visibleRows)

	// Calculate optimal widths based on ALL rows (not just filtered)
	widths := layout.CalculateColumnWidths(m.Rows, m.TermWidth)

	// Truncate rows according to widths
	truncatedRows := layout.TruncateRows(displayRows, widths)

	// Create the table
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)

			// Highlight current filter column with soft purple
			if col == m.FilterColumnIndex {
				style = style.Background(lipgloss.Color("#8B7BA8"))
			}

			return style.Foreground(lipgloss.Color("252"))
		})

	// Add all rows
	t.Headers(truncatedRows[0]...)
	for i := 1; i < len(truncatedRows); i++ {
		t.Row(truncatedRows[i]...)
	}

	// Build filter input display
	matchCount := len(m.FilteredRowIndices)
	filterInput := fmt.Sprintf("Filter [%s]: %s (%d matches)", columnName, m.FilterInput, matchCount)

	filterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Bold(true).
		Padding(0, 1)

	filterDisplay := filterStyle.Render(filterInput)

	// Help text
	helpText := "Type to search | Esc: Cancel | Enter: Apply"
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(helpText)

	return t.Render() + "\n" + filterDisplay + "\n" + help
}
