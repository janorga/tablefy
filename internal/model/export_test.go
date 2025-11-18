package model

import (
	"strings"
	"testing"
)

// TestGetExportData tests exporting all columns with proper alignment
func TestGetExportData(t *testing.T) {
	rows := [][]string{
		{"NAME", "AGE", "CITY"},
		{"John", "25", "London"},
		{"Alice", "30", "Paris"},
		{"Bob", "28", "Berlin"},
	}

	m := New(rows, 80, 24)

	output := m.GetExportData()
	lines := strings.Split(output, "\n")

	// Should have 3 lines (all rows except header)
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Verify each line has proper structure
	for i, line := range lines {
		if line == "" {
			t.Errorf("Line %d is empty", i)
		}
	}

	// Verify that columns are properly formatted
	if !strings.Contains(lines[0], "John") || !strings.Contains(lines[0], "25") || !strings.Contains(lines[0], "London") {
		t.Errorf("Line 0 should contain John, 25, and London")
	}
}

// TestGetExportDataWithFilter tests export with filtered data
func TestGetExportDataWithFilter(t *testing.T) {
	rows := [][]string{
		{"NAME", "AGE", "CITY"},
		{"John", "25", "London"},
		{"Alice", "30", "Paris"},
		{"Bob", "28", "Berlin"},
	}

	m := New(rows, 80, 24)

	// Simulate filtered data (only rows 1 and 3)
	m.FilteredRowIndices = []int{1, 3}

	output := m.GetExportData()
	lines := strings.Split(output, "\n")

	// Should have 2 lines (filtered rows only)
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	// Verify both lines have content
	for i, line := range lines {
		if line == "" {
			t.Errorf("Line %d is empty", i)
		}
	}

	// Verify content
	if !strings.Contains(lines[0], "John") {
		t.Errorf("Line 0 should contain John")
	}
	if !strings.Contains(lines[1], "Bob") {
		t.Errorf("Line 1 should contain Bob")
	}
}

// TestGetExportDataEmpty tests with empty data
func TestGetExportDataEmpty(t *testing.T) {
	rows := [][]string{}
	m := New(rows, 80, 24)

	output := m.GetExportData()
	if output != "" {
		t.Errorf("Expected empty string, got '%s'", output)
	}
}

// TestGetExportDataAlignment tests that data is properly aligned
func TestGetExportDataAlignment(t *testing.T) {
	rows := [][]string{
		{"NAME", "VALUE"},
		{"short", "1"},
		{"verylongname", "999"},
	}

	m := New(rows, 80, 24)

	output := m.GetExportData()
	lines := strings.Split(output, "\n")

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	// Verify both lines exist and have content
	for _, line := range lines {
		if len(line) == 0 {
			t.Errorf("Got empty line in aligned output")
		}
	}
}

// TestGetExportDataZoomView tests export in ZoomView with selected columns
func TestGetExportDataZoomView(t *testing.T) {
	rows := [][]string{
		{"NAME", "AGE", "CITY"},
		{"John", "25", "London"},
		{"Alice", "30", "Paris"},
		{"Bob", "28", "Berlin"},
	}

	m := New(rows, 80, 24)
	m.ViewMode = ZoomView

	// Select only NAME and CITY columns (indices 0 and 2)
	m.SelectedColumns[0] = true
	m.SelectedColumns[2] = true

	output := m.GetExportData()
	lines := strings.Split(output, "\n")

	// Should have 3 lines (all rows except header)
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Verify each line has only 2 columns (NAME and CITY)
	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			t.Errorf("Line %d: Expected 2 columns, got %d. Line: '%s'", i, len(parts), line)
		}
	}

	// Check actual values
	if !strings.Contains(lines[0], "John") || !strings.Contains(lines[0], "London") {
		t.Errorf("Line 0 should contain 'John' and 'London', got: '%s'", lines[0])
	}
}

// TestGetExportDataZoomViewSingleColumn tests export with single selected column
func TestGetExportDataZoomViewSingleColumn(t *testing.T) {
	rows := [][]string{
		{"NAME", "AGE", "CITY"},
		{"John", "25", "London"},
		{"Alice", "30", "Paris"},
	}

	m := New(rows, 80, 24)
	m.ViewMode = ZoomView

	// Select only AGE column (index 1)
	m.SelectedColumns[1] = true

	output := m.GetExportData()
	lines := strings.Split(output, "\n")

	// Should have 2 lines (2 data rows)
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	// Verify first line contains "25"
	if !strings.Contains(lines[0], "25") {
		t.Errorf("Expected '25' in first line, got '%s'", lines[0])
	}

	// Verify second line contains "30"
	if !strings.Contains(lines[1], "30") {
		t.Errorf("Expected '30' in second line, got '%s'", lines[1])
	}
}
