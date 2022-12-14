package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/form"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type state int

const (
	mainModel state = iota
	formModel
	infoModel
	helpModel
	viewModel
	menuModel
)

type Tui struct {
	width        int
	height       int
	state        state
	Style        Style
	fullScreen   bool
	KeyMap       keyMap
	Main         tea.Model
	Alt          tea.Model
	view         viewport.Model
	widget       viewport.Model
	showInfo     bool
	Info         *info.Info
	Help         Help
	showFullHelp bool
	showMenu     bool
	MainMenu     *Menu
	Menus        Menus
	CurrentMenu  *Menu
	ActionMenu   *Menu
	Form         *form.Form
}

func NewTui(main *list.List) Tui {
	main.AddUpdateFunc("i", "item meta", ShowItemMeta)
	ui := Tui{
		Main:       main,
		KeyMap:     DefaultKeyMap(),
		state:      mainModel,
		Style:      DefaultStyle(),
		Help:       NewHelp(),
		Info:       info.New(),
		Menus:      make(Menus),
		MainMenu:   MainMenu(),
		ActionMenu: ActionMenu(),
	}
	ui.view = viewport.New(ui.Width(), ui.Height())
	ui.widget = viewport.New(ui.Style.Widget.Width(), ui.Style.Widget.Height())
	ui.CurrentMenu = ui.MainMenu
	return ui
}

func (m Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case form.FormChangedMsg:
		m.Main = m.Alt
	case EditItemMetaMsg:
		m.Alt = m.Main
		m.Main = &msg.Item.Meta
		cmds = append(cmds, info.HideInfoCmd())
		m.state = mainModel
	case ShowItemMetaMsg:
		//m.Info = msg.Item.Meta.Info
		m.CurrentMenu = EditItemMetaMenu(msg.Item)
		//m.showInfo = true
		m.showMenu = true
		content := msg.Item.Meta.Info.Render()
		cmds = append(cmds, info.UpdateContentCmd(content))
		m.state = menuModel
		//m.state = formModel
	case ChangeMenuMsg:
		m.CurrentMenu = msg.Menu
	case info.HideInfoMsg:
		m.showInfo = false
		m.showFullHelp = false
		m.showMenu = false
		m.state = mainModel
	case info.UpdateContentMsg:
		m.view.SetContent(msg.Content)
		m.widget.SetContent(msg.Content)
	case tea.WindowSizeMsg:
		w := msg.Width - 1
		h := msg.Height - 2
		m.Style.Frame.SetSize(w, h)
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		switch m.state {
		case formModel:
			if m.showInfo {
				cmd = m.updateInfo(msg, m.Info)
				cmds = append(cmds, cmd)
			} else {
				switch {
				case key.Matches(msg, key.EditField):
					m.showInfo = false
				default:
					var model tea.Model
					switch main := m.Alt.(type) {
					case *form.Form:
						model, cmd = main.Update(msg)
						cmds = append(cmds, cmd)
					}
					m.Alt = model
				}
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
		case mainModel:
			switch {
			case key.Matches(msg, key.MenuKey):
				m.state = menuModel
				m.showMenu = true
				m.CurrentMenu = m.MainMenu
				cmds = append(cmds, info.UpdateContentCmd(m.CurrentMenu.Render()))
			case key.Matches(msg, key.HelpKey):
				cmd = GoToHelpView(&m)
				cmds = append(cmds, cmd)
			case key.Matches(msg, key.NewBinding("a", "action")):
				m.state = menuModel
				m.showMenu = true
				m.CurrentMenu = m.ActionMenu
				cmds = append(cmds, info.UpdateContentCmd(m.CurrentMenu.Render()))
			default:
				var model tea.Model
				switch main := m.Main.(type) {
				case *list.List:
					if main.SelectionList {
						cmds = append(cmds, ActionMenuCmd())
					}
					model, cmd = main.Update(msg)
					cmds = append(cmds, cmd)
				case *form.Form:
					model, cmd = main.Update(msg)
					cmds = append(cmds, cmd)
				}
				m.Main = model
			}
		}
	default:
		switch m.state {
		case mainModel:
			var model tea.Model
			switch main := m.Main.(type) {
			case *list.List:
				model, cmd = main.Update(msg)
				cmds = append(cmds, cmd)
			case *form.Form:
				model, cmd = main.Update(msg)
				cmds = append(cmds, cmd)
			}
			m.Main = model
		case infoModel:
			cmd = m.updateInfo(msg, m.Info)
			cmds = append(cmds, cmd)
		case helpModel:
			cmd = m.updateInfo(msg, m.Help.Info)
			cmds = append(cmds, cmd)
		case menuModel:
			cmd = m.updateMenu(msg)
			cmds = append(cmds, cmd)
		case formModel:
			if m.showInfo {
				cmd = m.updateInfo(msg, m.Info)
				cmds = append(cmds, cmd)
			} else {
			}
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
		widget = m.widget.View()
		availHeight -= lipgloss.Height(widget)
	}

	if m.showMenu {
		widget = m.widget.View()
		availHeight -= lipgloss.Height(widget)
	}

	m.SetSize(m.Width(), availHeight)
	var content string
	switch m.state {
	//case formModel:
	//content = m.Alt.View()
	default:
		content = m.Main.View()
	}
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
