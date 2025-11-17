package main

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	keys []string
}

func (m model) Init() tea.Cmd   { return nil }

func (m model) View() string {
	output := "Key Detection Test\n"
	output += "==================\n"
	output += fmt.Sprintf("Detected %d keys so far:\n", len(m.keys))
	for i, k := range m.keys {
		output += fmt.Sprintf("%d. %s\n", i+1, k)
	}
	output += "\nPress PgUp, PgDn, or other keys (q to quit)\n"
	return output
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		keyStr := fmt.Sprintf("String=%q Type=%v Runes=%v", 
			keyMsg.String(), keyMsg.Type, keyMsg.Runes)
		m.keys = append(m.keys, keyStr)
		
		if keyMsg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
