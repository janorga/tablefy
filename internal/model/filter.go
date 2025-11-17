package model

import (
	"github.com/schollz/closestmatch"
)

// ApplyFuzzyFilter applies fuzzy matching to filter rows based on a column and query
func ApplyFuzzyFilter(rows [][]string, columnIndex int, query string) []int {
	if len(rows) == 0 || columnIndex < 0 || columnIndex >= len(rows[0]) {
		return []int{}
	}

	// If query is empty, return all rows except header
	if query == "" {
		indices := make([]int, len(rows)-1)
		for i := 1; i < len(rows); i++ {
			indices[i-1] = i
		}
		return indices
	}

	// Extract column values (skip header row 0)
	var columnValues []string
	for i := 1; i < len(rows); i++ {
		if columnIndex < len(rows[i]) {
			columnValues = append(columnValues, rows[i][columnIndex])
		}
	}

	if len(columnValues) == 0 {
		return []int{}
	}

	// Create closest match finder
	cm := closestmatch.New(columnValues, []int{2})

	// Find matches for the query
	matches := cm.ClosestN(query, len(columnValues))

	// Convert matched values back to row indices
	var filteredIndices []int
	valueToRowIndex := make(map[string][]int)

	// Build map from value to row indices
	for i := 1; i < len(rows); i++ {
		if columnIndex < len(rows[i]) {
			value := rows[i][columnIndex]
			valueToRowIndex[value] = append(valueToRowIndex[value], i)
		}
	}

	// Collect indices from matched values
	for _, match := range matches {
		if indices, exists := valueToRowIndex[match]; exists {
			filteredIndices = append(filteredIndices, indices...)
		}
	}

	return filteredIndices
}

// GetFilteredRows extracts rows at specified indices
func GetFilteredRows(rows [][]string, filteredIndices []int) [][]string {
	if len(rows) == 0 {
		return [][]string{}
	}

	// Always include header
	result := [][]string{rows[0]}

	for _, idx := range filteredIndices {
		if idx > 0 && idx < len(rows) {
			result = append(result, rows[idx])
		}
	}

	return result
}

// ClearFilter resets all filter-related fields
func (m *Model) ClearFilter() {
	m.FilterInput = ""
	m.FilteredRowIndices = []int{}
	m.FilterColumnIndex = -1
	m.FilterScrollOffset = 0
}
