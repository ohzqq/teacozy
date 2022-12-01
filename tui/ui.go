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
	view         viewport.Model
	showInfo     bool
	Info         *info.Info
	info         *Info
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
		info:  NewInfo(),
	}
	ui.view = viewport.New(ui.Style.Widget.Width(), ui.Style.Widget.Height())
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
		case infoModel:
			//return updateInfo(msg, m)
			m.view, cmd = updateInfo(msg, m)
			cmds = append(cmds, cmd)
		case mainModel:
			switch {
			case key.Matches(msg, key.HelpKey):
				m.state = infoModel
				m.Info = m.Help.Info
				m.showInfo = true
				m.showFullHelp = true
				cmds = append(cmds, list.UpdateStatusCmd("info"))
			default:
				switch main := m.Main.(type) {
				case *list.List:
					if main.SelectionList {
						cmds = append(cmds, ActionMenuCmd())
					}
					m.Main, cmd = m.Main.Update(msg)
					cmds = append(cmds, list.UpdateStatusCmd("list"))
					cmds = append(cmds, cmd)
				}
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
				cmds = append(cmds, list.UpdateStatusCmd("list"))
				cmds = append(cmds, cmd)
			}
		case infoModel:
			m.view, cmd = updateInfo(msg, m)
			cmds = append(cmds, cmd)
			//return updateInfo(msg, m)
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
			//m.view = viewport.New(m.Width(), m.Height())
			//m.Help.Info.Model.SetContent(m.Help.Render())
			//m.Help.Render()
			widget = m.view.View()
		}
		availHeight -= widgetHeight
	}

	content := m.Main.View()
	sections = append(sections, content)

	//if m.showFullHelp {
	//  sections = append(sections, widget)
	//}

	if m.showInfo {
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
