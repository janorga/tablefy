package parser

import (
	"testing"
)

func TestParseTableWithTimestamp(t *testing.T) {
	input := `NAME                              NAMESPACE       REVISION        UPDATED                                         STATUS          CHART                                   APP VERSION
argocd                          argocd          2               2025-07-29 20:02:02.89383204 +0200 CEST         deployed        argo-cd-8.2.3                           v3.0.12`

	rows := ParseTable(input)

	// Check we have the correct number of rows (header + 1 data row)
	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}

	// Check header has 7 columns
	expectedHeader := []string{"NAME", "NAMESPACE", "REVISION", "UPDATED", "STATUS", "CHART", "APP VERSION"}
	if len(rows[0]) != len(expectedHeader) {
		t.Errorf("Expected %d columns in header, got %d", len(expectedHeader), len(rows[0]))
	}

	for i, col := range expectedHeader {
		if rows[0][i] != col {
			t.Errorf("Header column %d: expected %q, got %q", i, col, rows[0][i])
		}
	}

	// Check data row
	dataRow := rows[1]
	if len(dataRow) != len(expectedHeader) {
		t.Errorf("Expected %d columns in data row, got %d", len(expectedHeader), len(dataRow))
	}

	// Check specific values
	if dataRow[0] != "argocd" {
		t.Errorf("Column NAME: expected 'argocd', got %q", dataRow[0])
	}

	if dataRow[1] != "argocd" {
		t.Errorf("Column NAMESPACE: expected 'argocd', got %q", dataRow[1])
	}

	if dataRow[2] != "2" {
		t.Errorf("Column REVISION: expected '2', got %q", dataRow[2])
	}

	// The timestamp should be preserved with its internal spaces
	expectedTimestamp := "2025-07-29 20:02:02.89383204 +0200 CEST"
	if dataRow[3] != expectedTimestamp {
		t.Errorf("Column UPDATED: expected %q, got %q", expectedTimestamp, dataRow[3])
	}

	if dataRow[4] != "deployed" {
		t.Errorf("Column STATUS: expected 'deployed', got %q", dataRow[4])
	}

	if dataRow[5] != "argo-cd-8.2.3" {
		t.Errorf("Column CHART: expected 'argo-cd-8.2.3', got %q", dataRow[5])
	}

	if dataRow[6] != "v3.0.12" {
		t.Errorf("Column APP VERSION: expected 'v3.0.12', got %q", dataRow[6])
	}
}

// Test with actual helm list -A format (tab-separated)
func TestParseTableHelmFormat(t *testing.T) {
	input := `NAME	NAMESPACE	REVISION	UPDATED	STATUS	CHART	APP VERSION
argocd	argocd	2	2025-07-29 20:02:02.89383204 +0200 CEST	deployed	argo-cd-8.2.3	v3.0.12
longhorn	longhorn-system	2	2025-08-01 14:13:13.960555922 +0200 CEST	deployed	longhorn-1.9.1	v1.9.1`

	rows := ParseTable(input)

	// Check we have the correct number of rows (header + 2 data rows)
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(rows))
	}

	// Check header has 7 columns
	expectedHeader := []string{"NAME", "NAMESPACE", "REVISION", "UPDATED", "STATUS", "CHART", "APP VERSION"}
	if len(rows[0]) != len(expectedHeader) {
		t.Errorf("Expected %d columns in header, got %d", len(expectedHeader), len(rows[0]))
	}

	// Check all rows have the same number of columns
	for i, row := range rows {
		if len(row) != len(expectedHeader) {
			t.Errorf("Row %d has %d columns, expected %d", i, len(row), len(expectedHeader))
		}
	}

	// Check specific values from row 1
	if rows[1][0] != "argocd" {
		t.Errorf("Row 1, Column 0: expected 'argocd', got %q", rows[1][0])
	}

	expectedTimestamp := "2025-07-29 20:02:02.89383204 +0200 CEST"
	if rows[1][3] != expectedTimestamp {
		t.Errorf("Row 1, Column 3 (UPDATED): expected %q, got %q", expectedTimestamp, rows[1][3])
	}

	// Check specific values from row 2
	if rows[2][0] != "longhorn" {
		t.Errorf("Row 2, Column 0: expected 'longhorn', got %q", rows[2][0])
	}

	if rows[2][1] != "longhorn-system" {
		t.Errorf("Row 2, Column 1: expected 'longhorn-system', got %q", rows[2][1])
	}
}

func TestSplitByMultipleSpaces(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "col1  col2   col3",
			expected: []string{"col1", "col2", "col3"},
		},
		{
			input:    "one  two three  four",
			expected: []string{"one", "two three", "four"},
		},
		{
			input:    "2025-07-29 20:02:02  deployed",
			expected: []string{"2025-07-29 20:02:02", "deployed"},
		},
		{
			// Tab-separated (like helm output)
			input:    "NAME\tNAMESPACE\tUPDATED",
			expected: []string{"NAME", "NAMESPACE", "UPDATED"},
		},
		{
			// Tab with spaces inside values
			input:    "argocd\targocd\t2025-07-29 20:02:02.89383204 +0200 CEST",
			expected: []string{"argocd", "argocd", "2025-07-29 20:02:02.89383204 +0200 CEST"},
		},
	}

	for _, tt := range tests {
		result := splitByMultipleSpaces(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("For input %q: expected %d fields, got %d", tt.input, len(tt.expected), len(result))
		}
		for i, v := range tt.expected {
			if i >= len(result) || result[i] != v {
				t.Errorf("For input %q: field %d expected %q, got %q", tt.input, i, v, result[i])
			}
		}
	}
}
