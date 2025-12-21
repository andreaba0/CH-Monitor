package tui_view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model_home struct {
	Choices  []string
	cursor   int
	Selected map[int]struct{}
}

func (m Model_home) Init() tea.Cmd {
	return nil
}

func (m Model_home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.Selected[m.cursor]
			if ok {
				delete(m.Selected, m.cursor)
			} else {
				m.Selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m Model_home) View() string {
	s := "Hello, here are your options\n"
	for i, choice := range m.Choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit\n"

	return s
}

func home() {

}
