package info

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	*Form
}

func NewForm(data FormData) *Model {
	return &Model{Form: New(data)}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
	}

	m.Form, cmd = m.Form.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return m.Form.View()
}
