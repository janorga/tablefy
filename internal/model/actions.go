package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height
	}
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle FilterView input separately
	if m.ViewMode == FilterView {
		return m.handleFilterViewInput(msg)
	}

	// Handle special key types first (more efficient than string comparison)
	switch msg.Type {
	case tea.KeyPgUp:
		// Page Up (Re Pág) - scroll up by page size
		pageSize := m.GetPageSize()
		m.ScrollOffset -= pageSize
		if m.ScrollOffset < 0 {
			m.ScrollOffset = 0
		}
		return m, nil
	case tea.KeyType(-10): // Page Down (Av Pág)
		// Page Down - scroll down by page size
		pageSize := m.GetPageSize()
		maxScroll := m.GetMaxScroll()
		m.ScrollOffset += pageSize
		if m.ScrollOffset > maxScroll {
			m.ScrollOffset = maxScroll
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "q":
		if m.ViewMode == ZoomView {
			// Exit zoom mode
			m.ViewMode = NormalView
			m.ScrollOffset = 0 // Reset scroll when exiting zoom
			return m, nil
		}
		return m, tea.Quit
	case "f", "F":
		if m.ViewMode == NormalView {
			// Enter filter mode
			m.ViewMode = FilterView
			m.FilterColumnIndex = m.CurrentColumn
			m.FilterInput = ""
			m.FilterScrollOffset = 0
			// Apply initial filter (empty query shows all rows)
			m.FilteredRowIndices = ApplyFuzzyFilter(m.Rows, m.FilterColumnIndex, "")
			return m, nil
		}
	case "c", "C":
		// Clear filter
		if m.ViewMode == NormalView && len(m.FilteredRowIndices) > 0 {
			m.ClearFilter()
			m.ScrollOffset = 0
			return m, nil
		}
	case "o", "O":
		// Export visible table and quit
		m.ExportData = m.GetExportData()
		return m, tea.Quit
	case "left", "h":
		if m.ViewMode == NormalView && m.CurrentColumn > 0 {
			m.CurrentColumn--
		}
	case "right", "l":
		if m.ViewMode == NormalView && len(m.Rows) > 0 && m.CurrentColumn < len(m.Rows[0])-1 {
			m.CurrentColumn++
		}
	case "up", "k":
		// Scroll up
		if m.ScrollOffset > 0 {
			m.ScrollOffset--
		}
	case "down", "j":
		// Scroll down
		maxScroll := m.GetMaxScroll()
		if m.ScrollOffset < maxScroll {
			m.ScrollOffset++
		}
	case "s", "S":
		// Toggle selection of current column
		if m.ViewMode == NormalView {
			if m.SelectedColumns[m.CurrentColumn] {
				delete(m.SelectedColumns, m.CurrentColumn)
			} else {
				m.SelectedColumns[m.CurrentColumn] = true
			}
		}
	case "enter", " ":
		if m.ViewMode == NormalView && len(m.SelectedColumns) > 0 {
			// Enter zoom mode with selected columns
			m.ViewMode = ZoomView
			m.ScrollOffset = 0 // Reset scroll when entering zoom
		}
	}
	return m, nil
}

// GetMaxScroll calculates the maximum scroll offset
func (m Model) GetMaxScroll() int {
	// Account for header, borders, and help text (approximately 6 lines)
	visibleRows := m.TermHeight - 6
	if visibleRows < 1 {
		visibleRows = 1
	}

	// Determine total data rows based on filter status
	var dataRows int
	if len(m.FilteredRowIndices) > 0 {
		dataRows = len(m.FilteredRowIndices)
	} else {
		dataRows = len(m.Rows) - 1 // All rows except header
	}

	maxScroll := dataRows - visibleRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

// GetPageSize calculates how many rows fit in a page for pgup/pgdn
func (m Model) GetPageSize() int {
	// Account for header, borders, and help text (approximately 6 lines)
	pageSize := m.TermHeight - 6
	if pageSize < 1 {
		pageSize = 1
	}
	return pageSize
}

// GetMaxFilterScroll calculates the maximum scroll offset for filter view
func (m Model) GetMaxFilterScroll() int {
	// Account for header, borders, filter input, and help text (approximately 7 lines)
	visibleRows := m.TermHeight - 7
	if visibleRows < 1 {
		visibleRows = 1
	}

	// Filtered rows count (including header)
	totalRows := len(m.FilteredRowIndices) + 1

	maxScroll := totalRows - visibleRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

// handleFilterViewInput handles keyboard input while in FilterView
func (m Model) handleFilterViewInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle special key types first (more efficient than string comparison)
	switch msg.Type {
	case tea.KeyPgUp:
		// Page Up (Re Pág) - scroll up by page size in filter view
		pageSize := m.GetPageSize()
		m.FilterScrollOffset -= pageSize
		if m.FilterScrollOffset < 0 {
			m.FilterScrollOffset = 0
		}
		return m, nil
	case tea.KeyType(-10): // Page Down (Av Pág)
		// Page Down - scroll down by page size in filter view
		pageSize := m.GetPageSize()
		maxFilterScroll := m.GetMaxFilterScroll()
		m.FilterScrollOffset += pageSize
		if m.FilterScrollOffset > maxFilterScroll {
			m.FilterScrollOffset = maxFilterScroll
		}
		return m, nil
	}

	switch msg.String() {
	case "esc":
		// Cancel filter and return to normal view
		m.ClearFilter()
		m.ViewMode = NormalView
		return m, nil
	case "enter":
		// Apply filter and return to normal view
		m.ViewMode = NormalView
		// Keep the filter applied - don't clear it
		return m, nil
	case "backspace":
		// Remove last character
		if len(m.FilterInput) > 0 {
			m.FilterInput = m.FilterInput[:len(m.FilterInput)-1]
			m.FilteredRowIndices = ApplyFuzzyFilter(m.Rows, m.FilterColumnIndex, m.FilterInput)
			m.FilterScrollOffset = 0
		}
	case "up", "k":
		// Scroll up in filter view
		if m.FilterScrollOffset > 0 {
			m.FilterScrollOffset--
		}
	case "down", "j":
		// Scroll down in filter view
		maxFilterScroll := m.GetMaxFilterScroll()
		if m.FilterScrollOffset < maxFilterScroll {
			m.FilterScrollOffset++
		}
	case "o", "O":
		// Export filtered table and quit
		m.ExportData = m.GetExportData()
		return m, tea.Quit
	default:
		// Add character to filter input
		if len(msg.String()) == 1 {
			m.FilterInput += msg.String()
			m.FilteredRowIndices = ApplyFuzzyFilter(m.Rows, m.FilterColumnIndex, m.FilterInput)
			m.FilterScrollOffset = 0
		}
	}
	return m, nil
}
