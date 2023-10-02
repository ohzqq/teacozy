package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/textinput"
)

type Model struct {
	textinput.Model
	Enter EnterInput
}

type EnterInput func(string) tea.Cmd

type ResetInputMsg struct{}
type FocusInputMsg struct{}

func New() *Model {
	m := &Model{
		Model: textinput.New(),
	}
	return m
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Focused() {
			switch msg.Type {
			case tea.KeyCtrlC:
				cmds = append(cmds, tea.Quit)
			case tea.KeyEsc:
				cmds = append(cmds, m.Reset)
			case tea.KeyEnter:
				val := m.Value()
				cmd := m.Enter(val)
				cmds = append(cmds, cmd)
				cmds = append(cmds, m.Reset)
			}
		}
	}
	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) Reset() tea.Msg {
	m.Model.Reset()
	m.Model.Blur()
	return ResetInputMsg{}
}

func Focus() tea.Msg {
	return FocusInputMsg{}
}

func Reset() tea.Msg {
	return ResetInputMsg{}
}

func (m Model) View() string {
	return m.Model.View()
}

func (m Model) Init() tea.Cmd {
	return nil
}
