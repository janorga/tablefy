package model

import (
	"strings"
	"testing"
)

// TestExportDataZoomViewOnlySelectedColumns verifies that zoom exports ONLY selected columns
func TestExportDataZoomViewOnlySelectedColumns(t *testing.T) {
	rows := [][]string{
		{"NAME", "AGE", "CITY", "COUNTRY"},
		{"John", "25", "London", "UK"},
		{"Alice", "30", "Paris", "France"},
		{"Bob", "28", "Berlin", "Germany"},
	}

	m := New(rows, 80, 24)

	// Select only NAME and CITY columns (indices 0 and 2)
	m.SelectedColumns[0] = true
	m.SelectedColumns[2] = true

	// Enter zoom mode
	m.ViewMode = ZoomView

	output := m.GetExportData()
	
	t.Logf("ViewMode: %v (ZoomView=%v)", m.ViewMode, ZoomView)
	t.Logf("SelectedColumns: %v", m.SelectedColumns)
	t.Logf("\nExported output:\n%s\n", output)

	// Should NOT contain AGE or COUNTRY
	if strings.Contains(output, "25") || strings.Contains(output, "30") || strings.Contains(output, "28") {
		t.Errorf("ERROR: Export contains AGE values (25, 30, 28) when it should not!")
		t.Errorf("Full output:\n%s", output)
	}

	if strings.Contains(output, "UK") || strings.Contains(output, "France") || strings.Contains(output, "Germany") {
		t.Errorf("ERROR: Export contains COUNTRY values when it should not!")
		t.Errorf("Full output:\n%s", output)
	}

	// Should contain NAME and CITY
	if !strings.Contains(output, "John") || !strings.Contains(output, "London") {
		t.Errorf("ERROR: Export missing NAME or CITY columns!")
		t.Errorf("Full output:\n%s", output)
	}

	// Verify each line has exactly 2 values (NAME and CITY)
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			t.Errorf("Line %d has %d columns (expected 2): %s", i, len(parts), line)
		}
	}
}

// TestExportDataNormalViewWithAllColumns tests that normal view exports ALL columns
func TestExportDataNormalViewWithAllColumns(t *testing.T) {
	rows := [][]string{
		{"NAME", "AGE", "CITY", "COUNTRY"},
		{"John", "25", "London", "UK"},
		{"Alice", "30", "Paris", "France"},
	}

	m := New(rows, 80, 24)
	m.ViewMode = NormalView

	output := m.GetExportData()
	
	t.Logf("ViewMode: NORMAL")
	t.Logf("\nExported output:\n%s\n", output)

	// Should contain ALL columns
	if !strings.Contains(output, "John") || !strings.Contains(output, "25") || 
	   !strings.Contains(output, "London") || !strings.Contains(output, "UK") {
		t.Errorf("ERROR: Normal view should export all columns!")
		t.Errorf("Full output:\n%s", output)
	}

	// Verify each line has exactly 4 values (all columns)
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 4 {
			t.Errorf("Line %d in normal view has %d columns (expected 4): %s", i, len(parts), line)
		}
	}
}
