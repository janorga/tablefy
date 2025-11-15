package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

// ViewMode represents the current view mode
type ViewMode int

const (
	NormalView ViewMode = iota
	ZoomView
)

// model represents the application state
type model struct {
	rows            [][]string
	currentColumn   int
	selectedColumns map[int]bool
	viewMode        ViewMode
	termWidth       int
	termHeight      int
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "q":
			if m.viewMode == ZoomView {
				// Exit zoom mode
				m.viewMode = NormalView
				return m, nil
			}
			return m, tea.Quit
		case "left", "h":
			if m.viewMode == NormalView && m.currentColumn > 0 {
				m.currentColumn--
			}
		case "right", "l":
			if m.viewMode == NormalView && m.currentColumn < len(m.rows[0])-1 {
				m.currentColumn++
			}
		case "s", "S":
			// Toggle selection of current column
			if m.viewMode == NormalView {
				if m.selectedColumns[m.currentColumn] {
					delete(m.selectedColumns, m.currentColumn)
				} else {
					m.selectedColumns[m.currentColumn] = true
				}
			}
		case "enter", " ":
			if m.viewMode == NormalView && len(m.selectedColumns) > 0 {
				// Enter zoom mode with selected columns
				m.viewMode = ZoomView
			}
		}
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	if m.viewMode == ZoomView {
		return m.renderZoomView()
	}
	return m.renderNormalView()
}

// renderNormalView renders the table with all columns
func (m model) renderNormalView() string {
	if len(m.rows) == 0 {
		return "No data to display"
	}

	// Calculate optimal widths
	widths := calculateColumnWidths(m.rows, m.termWidth)

	// Truncate rows according to widths
	truncatedRows := truncateRows(m.rows, widths)

	// Create a new table
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)

			// Highlight current column
			if col == m.currentColumn {
				style = style.Background(lipgloss.Color("#3D3D3D"))
			}

			// Mark selected columns with different background
			if m.selectedColumns[col] {
				style = style.Background(lipgloss.Color("#5A4E8C"))
			}

			return style.Foreground(lipgloss.Color("252"))
		})

	// Add all rows
	t.Headers(truncatedRows[0]...)
	for i := 1; i < len(truncatedRows); i++ {
		t.Row(truncatedRows[i]...)
	}

	selectedCount := len(m.selectedColumns)
	helpText := fmt.Sprintf("\n← → / h l: Navigate | s: Toggle select (%d selected) | Enter: Zoom | q: Quit", selectedCount)
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(helpText)

	return t.Render() + help
}

// renderZoomView renders selected columns in full width
func (m model) renderZoomView() string {
	if len(m.rows) == 0 || len(m.selectedColumns) == 0 {
		return "No data to display"
	}

	// Get sorted list of selected column indices
	var selectedIndices []int
	for colIdx := range m.selectedColumns {
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

	// Extract selected columns
	zoomedRows := make([][]string, len(m.rows))
	for i, row := range m.rows {
		zoomedRows[i] = make([]string, 0, len(selectedIndices))
		for _, colIdx := range selectedIndices {
			if colIdx < len(row) {
				zoomedRows[i] = append(zoomedRows[i], row[colIdx])
			} else {
				zoomedRows[i] = append(zoomedRows[i], "")
			}
		}
	}

	// Calculate optimal widths for zoomed table
	widths := calculateColumnWidths(zoomedRows, m.termWidth)

	// Truncate rows according to widths
	truncatedRows := truncateRows(zoomedRows, widths)

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
	t.Headers(truncatedRows[0]...)
	for i := 1; i < len(truncatedRows); i++ {
		t.Row(truncatedRows[i]...)
	}

	// Build column names for title
	var columnNames []string
	for _, colIdx := range selectedIndices {
		if colIdx < len(m.rows[0]) {
			columnNames = append(columnNames, m.rows[0][colIdx])
		}
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9D4EDD")).
		Render(fmt.Sprintf("Zoomed: %s\n\n", strings.Join(columnNames, ", ")))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\nq: Exit zoom")

	return title + t.Render() + help
}

// parseTable parses the input and converts it to rows and columns
func parseTable(input string) [][]string {
	lines := strings.Split(input, "\n")
	var validLines []string

	// Filter empty lines
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			validLines = append(validLines, line)
		}
	}

	if len(validLines) == 0 {
		return nil
	}

	// Get the header and count how many columns it has
	header := validLines[0]
	headerFields := strings.Fields(header)
	numHeaderCols := len(headerFields)

	if numHeaderCols == 0 {
		return nil
	}

	var rows [][]string
	rows = append(rows, headerFields)

	// Process each data line
	for i := 1; i < len(validLines); i++ {
		fields := strings.Fields(validLines[i])

		if len(fields) == 0 {
			continue
		}

		var row []string

		if len(fields) >= numHeaderCols {
			// If it has at least as many fields as columns in the header,
			// take the first numHeaderCols-1 and join the rest in the last one
			row = append(row, fields[:numHeaderCols-1]...)
			row = append(row, strings.Join(fields[numHeaderCols-1:], " "))
		} else {
			// If it has fewer fields, copy what's there and fill with empty strings
			row = make([]string, numHeaderCols)
			copy(row, fields)
		}

		rows = append(rows, row)
	}

	return rows
}

// getTerminalWidth gets the terminal width
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // Default fallback
	}
	return width
}

// truncateCell truncates a cell to the maximum width
func truncateCell(cell string, maxWidth int) string {
	if len(cell) <= maxWidth {
		return cell
	}
	if maxWidth <= 3 {
		return cell[:maxWidth]
	}
	return cell[:maxWidth-3] + "..."
}

// calculateColumnWidths calculates the optimal width for each column
func calculateColumnWidths(rows [][]string, termWidth int) []int {
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

// truncateRows truncates rows according to column widths
func truncateRows(rows [][]string, widths []int) [][]string {
	truncated := make([][]string, len(rows))

	for i, row := range rows {
		truncated[i] = make([]string, len(row))
		for j, cell := range row {
			if j < len(widths) {
				truncated[i][j] = truncateCell(cell, widths[j])
			} else {
				truncated[i][j] = cell
			}
		}
	}

	return truncated
}

// formatTable formats rows using lipgloss
func formatTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	// Get terminal width
	termWidth := getTerminalWidth()

	// Calculate optimal widths
	widths := calculateColumnWidths(rows, termWidth)

	// Truncate rows according to widths
	truncatedRows := truncateRows(rows, widths)

	// Create a new table
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			// row == 0 is header in lipgloss table rendering
			// All data rows should use normal style
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(0, 1)
		})

	// Add all rows
	t.Headers(truncatedRows[0]...)
	for i := 1; i < len(truncatedRows); i++ {
		t.Row(truncatedRows[i]...)
	}

	return t.Render()
}

func main() {
	// Read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var input strings.Builder

	for scanner.Scan() {
		input.WriteString(scanner.Text())
		input.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Parse the table
	rows := parseTable(input.String())

	if len(rows) == 0 {
		fmt.Println("No data found to format")
		return
	}

	// Get initial terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
		height = 24
	}

	// Initialize model
	m := model{
		rows:            rows,
		currentColumn:   0,
		selectedColumns: make(map[int]bool),
		viewMode:        NormalView,
		termWidth:       width,
		termHeight:      height,
	}

	// Start bubbletea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
