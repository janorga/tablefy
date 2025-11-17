package model

import (
	"testing"
)

func TestApplyFuzzyFilter(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS", "PORT"},
		{"web-1", "running", "8080"},
		{"web-2", "stopped", "8081"},
		{"db-prod", "running", "5432"},
		{"db-test", "stopped", "5433"},
		{"cache-1", "running", "6379"},
	}

	// Test 1: Filter by "running" status
	indices := ApplyFuzzyFilter(rows, 1, "run")
	if len(indices) != 3 {
		t.Errorf("Expected 3 running entries, got %d", len(indices))
	}

	// Test 2: Filter by "stop" status
	indices = ApplyFuzzyFilter(rows, 1, "stop")
	if len(indices) != 2 {
		t.Errorf("Expected 2 stopped entries, got %d", len(indices))
	}

	// Test 3: Filter by "web" in NAME
	indices = ApplyFuzzyFilter(rows, 0, "web")
	if len(indices) != 2 {
		t.Errorf("Expected 2 web entries, got %d", len(indices))
	}

	// Test 4: Empty query returns all rows
	indices = ApplyFuzzyFilter(rows, 1, "")
	if len(indices) != 5 {
		t.Errorf("Expected 5 rows for empty query, got %d", len(indices))
	}
}

func TestGetFilteredRows(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS"},
		{"web-1", "running"},
		{"web-2", "stopped"},
		{"db-prod", "running"},
	}

	// Get rows at indices 1 and 3
	filtered := GetFilteredRows(rows, []int{1, 3})

	// Should have header + 2 data rows
	if len(filtered) != 3 {
		t.Errorf("Expected 3 rows (header + 2 data), got %d", len(filtered))
	}

	// Check header is included
	if filtered[0][0] != "NAME" {
		t.Errorf("Expected header, got %v", filtered[0])
	}
}
