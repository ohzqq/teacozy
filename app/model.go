package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/input"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
	"golang.org/x/term"
)

type State int

const (
	List State = iota + 1
	Cmd
	Pager
)

type Layout int

const (
	Vertical Layout = iota
	Horizontal
)

type Command struct {
	Name string
	Key  key.Binding
	Cmd  func(string) tea.Cmd
}

type Model struct {
	state  State
	layout Layout
	width  int
	height int
	KeyMap KeyMap

	List *list.Model

	// input
	Command    *input.Model
	Commands   []Command
	inputValue string
	showInput  bool

	// view
	Pager *pager.Model
}

func New() *Model {
	m := &Model{
		KeyMap: DefaultKeyMap(),
		layout: Vertical,
		Commands: []Command{
			Command{
				Name: "",
				Cmd:  func(string) tea.Cmd { return nil },
			},
		},
	}

	m.Command = input.New()
	m.Command.Prompt = ":"
	m.KeyMap.Input = m.Command.KeyMap

	m.SetSize(TermSize())
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
	m.height = h - 2

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

	m.Command.Width = m.Width()
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

func (m *Model) SetPager(l *pager.Model) *Model {
	m.state = Pager
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
		//return m.List.Focus()
	case m.HasPager():
		m.state = Pager
		//return m.Pager.Focus()
	case m.HasInput():
		m.state = Cmd
		//return m.Command.Focus()
	}
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case input.InputValueMsg:
		if c, arg, ok := strings.Cut(msg.Value, " "); ok {
			cmd = m.RunCommand(c, arg)
			cmds = append(cmds, cmd)
		}
		//m.Command.Reset()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case key.Matches(msg, m.KeyMap.Command):
			if !m.showFilter() {
				cmds = append(cmds, m.SetFocus(Cmd))
			}
		}
	}

	switch m.State() {
	case Cmd:
		cmd = m.updateCommand(msg)
		cmds = append(cmds, cmd)
	case Pager:
		if m.HasPager() {
			cmd = m.updatePager(msg)
			cmds = append(cmds, cmd)
		}
	case List:
		if m.HasList() {
			cmd = m.updateList(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updatePager(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ToggleFocus):
			if m.HasList() {
				switch m.State() {
				case Pager:
					cmds = append(cmds, m.SetFocus(List))
				case List:
					cmds = append(cmds, m.SetFocus(Pager))
				}
			}
		}
	}

	m.Pager, cmd = m.Pager.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) updateList(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case list.ItemsChosenMsg:
		return tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ToggleFocus):
			if m.HasPager() {
				switch m.State() {
				case Pager:
					cmds = append(cmds, m.SetFocus(List))
				case List:
					cmds = append(cmds, m.SetFocus(Pager))
				}
			}
		}
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

type RunCommandMsg Command

func (m *Model) RunCommand(cmd, arg string) tea.Cmd {
	return func() tea.Msg {
		for _, c := range m.Commands {
			if c.Name == cmd {
				m.inputValue = fmt.Sprintf("cmd %s val %s\n", cmd, arg)
				return c.Cmd(arg)
			}
		}
		return nil
	}
}

func (m *Model) updateCommand(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case input.FocusMsg:
		//cmds = append(cmds, m.SetFocus(Cmd))
	case input.UnfocusMsg:
		switch {
		case m.HasPager():
			cmds = append(cmds, m.SetFocus(Pager))
		case m.HasList():
			cmds = append(cmds, m.SetFocus(List))
		}
	case input.InputValueMsg:
		m.inputValue = msg.Value

	case tea.KeyMsg:
		//switch {
		//case key.Matches(msg, m.Input.FocusKey):
		//  if !m.showFilter() {
		//    cmds = append(cmds, m.SetFocus(Input))
		//  }
		//}

	}

	m.Command, cmd = m.Command.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
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
	default:
		if m.Command.Focused() {
			cmds = append(cmds, m.Command.Unfocus())
		}
	}

	switch focus {
	case Cmd:
		m.state = Cmd
		if m.HasPager() {
			m.Pager.Unfocus()
		}
		if m.HasList() {
			m.List.Unfocus()
		}
		cmds = append(cmds, m.Command.Focus())
	case Pager:
		m.state = Pager
		m.Command.Unfocus()
		if m.HasList() {
			m.List.Unfocus()
		}
		cmds = append(cmds, m.Pager.Focus())
	case List:
		m.state = List
		if m.HasPager() {
			m.Pager.Unfocus()
		}
		m.Command.Unfocus()
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
	return m.Command != nil
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
	case Cmd:
		in := m.Command.View()
		sections = append(sections, in)
	case List:
		if m.HasList() {
			switch {
			case m.showFilter():
				sections = append(sections, m.List.FilterInput.View())
			case m.List.State() == list.Input:
				sections = append(sections, m.List.Input.View())
			default:
				sections = append(sections, "")
			}
		}
	default:
		status := ""
		//status = fmt.Sprintf("state %s has list %v", m.state, m.HasList())
		sections = append(sections, status)
		//sections = append(sections, m.state.String())
	}
	if m.inputValue != "" {
		sections = append(sections, m.inputValue)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (s State) String() string {
	switch s {
	case List:
		return "list"
	case Cmd:
		return "command"
	case Pager:
		return "pager"
	}
	return ""
}

func TermSize() (int, int) {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return w, h
}
