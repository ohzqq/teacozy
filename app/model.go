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
	List State = iota + 1
	Input
	Pager
)

type Layout int

const (
	Vertical Layout = iota
	Horizontal
)

type Model struct {
	state  State
	layout Layout
	KeyMap KeyMap

	List *list.Model

	// input
	Input      *input.Model
	inputValue string

	// view
	Pager *pager.Model
}

func New() *Model {
	m := &Model{
		KeyMap: DefaultKeyMap(),
	}
	m.state = List
	return m
}

func (m *Model) SetList(l *list.Model) *Model {
	m.List = l
	m.KeyMap.List = l.KeyMap
	return m
}

func (m *Model) SetInput(prompt string) *Model {
	m.Input = input.New()
	m.Input.Prompt = prompt
	m.KeyMap.Input = m.Input.KeyMap
	return m
}

func (m *Model) SetPager(l *pager.Model) *Model {
	m.Pager = l
	m.KeyMap.Pager = l.KeyMap
	return m
}

func (m *Model) Init() tea.Cmd {
	switch {
	case m.HasList():
		m.state = List
		return m.List.Focus()
	case m.HasPager():
		m.state = Pager
		return m.Pager.Focus()
	case m.HasInput():
		m.state = Input
		return m.Input.Focus()
	}
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w := msg.Width
		h := msg.Height
		if m.HasPager() {
			switch m.layout {
			case Vertical:
				h = h / 2
			case Horizontal:
				w = w / 2
			}
			m.Pager.SetSize(w, h)
		}
		if m.HasList() {
			m.List.SetSize(w, h)
		}

	case input.FocusMsg:
		if m.HasInput() {
			m.SetShowInput(true)
			cmds = append(cmds, m.SetFocus(Input))
		}
	case input.UnfocusMsg:
		if m.HasInput() {
			m.Input.Reset()
			m.Input.Blur()
			m.SetShowInput(false)
			cmds = append(cmds, m.SetFocus(List))
		}
	case input.InputValueMsg:
		m.inputValue = msg.Value

	case pager.UnfocusMsg:
	case pager.FocusMsg:
	case list.FocusMsg:

	case list.ItemsChosenMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case key.Matches(msg, m.KeyMap.ToggleFocus):
			if m.HasPager() {
				switch {
				case m.Pager.Focused():
					cmds = append(cmds, m.SetFocus(List))
				case m.List.Focused():
					cmds = append(cmds, m.SetFocus(Pager))
				}
			}
		}
	}

	switch m.State() {
	case Input:
		if m.HasInput() {
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		}
	case Pager:
		if m.HasPager() {
			m.Pager, cmd = m.Pager.Update(msg)
			cmds = append(cmds, cmd)
		}
	case List:
		if m.HasList() {
			m.List, cmd = m.List.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) SetFocus(focus State) tea.Cmd {
	var cmds []tea.Cmd

	switch {
	case m.HasPager():
		if m.Pager.Focused() {
			cmds = append(cmds, m.Pager.Unfocus())
		}
	case m.HasList():
		if m.List.Focused() {
			cmds = append(cmds, m.List.Unfocus())
		}
	case m.HasInput():
		if m.Input.Focused() {
			cmds = append(cmds, m.Input.Unfocus())
		}
	}

	switch focus {
	case Input:
		m.state = Input
		if m.HasPager() {
			m.Pager.Unfocus()
		}
		if m.HasList() {
			m.List.Unfocus()
		}
		cmds = append(cmds, m.Input.Focus())
	case Pager:
		m.state = Pager
		if m.HasInput() {
			m.Input.Unfocus()
		}
		if m.HasList() {
			m.List.Unfocus()
		}
		cmds = append(cmds, m.Pager.Focus())
	case List:
		m.state = List
		if m.HasPager() {
			m.Pager.Unfocus()
		}
		if m.HasInput() {
			m.Input.Unfocus()
		}
		cmds = append(cmds, m.List.Focus())
	}

	return tea.Batch(cmds...)
}

func (m Model) State() State {
	return m.state
}

func (m Model) HasPager() bool {
	return m.Pager != nil
}

func (m Model) HasInput() bool {
	return m.Input != nil
}

func (m Model) HasList() bool {
	return m.List != nil
}

func (m *Model) View() string {
	var views []string

	if m.HasInput() {
		m.List.SetShowFilter(true)
		if m.Input.Focused() {
			m.List.SetShowFilter(false)
			in := m.Input.View()
			views = append(views, in)
		}
	}

	if m.HasList() {
		li := m.List.View()
		views = append(views, li)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, views...)

	var p string
	if m.HasPager() {
		p = m.Pager.View()
		switch m.layout {
		case Vertical:
			view = lipgloss.JoinVertical(lipgloss.Right, view, p)
		case Horizontal:
			view = lipgloss.JoinHorizontal(lipgloss.Center, p, view)
		}
	}
	if m.inputValue != "" {
		view += "\n"
		view += m.inputValue
	}
	return view
}

// SetShowInput shows or hides the input model.
func (m *Model) SetShowInput(show bool) {
	m.List.SetShowTitle(!show)
	if show {
		m.List.SetHeight(m.List.Height() - 1)
		m.SetFocus(Input)
		return
	}
	m.List.SetHeight(m.List.Height() + 1)
	m.SetFocus(List)
}

// ResetInput resets the current input state.
func (m *Model) ResetInput() {
	m.resetInput()
}

func (m *Model) resetInput() {
	if m.state == List {
		return
	}
	m.Input.Reset()
	m.Input.Blur()
	m.SetShowInput(false)
}

func (s State) String() string {
	switch s {
	case List:
		return "list"
	case Input:
		return "input"
	case Pager:
		return "pager"
	}
	return ""
}
