package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/list"
)

type Tui struct {
	Main tea.Model
	//Info   *info.Info
	Style  Style
	state  state
	width  int
	height int
}

func NewTui(main *list.List) Tui {
	ui := Tui{
		Main:  main,
		state: mainModel,
		Style: DefaultStyle(),
		//Info:  info.New(),
	}
	return ui
}

func (m Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w := msg.Width - 1
		h := msg.Height - 2
		m.Style.Frame.SetSize(w, h)
		//m.SetSize(w, h)
	case tea.KeyMsg:
		//cmds = append(cmds, list.UpdateStatusCmd(msg.String()))
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		switch m.state {
		case mainModel:
			switch main := m.Main.(type) {
			case *list.List:
				if main.SelectionList {
					cmds = append(cmds, ActionMenuCmd())
				}
				m.Main, cmd = m.Main.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	default:
		switch m.state {
		case mainModel:
			switch main := m.Main.(type) {
			case *list.List:
				if main.SelectionList {
					cmds = append(cmds, ActionMenuCmd())
				}
				m.Main, cmd = m.Main.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func updateList(msg tea.Msg, m *list.List) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	var li tea.Model
	li, cmd = m.Update(msg)
	cmds = append(cmds, cmd)

	//switch msg := msg.(type) {
	//case tea.KeyMsg:
	//}

	return li, tea.Batch(cmds...)
}

func (m Tui) View() string {
	var (
		sections    []string
		availHeight = m.Height()
		//widgetWidth  = m.Style.Widget.Width()
		//widgetHeight = m.Style.Widget.Height()
	)
	m.SetSize(m.Width(), availHeight)

	content := m.Main.View()
	sections = append(sections, content)

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m Tui) Init() tea.Cmd {
	return nil
}

func (m Tui) Start() Tui {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m
}

func (ui *Tui) SetSize(w, h int) *Tui {
	switch main := ui.Main.(type) {
	case *list.List:
		main.SetSize(w, h)
	}
	return ui
}

func (ui Tui) Width() int {
	return ui.Style.Frame.Width()
}

func (ui Tui) Height() int {
	return ui.Style.Frame.Height()
}
