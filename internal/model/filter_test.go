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

func TestFuzzyFilterRefinement(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS"},
		{"pod-1", "running"},
		{"pod-2", "running"},
		{"pod-3", "pending"},
		{"pod-4", "running"},
		{"pod-5", "completed"},
	}

	// Test refinement: more specific queries should return fewer results
	allMatches := ApplyFuzzyFilter(rows, 1, "r")       // matches "running" (3 matches)
	runMatches := ApplyFuzzyFilter(rows, 1, "run")     // matches "running" (3 matches)
	runningMatches := ApplyFuzzyFilter(rows, 1, "run") // matches "running" (3 matches - exact substring)

	// "r" should match "running" (3) + "pending" (0, doesn't have r) = 3
	if len(allMatches) != 3 {
		t.Errorf("Query 'r': expected 3 matches, got %d", len(allMatches))
	}

	// "run" should match same as "r"
	if len(runMatches) != 3 {
		t.Errorf("Query 'run': expected 3 matches, got %d", len(runMatches))
	}

	// "running" exact - should match 3
	if len(runningMatches) != 3 {
		t.Errorf("Query 'running': expected 3 matches, got %d", len(runningMatches))
	}

	// "p" should match "pending" + "completed" (2 matches with 'p')
	pMatches := ApplyFuzzyFilter(rows, 1, "p")
	if len(pMatches) != 2 {
		t.Errorf("Query 'p': expected 2 matches (pending, completed), got %d", len(pMatches))
	}

	// "pe" should match only "pending" and "completed" (both have 'p' then 'e')
	peMatches := ApplyFuzzyFilter(rows, 1, "pe")
	if len(peMatches) != 2 {
		t.Errorf("Query 'pe': expected 2 matches, got %d", len(peMatches))
	}

	// "pen" should match only "pending"
	penMatches := ApplyFuzzyFilter(rows, 1, "pen")
	if len(penMatches) != 1 {
		t.Errorf("Query 'pen': expected 1 match (pending), got %d", len(penMatches))
	}
}

func TestFuzzyMatchCaseSensitivity(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS"},
		{"pod-1", "Running"},
		{"pod-2", "RUNNING"},
		{"pod-3", "running"},
		{"pod-4", "Pending"},
	}

	// Should match case-insensitively
	matches := ApplyFuzzyFilter(rows, 1, "run")
	if len(matches) != 3 {
		t.Errorf("Expected 3 case-insensitive matches for 'run', got %d", len(matches))
	}

	// Should match case-insensitively for uppercase query
	matches = ApplyFuzzyFilter(rows, 1, "RUN")
	if len(matches) != 3 {
		t.Errorf("Expected 3 case-insensitive matches for 'RUN', got %d", len(matches))
	}
}

func TestFuzzyMatchPartialMatching(t *testing.T) {
	rows := [][]string{
		{"NAME", "STATUS"},
		{"pod-1", "running"},
		{"pod-2", "runner"},
		{"pod-3", "runtime"},
		{"pod-4", "pending"},
	}

	// "run" should match all three: running, runner, runtime
	matches := ApplyFuzzyFilter(rows, 1, "run")
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches for 'run', got %d", len(matches))
	}

	// "runn" should match "running" and "runner" (both have r, u, n, n)
	matches = ApplyFuzzyFilter(rows, 1, "runn")
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'runn', got %d (running, runner)", len(matches))
	}

	// "rung" should match only "running" (has r, u, n, g in order)
	matches = ApplyFuzzyFilter(rows, 1, "rung")
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for 'rung', got %d", len(matches))
	}

	// "runt" should match only "runtime" (has r, u, n, t in order)
	// "running" doesn't have 't' and "runner" doesn't have 't'
	matches = ApplyFuzzyFilter(rows, 1, "runt")
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for 'runt' (only runtime), got %d", len(matches))
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
