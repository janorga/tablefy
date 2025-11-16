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

	// Number of data rows (excluding header)
	dataRows := len(m.Rows) - 1

	maxScroll := dataRows - visibleRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}
