package app

import (
	"fmt"
	"os"
	"strings"
	"time"

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
	Cmd State = iota + 1
	List
	Pager
)

const (
	Vertical = iota
	Horizontal
)

type Command struct {
	Name string
	Key  key.Binding
	Cmd  func(string) tea.Cmd
}

var noItems = list.NewItems(func() []*list.Item { return []*list.Item{} })

type Model struct {
	mainView State
	state    State
	layout   *Layout
	KeyMap   KeyMap

	List         *list.Model
	showList     bool
	showItemDesc bool

	// input
	Command     *input.Model
	Commands    []Command
	showCommand bool

	// view
	Pager     *pager.Model
	showPager bool

	StatusMsg
}

func New(opts ...Option) *Model {
	m := &Model{
		KeyMap:   DefaultKeyMap(),
		mainView: Cmd,
		StatusMsg: StatusMsg{
			StatusMessageLifetime: time.Second,
		},
		Commands: []Command{
			Command{
				Name: "",
				Cmd:  NewStatusMessage,
			},
		},
	}
	for _, opt := range opts {
		opt(m)
	}

	m.Command = input.New()
	m.Command.Prompt = ":"
	m.KeyMap.Input = m.Command.KeyMap

	if m.layout == nil {
		m.layout = NewLayout()
		if m.HasList() && m.HasPager() {
			m.layout.Half()
			m.layout.Position(Top)
		}
	}
	if m.HasPager() {
		m.mainView = Pager
	}
	if m.HasList() {
		m.mainView = List
		m.List.SetShowInput(false)
		m.List.Styles.FilterPrompt = m.Command.Style.Prompt
	}

	m.SetSize(TermSize())

	return m
}

func (m *Model) Run() error {
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		return err
	}
	return nil
}

func (m Model) Height() int {
	return m.layout.Height()
}

func (m Model) Width() int {
	return m.layout.Width()
}

func (m *Model) SetSize(w, h int) {
	m.layout.SetSize(w, h)

	switch m.mainView {
	case List:
		if m.HasList() {
			m.List.SetSize(m.layout.main())
		}
		if m.HasPager() {
			m.Pager.SetSize(m.layout.sub())
		}
	case Pager:
		if m.HasPager() {
			m.Pager.SetSize(m.layout.main())
		}
		if m.HasList() {
			m.List.SetSize(m.layout.sub())
		}
	}

	m.Command.Width = m.Width()
}

func (m *Model) AddCommands(cmds ...Command) *Model {
	m.Commands = append(m.Commands, cmds...)
	return m
}

func (m *Model) SetList(l *list.Model) *Model {
	m.state = List
	m.showList = true
	l.SetShowFilter(false)
	l.SetShowInput(false)
	l.SetShowHelp(false)
	m.List = l
	m.KeyMap.List = l.KeyMap
	return m
}

func (m *Model) SetPager(p *pager.Model) *Model {
	m.showPager = true
	m.state = Pager
	p.SetContent(p.Render())
	m.Pager = p
	m.KeyMap.Pager = p.KeyMap
	return m
}

func (m *Model) Init() tea.Cmd {
	switch {
	case m.HasList():
		m.state = List
	case m.HasPager():
		m.state = Pager
	case m.ShowCommand():
		m.state = Cmd
		return m.SetFocus(Cmd)
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
		status := fmt.Sprintf("%s is not a command", msg.Value)
		if c, arg, ok := strings.Cut(msg.Value, " "); ok {
			if c == "" {
				return m, m.NewStatusMessage(status)
			}
			for _, co := range m.Commands {
				if co.Name == c {
					return m, co.Cmd(arg)
				}
			}
		}
		return m, m.NewStatusMessage(status)
	case statusMessageTimeoutMsg:
		m.hideStatusMessage()
	case statusMsg:
		return m, m.NewStatusMessage(msg.Value)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case key.Matches(msg, m.KeyMap.Command):
			if !m.showFilter() && m.ShowCommand() {
				m.hideStatusMessage()
				cmds = append(cmds, m.SetFocus(Cmd))
			}
		}
		if m.ShowCommand() {
			for _, c := range m.Commands {
				if key.Matches(msg, c.Key) {
					m.Command.SetValue(c.Name + " ")
					cmds = append(cmds, m.SetFocus(Cmd))
				}
			}
		}
	}

	switch m.State() {
	case Cmd:
		if m.ShowCommand() {
			cmd = m.updateCommand(msg)
			cmds = append(cmds, cmd)
		}
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

	//cmds = append(cmds, m.NewStatusMessage(m.state.String()))

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
		case key.Matches(msg, m.List.KeyMap.Filter):
			m.hideStatusMessage()
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

	if m.showItemDesc && m.HasPager() {
		item := m.List.CurrentItem()
		m.Pager.SetText(item.Description())
	}

	return tea.Batch(cmds...)
}

func (m *Model) updateCommand(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg.(type) {
	case input.UnfocusMsg:
		switch {
		case m.HasPager():
			cmds = append(cmds, m.SetFocus(Pager))
		case m.HasList():
			cmds = append(cmds, m.SetFocus(List))
		}
	case tea.KeyMsg:
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
		if m.HasList() {
			m.List.Unfocus()
		}
		m.Command.Unfocus()
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

func (m Model) ShowCommand() bool {
	return m.showCommand
}

func (m *Model) SetShowCommand(show bool) {
	m.showCommand = show
}

func (m Model) HasPager() bool {
	return m.Pager != nil
}

func (m Model) ShowPager() bool {
	if m.HasPager() {
		return m.showPager
	}
	return false
}

func (m *Model) SetShowPager(show bool) {
	m.showPager = show
}

func (m Model) ShowList() bool {
	if m.HasList() {
		return m.showList
	}
	return false
}

func (m *Model) SetShowList(show bool) {
	m.showList = show
}

func (m Model) HasList() bool {
	return m.List != nil
}

func (m Model) showFilter() bool {
	if m.HasList() {
		return m.List.SettingFilter()
	}
	return false
}

func (m *Model) View() string {
	var sections []string

	var li string
	if m.ShowList() {
		li = m.List.View()
	}

	var page string
	if m.ShowPager() {
		page = m.Pager.View()
	}

	main := m.layout.Join(li, page)
	sections = append(sections, main)

	var footer string
	if !m.showFilter() && !m.Command.Focused() {
		footer = m.statusMessage
	}
	switch m.State() {
	case Cmd:
		if m.statusMessage == "" {
			footer = m.Command.View()
		}
	case List:
		if m.ShowList() {
			switch {
			case m.showFilter():
				footer = m.List.FilterInput.View()
			case m.List.State() == list.Input:
				footer = m.List.Input.View()
			}
		}
	}
	sections = append(sections, footer)

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

func TermHeight() int {
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return h
}

func TermWidth() int {
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))
	return w
}
