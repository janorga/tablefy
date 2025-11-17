package model

import (
	"strings"
)

// ApplyFuzzyFilter applies fuzzy matching to filter rows based on a column and query
// Returns row indices where the column value fuzzy-matches the query
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

	// Normalize query for matching (lowercase, trim)
	queryLower := strings.ToLower(strings.TrimSpace(query))

	// Extract column values and find matches (skip header row 0)
	var filteredIndices []int
	for i := 1; i < len(rows); i++ {
		if columnIndex < len(rows[i]) {
			value := rows[i][columnIndex]
			valueLower := strings.ToLower(value)

			// Use fuzzy matching: check if query characters appear in order in value
			if fuzzyMatch(queryLower, valueLower) {
				filteredIndices = append(filteredIndices, i)
			}
		}
	}

	return filteredIndices
}

// fuzzyMatch checks if all characters in query appear in value in the same order
// This implements subsequence matching (e.g., "run" matches "rUnning" when case-insensitive)
func fuzzyMatch(query, value string) bool {
	queryIdx := 0
	for _, char := range value {
		if queryIdx < len(query) && char == rune(query[queryIdx]) {
			queryIdx++
		}
	}
	return queryIdx == len(query)
}

// fuzzyMatchSmart checks fuzzy match with scoring for better relevance
// Prefers matches where characters are closer together
func fuzzyMatchSmart(query, value string) (bool, float64) {
	queryIdx := 0
	valueIdx := 0
	charDistance := 0

	for valueIdx < len(value) && queryIdx < len(query) {
		if value[valueIdx] == query[queryIdx] {
			queryIdx++
		}
		charDistance++
		valueIdx++
	}

	// Only a match if all query characters were found
	if queryIdx != len(query) {
		return false, 0
	}

	// Score based on how many characters we had to skip
	// Lower score (closer match) is better
	score := float64(charDistance) / float64(len(query))
	return true, score
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
