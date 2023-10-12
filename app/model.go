package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/input"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
	"github.com/ohzqq/teacozy/util"
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
	width  int
	height int
	KeyMap KeyMap

	List *list.Model

	// input
	Input      *input.Model
	inputValue string
	showInput  bool

	// view
	Pager *pager.Model
}

func New() *Model {
	m := &Model{
		KeyMap: DefaultKeyMap(),
		layout: Vertical,
	}
	m.SetSize(util.TermSize())
	return m
}

func (m Model) Height() int {
	return m.height
}

func (m Model) Width() int {
	return m.width
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h - 1

	nw := m.Width()
	nh := m.Height()
	if m.HasPager() && m.HasList() {
		switch m.layout {
		case Vertical:
			nh = m.Height() / 2
		case Horizontal:
			nw = m.Width() / 2
		}
	}

	if m.HasPager() {
		m.Pager.SetSize(nw, nh)
	}
	if m.HasList() {
		m.List.SetSize(nw, nh)
	}
	if m.HasInput() {
		m.Input.Width = m.Width()
	}
}

func (m *Model) SetList(l *list.Model) *Model {
	m.state = List
	l.SetHeight(m.Height())
	l.SetShowFilter(false)
	l.SetShowInput(false)
	l.SetShowHelp(false)
	m.List = l
	m.KeyMap.List = l.KeyMap
	return m
}

func (m *Model) SetInput(prompt string, k key.Binding) *Model {
	m.Input = input.New()
	m.Input.Prompt = prompt
	m.KeyMap.Input = m.Input.KeyMap
	m.Input.FocusKey = k
	return m
}

func (m *Model) SetPager(l *pager.Model) *Model {
	m.Pager = l
	m.KeyMap.Pager = l.KeyMap
	return m
}

func (m Model) showFilter() bool {
	if m.HasList() {
		return m.List.SettingFilter()
	}
	return false
}

func (m *Model) Init() tea.Cmd {
	m.SetSize(m.Width(), m.Height())
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
		m.SetSize(msg.Width, msg.Height)

	case input.FocusMsg:
		if m.HasInput() {
			cmds = append(cmds, m.SetFocus(Input))
		}
	case input.UnfocusMsg:
		if m.HasInput() {
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
		case key.Matches(msg, m.Input.FocusKey):
			if m.HasInput() && !m.showFilter() {
				cmds = append(cmds, m.SetFocus(Input))
			}
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
		switch {
		case key.Matches(msg, m.List.KeyMap.CancelWhileFiltering, m.List.KeyMap.ClearFilter, m.List.KeyMap.AcceptWhileFiltering):
		case key.Matches(msg, m.List.KeyMap.Filter):
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
	var sections []string

	var li string
	if m.HasList() {
		li = m.List.View()
	}

	var page string
	if m.HasPager() {
		page = m.Pager.View()
	}

	var main string
	switch m.layout {
	case Vertical:
		main = lipgloss.JoinVertical(lipgloss.Left, li, page)
	case Horizontal:
		main = lipgloss.JoinHorizontal(lipgloss.Center, page, li)
	}
	sections = append(sections, main)

	switch m.State() {
	case Input:
		in := m.Input.View()
		sections = append(sections, in)
	case List:
		switch {
		case m.showFilter():
			sections = append(sections, m.List.FilterInput.View())
		case m.List.State() == list.Input:
			sections = append(sections, m.List.Input.View())
		default:
			sections = append(sections, "")
		}
	default:
		sections = append(sections, "")
	}

	//views = append(views, fmt.Sprintf("filter is applied %v\n", m.List.IsFiltered()))
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// SetShowInput shows or hides the input model.
func (m *Model) SetShowInput(show bool) {
	if show {
		m.SetFocus(Input)
		m.SetSize(m.Width(), m.Height()-1)
		return
	}
	m.SetFocus(List)
	m.SetSize(m.Width(), m.Height()+1)
}

// ResetInput resets the current input state.
//func (m *Model) ResetInput() {
//m.resetInput()
//}

//func (m *Model) resetInput() {
//if m.state == List {
//return
//}
//m.Input.Reset()
//m.Input.Blur()
//m.SetShowInput(false)
//}

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
