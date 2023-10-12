package input

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textinput.Model
	Enter    EnterInput
	FocusKey key.Binding
}

type EnterInput func(string) tea.Cmd

type ResetInputMsg struct{}
type FocusMsg struct{}
type UnfocusMsg struct{}

func New() *Model {
	m := &Model{
		Model: textinput.New(),
		Enter: InputValue,
	}
	return m
}

func (m *Model) WithKey(k key.Binding) *Model {
	m.FocusKey = k
	return m
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		}
		if m.Focused() {
			switch msg.Type {
			case tea.KeyEsc:
				cmds = append(cmds, m.Reset)
			case tea.KeyEnter:
				val := m.Value()
				cmd := m.Enter(val)
				cmds = append(cmds, cmd)
				cmds = append(cmds, m.Unfocus())
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

func (m *Model) Focus() tea.Cmd {
	return m.Model.Focus()
}

func (m *Model) Unfocus() tea.Cmd {
	return func() tea.Msg {
		m.Model.Reset()
		m.Model.Blur()
		return UnfocusMsg{}
	}
}

type InputValueMsg struct {
	Value string
}

func InputValue(val string) tea.Cmd {
	return func() tea.Msg {
		return InputValueMsg{
			Value: val,
		}
	}
}

func Focus() tea.Msg {
	return FocusMsg{}
}

func Unfocus() tea.Msg {
	return UnfocusMsg{}
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
