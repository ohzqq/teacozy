package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Tui struct {
	width        int
	height       int
	state        state
	Style        Style
	KeyMap       keyMap
	Main         tea.Model
	view         viewport.Model
	showInfo     bool
	Info         *info.Info
	info         *Info
	Help         Help
	showFullHelp bool
	showMenu     bool
	MainMenu     *Menu
	Menus        Menus
	CurrentMenu  *Menu
}

func NewTui(main *list.List) Tui {
	mk := key.NewKey("m", "menu")
	ui := Tui{
		Main:     main,
		KeyMap:   DefaultKeyMap(),
		state:    mainModel,
		Style:    DefaultStyle(),
		Help:     NewHelp(),
		Info:     info.New(),
		Menus:    make(Menus),
		MainMenu: NewMenu(mk),
	}
	ui.view = viewport.New(ui.Style.Widget.Width(), ui.Style.Widget.Height())
	ui.MainMenu.AddKey(ui.Help.Toggle, GoToHelp)
	ui.CurrentMenu = ui.MainMenu
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
		m.showFullHelp = false
		m.state = mainModel
	case info.UpdateContentMsg:
		m.view.SetContent(msg.Content)
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
			cmd = m.updateInfo(msg, m.Info)
			cmds = append(cmds, cmd)
		case helpModel:
			cmd = m.updateInfo(msg, m.Help.Info)
			cmds = append(cmds, cmd)
		case menuModel:
			cmd = m.updateMenu(msg)
			cmds = append(cmds, cmd)
		case mainModel:
			switch {
			case key.Matches(msg, key.MenuKey):
				m.state = menuModel
				m.showMenu = true
				m.CurrentMenu = m.MainMenu
				cmds = append(cmds, info.UpdateContentCmd(m.CurrentMenu.Render()))
			case key.Matches(msg, key.HelpKey):
				m.state = helpModel
				m.showFullHelp = true
				m.view = m.Help.Info.Model
				cmds = append(cmds, info.UpdateContentCmd(m.Help.Render()))
			default:
				switch main := m.Main.(type) {
				case *list.List:
					if main.SelectionList {
						cmds = append(cmds, ActionMenuCmd())
					}
					m.Main, cmd = m.Main.Update(msg)
					//cmds = append(cmds, list.UpdateStatusCmd("list"))
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
				//cmds = append(cmds, list.UpdateStatusCmd("list"))
				cmds = append(cmds, cmd)
			}
		case infoModel:
			cmd = m.updateInfo(msg, m.Info)
			cmds = append(cmds, cmd)
		case helpModel:
			cmd = m.updateInfo(msg, m.Help.Info)
			cmds = append(cmds, cmd)
		case menuModel:
			cmd = m.updateMenu(msg)
			cmds = append(cmds, cmd)
		}

	}

	return m, tea.Batch(cmds...)
}

func (m Tui) View() string {
	var (
		sections    []string
		availHeight = m.Height()
		//widgetWidth  = m.Style.Widget.Width()
		//widgetHeight = m.Style.Widget.Height()
	)

	if m.showFullHelp {
		return m.view.View()
	}

	var widget string
	if m.showInfo {
		widget = m.view.View()
		availHeight -= lipgloss.Height(widget)
	}

	if m.showMenu {
		widget = m.view.View()
		availHeight -= lipgloss.Height(widget)
		fmt.Println(widget)
	}

	m.SetSize(m.Width(), availHeight)
	content := m.Main.View()
	sections = append(sections, content)

	if m.showInfo || m.showMenu {
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
