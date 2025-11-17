package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"tablefy/internal/model"
	"tablefy/internal/parser"
	"tablefy/internal/terminal"
	"tablefy/internal/view"
)

// Config holds application configuration
type Config struct {
	AutoExpand bool
}

// Run starts the application
func Run(config Config) error {
	// Read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var input strings.Builder

	for scanner.Scan() {
		input.WriteString(scanner.Text())
		input.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Parse the table
	rows := parser.ParseTable(input.String())

	if len(rows) == 0 {
		fmt.Println("No data found to format")
		return nil
	}

	// Get initial terminal size
	width, height, err := terminal.GetSize()
	if err != nil {
		width = 80
		height = 24
	}

	// Initialize model
	m := model.New(rows, width, height)
	m.AutoExpand = config.AutoExpand
	m.SetRenderer(view.Render)

	// Start bubbletea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
