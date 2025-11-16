package parser

import "strings"

// ParseTable parses the input and converts it to rows and columns
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
