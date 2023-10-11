package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

	case list.ItemsChosenMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.list.KeyMap.SwitchPane):
			switch {
			case m.Pager.Focused():
				cmds = append(cmds, m.list.Focus())
				//m.state = Browsing
			case m.list.Focused():
				cmds = append(cmds, m.Pager.Focus())
				//m.state = Paging
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

func (m *Model) View() string {
	return m.list.View()
}
