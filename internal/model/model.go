package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
	Rows               [][]string
	CurrentColumn      int
	SelectedColumns    map[int]bool
	ViewMode           ViewMode
	ScrollOffset       int
	TermWidth          int
	TermHeight         int
	AutoExpand         bool
	FilterInput        string
	FilteredRowIndices []int
	FilterColumnIndex  int
	FilterScrollOffset int
	ExportData         string // Data to export when quitting with 'o'
	renderer           func(Model) string
}

// New creates a new model with the given rows
func New(rows [][]string, termWidth, termHeight int) Model {
	return Model{
		Rows:            rows,
		CurrentColumn:   0,
		SelectedColumns: make(map[int]bool),
		ViewMode:        NormalView,
		ScrollOffset:    0,
		TermWidth:       termWidth,
		TermHeight:      termHeight,
		AutoExpand:      false,
	}
}

// SetRenderer sets the view renderer function
func (m *Model) SetRenderer(renderer func(Model) string) {
	m.renderer = renderer
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// View renders the UI using the provided renderer
func (m Model) View() string {
	if m.renderer != nil {
		return m.renderer(m)
	}
	return ""
}
