package parser

import (
	"regexp"
	"strings"
)

// ColumnPosition represents a column with its starting position and name
type ColumnPosition struct {
	Start int    // Starting position of the column in the line
	Name  string // Column name from the header
}

// splitByMultipleSpaces splits a line by tabs first, then falls back to 2+ consecutive spaces
func splitByMultipleSpaces(line string) []string {
	var parts []string

	// First try to split by tabs (most common in command output like helm, kubectl)
	if strings.Contains(line, "\t") {
		parts = strings.Split(line, "\t")
	} else {
		// Fallback: split by 2 or more consecutive spaces
		re := regexp.MustCompile(`\s{2,}`)
		parts = re.Split(line, -1)
	}

	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// extractColumnPositions extracts the start position of each column from the header line
// by finding where 2+ consecutive spaces occur and marking column starts
func extractColumnPositions(headerLine string) []ColumnPosition {
	// Get column names using the proven split logic
	columnNames := splitByMultipleSpaces(headerLine)
	if len(columnNames) == 0 {
		return nil
	}

	// Find where 2+ spaces occur in the header to identify column starts
	re := regexp.MustCompile(`\s{2,}`)
	matches := re.FindAllStringIndex(headerLine, -1)

	// Column positions: start at 0, and right after each multi-space separator
	var columnStarts []int
	columnStarts = append(columnStarts, 0)
	for _, match := range matches {
		columnStarts = append(columnStarts, match[1])
	}

	// Create ColumnPosition objects with the column starts and names
	var positions []ColumnPosition
	for i, start := range columnStarts {
		if i < len(columnNames) {
			positions = append(positions, ColumnPosition{
				Start: start,
				Name:  columnNames[i],
			})
		}
	}

	return positions
}

// extractValuesByPosition extracts column values from a data line using column positions
// Each value starts at the column's start position and extends until the next column starts
func extractValuesByPosition(line string, positions []ColumnPosition) []string {
	var result []string

	for i, pos := range positions {
		var value string
		var endPos int

		// Determine where this column ends (start of next column or end of line)
		if i+1 < len(positions) {
			endPos = positions[i+1].Start
		} else {
			endPos = len(line)
		}

		// Make sure we don't go past the line length
		start := pos.Start
		if start > len(line) {
			start = len(line)
		}
		if endPos > len(line) {
			endPos = len(line)
		}

		// Extract substring and trim whitespace
		if start < len(line) && endPos > start {
			value = strings.TrimSpace(line[start:endPos])
		}

		result = append(result, value)
	}

	return result
}

// ParseTable parses the input and converts it to rows and columns
// It uses column positions from the header to align data rows correctly when space-separated
// For tab-separated data, it uses simple tab splitting
func ParseTable(input string) [][]string {
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

	// Check if the data is tab-separated (like helm output)
	header := validLines[0]
	isTabSeparated := strings.Contains(header, "\t")

	if isTabSeparated {
		// For tab-separated data, use simple splitting
		var rows [][]string
		for _, line := range validLines {
			fields := strings.Split(line, "\t")
			rows = append(rows, fields)
		}
		return rows
	}

	// For space-separated data, use position-based alignment
	columnPositions := extractColumnPositions(header)

	if len(columnPositions) == 0 {
		return nil
	}

	// Extract header row
	var rows [][]string
	headerRow := make([]string, len(columnPositions))
	for i, pos := range columnPositions {
		headerRow[i] = pos.Name
	}
	rows = append(rows, headerRow)

	// Process each data line using column positions
	for i := 1; i < len(validLines); i++ {
		row := extractValuesByPosition(validLines[i], columnPositions)

		// Ensure row has the correct number of columns (fill missing ones with empty strings)
		if len(row) < len(columnPositions) {
			for len(row) < len(columnPositions) {
				row = append(row, "")
			}
		}

		rows = append(rows, row)
	}

	return rows
}
