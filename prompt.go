package teacozy

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Prompt struct {
	*List
}

func NewPrompt(title string, items Items) Prompt {
	return Prompt{
		List: NewList(title, items),
	}
}

func (m *Prompt) Start() *Prompt {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Quit):
			cmds = append(cmds, tea.Quit)
		}
	}
	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
