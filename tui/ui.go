package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Tui struct {
	state        state
	Main         tea.Model
	info         viewport.Model
	showInfo     bool
	Info         *info.Info
	Help         Help
	showFullHelp bool
	Style        Style
	width        int
	height       int
}

func NewTui(main *list.List) Tui {
	ui := Tui{
		Main:  main,
		state: mainModel,
		Style: DefaultStyle(),
		Help:  NewHelp(),
		Info:  info.New(),
	}
	ui.info = viewport.New(ui.Style.Widget.Width(), ui.Style.Widget.Height())
	return ui
}

func (m Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case info.HideInfoMsg:
		m.showInfo = false
		m.Info.Hide()
		m.state = mainModel
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
			switch {
			case key.Matches(msg, key.HelpKey):
				m.state = infoModel
				m.Info = m.Help.Info
				m.showInfo = true
				m.showFullHelp = true
			}
			switch main := m.Main.(type) {
			case *list.List:
				if main.SelectionList {
					cmds = append(cmds, ActionMenuCmd())
				}
				m.Main, cmd = m.Main.Update(msg)
				cmds = append(cmds, cmd)
			}
		case infoModel:
			m.Info, cmd = m.Info.Update(msg)
			cmds = append(cmds, cmd)
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
		case infoModel:
			m.Info, cmd = m.Info.Update(msg)
			cmds = append(cmds, cmd)
		}

	}

	return m, tea.Batch(cmds...)
}

func (m Tui) View() string {
	var (
		sections     []string
		availHeight  = m.Height()
		widgetWidth  = m.Style.Widget.Width()
		widgetHeight = m.Style.Widget.Height()
	)
	m.SetSize(m.Width(), availHeight)

	var widget string
	if m.showInfo {
		m.Info.SetSize(widgetWidth, widgetHeight)
		if m.showFullHelp {
			//m.info = viewport.New(m.Width(), m.Height())
			//m.Help.Info.Model.SetContent(m.Help.Render())
			//m.Help.Render()
			widget = m.Info.View()
		}
		availHeight -= widgetHeight
	}

	content := m.Main.View()
	sections = append(sections, content)

	if m.showFullHelp {
		sections = append(sections, widget)
	}

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
