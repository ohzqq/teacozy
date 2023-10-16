package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/input"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
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

	//help
	Help     help.Model
	showHelp bool

	StatusMsg
}

func New(opts ...Option) *Model {
	m := &Model{
		KeyMap:   DefaultKeyMap(),
		mainView: Cmd,
		Help:     help.New(),
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

	if m.layout == nil {
		m.layout = NewLayout()
		if m.HasList() && m.HasPager() {
			m.layout.Half()
			m.layout.Position(Top)
		}
	}
	if m.HasPager() {
		m.mainView = Pager
		m.SetShowHelp(true)
	}
	if m.HasList() {
		m.mainView = List
		m.List.SetShowInput(false)
		m.SetShowHelp(true)
	}

	m.SetSize(teacozy.TermSize())

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

	mw, mh := m.layout.main()
	sw, sh := m.layout.sub()
	if !m.Help.ShowAll {
		mh--
	} else {
		mh++
		sh++
	}

	switch m.mainView {
	case List:
		if m.HasList() {
			m.List.SetSize(mw, mh)
		}
		if m.HasPager() {
			m.Pager.SetSize(sw, sh)
		}
	case Pager:
		if m.HasPager() {
			m.Pager.SetSize(mw, mh)
		}
		if m.HasList() {
			m.List.SetSize(sw, sh)
		}
	}

	m.Command.Width = m.Width()
	m.Help.Width = m.Width()
}

func (m *Model) AddCommands(cmds ...Command) *Model {
	m.SetShowCommand(true)
	m.Commands = append(m.Commands, cmds...)
	return m
}

func (m *Model) SetList(parser list.ParseItems, opts ...list.Option) *Model {
	m.showList = true
	items := list.NewItems(parser)
	l := list.New(items, opts...)
	l.SetShowFilter(false)
	l.SetShowInput(false)
	l.SetShowHelp(false)
	m.List = l
	return m
}

func (m *Model) SetPager(render pager.Renderer, text ...string) *Model {
	m.showPager = true
	p := pager.New(render)
	if len(text) > 0 {
		p.SetText(text[0])
	}
	p.SetContent(p.Render())
	m.Pager = p
	return m
}

func (m *Model) SetLayout(l *Layout) *Model {
	m.layout = l
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
		case key.Matches(msg, m.KeyMap.ToggleFocus):
			if m.Help.ShowAll {
				m.toggleHelp()
			}
			return m, m.toggleView()
		case key.Matches(msg, m.KeyMap.FullHelp):
			m.toggleHelp()
			cmds = append(cmds, tea.ClearScreen)
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
		cmds = append(cmds, m.focusMain())
	case tea.KeyMsg:
	}

	m.Command, cmd = m.Command.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) focusMain() tea.Cmd {
	return m.SetFocus(m.mainView)
}

func (m *Model) toggleView() tea.Cmd {
	switch m.State() {
	case List:
		if m.HasPager() {
			return m.SetFocus(Pager)
		}
	case Pager:
		if m.HasList() {
			return m.SetFocus(List)
		}
	}
	return nil
}

func (m *Model) toggleHelp() {
	m.Help.ShowAll = !m.Help.ShowAll
	w, h := teacozy.TermSize()
	if m.Help.ShowAll {
		m.SetSize(w, h-m.helpViewHeight())
	} else {
		m.SetSize(w, h)
	}
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
			cmds = append(cmds, input.Unfocus)
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
		cmds = append(cmds, input.Focus)
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

func (m Model) ShowHelp() bool {
	return m.showHelp
}

func (m *Model) SetShowHelp(show bool) {
	m.showHelp = show
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

func (m Model) helpView() string {
	return m.Help.View(m)
}

func (m Model) helpViewHeight() int {
	return lipgloss.Height(m.helpView())
}

func (m *Model) View() string {
	var sections []string
	//availHeight := m.Height()

	var li string
	if m.ShowList() {
		li = m.List.View()
	}

	var page string
	if m.ShowPager() {
		page = m.Pager.View()
	}

	//var help string
	//if m.ShowHelp() {
	//help = m.Help.View(m)
	//availHeight -= lipgloss.Height(help)
	//}

	main := m.layout.Join(li, page)
	//main = lipgloss.NewStyle().
	//Height(availHeight).
	//Render(main)
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

	if m.ShowHelp() {
		help := m.Help.View(m)
		sections = append(sections, help)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) ShortHelp() []key.Binding {
	var h []key.Binding

	switch m.State() {
	case List:
		h = append(h, m.List.ShortHelp()...)
	case Pager:
		h = append(h, m.Pager.ShortHelp()...)
	}

	h = append(h, m.appHelp()...)
	h = append(h, m.KeyMap.Quit)
	return h
}

func (m Model) appHelp() []key.Binding {
	var h []key.Binding

	if m.HasList() && m.HasPager() {
		h = append(h, m.KeyMap.ToggleFocus)
	}

	if m.ShowCommand() {
		h = append(h, m.KeyMap.Command)
	}

	if m.ShowHelp() {
		h = append(h, m.KeyMap.FullHelp)
	}

	return h
}

func (m Model) FullHelp() [][]key.Binding {
	var h [][]key.Binding
	h = append(h, m.appHelp())

	if m.HasList() {
		for _, help := range m.List.FullHelp() {
			h = append(h, help)
		}
	}

	if m.HasPager() {
		h = append(h, m.Pager.FullHelp()...)
	}

	return h
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

var NoItems = list.NewItems(func() []*list.Item { return []*list.Item{} })
