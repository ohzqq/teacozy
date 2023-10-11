package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/input"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
)

type State int

const (
	Browsing State = iota
	Input
	Paging
)

type Layout int

const (
	Vertical Layout = iota
	Horizontal
)

type Model struct {
	state  State
	layout Layout

	list *list.Model

	// input
	Input    *input.Model
	hasInput bool

	// view
	Pager    *pager.Model
	hasPager bool
}

func New(l *list.Model) *Model {
	m := &Model{}
	m.list = l
	m.Input = l.Input
	m.hasInput = true
	m.hasPager = true
	m.Pager = l.Pager
	return m
}

func (m *Model) Init() tea.Cmd {
	return m.list.Focus()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w := msg.Width
		h := msg.Height
		if m.hasPager {
			switch m.layout {
			case Vertical:
				h = h / 2
			case Horizontal:
				w = w / 2
			}
			m.Pager.SetSize(w, h)
		}
		m.list.SetSize(w, h)

	case input.FocusInputMsg:
		if m.hasInput {
			m.list.SetShowInput(true)
			cmds = append(cmds, m.Input.Focus())
		}
	case input.ResetInputMsg:
		m.list.ResetInput()
		m.Input.Unfocus()

	case pager.UnfocusPagerMsg:
		//m.Pager.Unfocus()
		cmds = append(cmds, m.list.NewStatusMessage(m.CurrentFocus()))
		//return m, m.list.Focus()

	case list.ItemsChosenMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.list.KeyMap.SwitchPane):
			switch {
			case m.Pager.Focused():
				m.list.Focus()
				cmds = append(cmds, m.SetFocus("list"))
			case m.list.Focused():
				cmds = append(cmds, m.SetFocus("pager"))
			}
		case key.Matches(msg, m.list.KeyMap.Filter):
			//m.state = Input
		}
	}

	switch {
	case m.Input.Focused():
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	case m.Pager.Focused():
		m.Pager, cmd = m.Pager.Update(msg)
		cmds = append(cmds, cmd)
	case m.list.Focused():
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) SetFocus(focus string) tea.Cmd {
	var cmds []tea.Cmd

	switch {
	case m.Pager.Focused():
		cmds = append(cmds, m.Pager.Unfocus())
	case m.list.Focused():
		cmds = append(cmds, m.list.Unfocus())
	case m.Input.Focused():
		cmds = append(cmds, m.Input.Unfocus())
	}

	switch focus {
	case "input":
		m.Pager.Unfocus()
		m.list.Unfocus()
		cmds = append(cmds, m.Input.Focus())
	case "pager":
		m.Input.Unfocus()
		m.list.Unfocus()
		cmds = append(cmds, m.Pager.Focus())
	case "list":
		m.Pager.Unfocus()
		m.Input.Unfocus()
		cmds = append(cmds, m.list.Focus())
	}

	return tea.Batch(cmds...)
}

func (m Model) CurrentFocus() string {
	switch {
	case m.Input.Focused():
		return "input"
	case m.Pager.Focused():
		return "pager"
	case m.list.Focused():
		return "list"
	}
	return ""
}

func (m *Model) View() string {
	view := m.list.View()
	return lipgloss.JoinVertical(lipgloss.Left, view, m.CurrentFocus())
}
