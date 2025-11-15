package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

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

	// Format and display
	formatted := formatTable(rows)
	fmt.Println(formatted)
}
