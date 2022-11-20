package teacozy

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TUI struct {
	Main            *List
	Alt             *List
	Input           textarea.Model
	view            viewport.Model
	prompt          textinput.Model
	info            *Info
	Title           string
	FocusedView     string
	showMenu        bool
	showInfo        bool
	showHelp        bool
	currentListItem int
	Style           TUIStyle
	width           int
	height          int
	fullScreen      bool
	Hash            map[string]string
	ShortHelp       Help
	Help            *Info
	MainMenu        *Menu
	ActionMenu      *Menu
	Menus           Menus
	CurrentMenu     *Menu
}

func New(main *List) TUI {
	ui := TUI{
		Main:        main,
		Menus:       make(Menus),
		FocusedView: "list",
		Style:       DefaultTuiStyle(),
		MainMenu:    DefaultMenu().SetToggle("m", "menu").SetLabel("menu"),
		ActionMenu:  ActionMenu(),
		showHelp:    true,
	}
	ui.SetHelp(Keys.SortList, Keys.Menu, Keys.Help)
	ui.AddMenu(SortListMenu())
	return ui
}

func (ui *TUI) SetSize(w, h int) *TUI {
	ui.Main.SetSize(w, h)
	return ui
}

func (ui TUI) Width() int {
	return ui.Style.Frame.Width()
}

func (ui TUI) Height() int {
	return ui.Style.Frame.Height()
}

func (ui *TUI) SetHelp(keys ...Key) *TUI {
	ui.ShortHelp = NewHelp(keys...)
	return ui
}

func (l *TUI) AddMenu(menu *Menu) {
	l.MainMenu.NewKey(
		menu.Toggle.Help().Key,
		menu.Toggle.Help().Desc,
		GoToMenuCmd(menu),
	)
	l.Menus[menu.Label] = menu
}

func (l *TUI) ShowMenu() {
	l.showMenu = true
}

func (l *TUI) HideMenu() {
	l.showMenu = false
}

func (ui *TUI) ShowInfo() {
	ui.showInfo = true
}

func (l *TUI) HideInfo() {
	l.showInfo = false
}

func (m *TUI) Start() *TUI {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		focus = m.FocusedView
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		if m.Main.Editable {
			switch {
			case Keys.ExitScreen.Matches(msg):
				m.Main = m.Alt
			case Keys.PrevScreen.Matches(msg):
				m.Main = m.Alt
			case Keys.SaveAndExit.Matches(msg):
				m.HideInfo()
			}
		}
		if m.showMenu {
			cmds = append(cmds, m.UpdateMenu(msg))
		}
		switch {
		case Keys.PrevScreen.Matches(msg):
			if m.Main.SelectionList {
				m.Main.SelectionList = false
				cmds = append(cmds, m.Main.ShowVisibleItemsCmd())
			}
		case Keys.FullScreen.Matches(msg):
			m.fullScreen = !m.fullScreen
			cmds = append(cmds, m.ToggleFullScreenCmd())
		case Keys.Info.Matches(msg):
			cmds = append(cmds, HideInfoCmd())
		case Keys.Quit.Matches(msg):
			cmds = append(cmds, tea.Quit)
		case Keys.Help.Matches(msg):
			if focus == "help" {
				cmds = append(cmds, HideInfoCmd())
			} else {
				m.info = Keys.FullHelp()
				cmds = append(cmds, ShowInfoCmd())
			}
		case Keys.Menu.Matches(msg):
			if focus == "menu" {
				cmds = append(cmds, HideMenuCmd())
			} else {
				cmds = append(cmds, ChangeMenuCmd(m.MainMenu))
			}
		default:
			for label, menu := range m.Menus {
				if key.Matches(msg, menu.Toggle) && len(menu.Items) > 0 {
					m.CurrentMenu = menu
					m.ShowMenu()
					m.HideInfo()
					cmds = append(cmds, SetFocusedViewCmd(label))
				}
			}
		}
	case tea.WindowSizeMsg:
		w := msg.Width - 1
		h := msg.Height - 2
		m.Style.Frame.SetSize(w, h)
		m.SetSize(w, h)
	case EditInfoMsg:
		cur := m.Main.SelectedItem()
		m.Alt = m.Main
		m.Main = cur.Fields.Edit()
		cmds = append(cmds, HideInfoCmd())
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case ShowItemInfoMsg:
		m.info = NewForm().SetData(msg.Fields)
		m.currentListItem = m.Main.Model.Index()
		cmds = append(cmds, ShowInfoCmd())
	case HideInfoMsg:
		m.HideInfo()
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case ShowInfoMsg:
		m.ShowInfo()
		m.HideMenu()
		cmds = append(cmds, SetFocusedViewCmd("info"))
	case ChangeMenuMsg:
		m.CurrentMenu = msg.Menu
		cmds = append(cmds, ShowMenuCmd(m.CurrentMenu))
	case ShowMenuMsg:
		m.ShowMenu()
		m.HideInfo()
		cmds = append(cmds, SetFocusedViewCmd(m.CurrentMenu.Label))
	case HideMenuMsg:
		m.HideMenu()
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
	case ActionMenuMsg:
		m.CurrentMenu = m.ActionMenu
		m.ShowMenu()
		m.HideInfo()
		cmds = append(cmds, SetFocusedViewCmd("action"))
	case FormChangedMsg:
		m.Main = m.Alt
		m.Main.Model.Select(m.currentListItem)
		cur := m.Main.SelectedItem()
		m.info = NewForm().SetData(cur.Fields)
		m.ShowInfo()
		cmds = append(cmds, ItemChangedCmd(cur))
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case SaveAndExitFormMsg:
		cmds = append(cmds, msg.Save(m.Main))
	}

	switch focus {
	case "info":
		m.info, cmd = m.info.Update(msg)
		cmds = append(cmds, cmd)
	case "list":
		if m.Main.SelectionList {
			cmds = append(cmds, ActionMenuCmd())
		}
		m.Main, cmd = m.Main.Update(msg)
		cmds = append(cmds, cmd)
	default:
		for label, _ := range m.Menus {
			if focus == label {
				cmds = append(cmds, m.UpdateMenu(msg))
			}
		}
	}

	//cmds = append(cmds, UpdateStatusCmd(m.FocusedView))
	return m, tea.Batch(cmds...)
}

func (m *TUI) UpdateMenu(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.CurrentMenu.Toggle):
			m.HideMenu()
			cmds = append(cmds, SetFocusedViewCmd("list"))
		default:
			for _, item := range m.CurrentMenu.Items {
				if key.Matches(msg, item.Bind) {
					m.HideMenu()
					cmds = append(cmds, item.Cmd(m))
				}
			}
			m.HideMenu()
			cmds = append(cmds, SetFocusedViewCmd("list"))
		}
	}
	m.CurrentMenu.Model, cmd = m.CurrentMenu.Model.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *TUI) Init() tea.Cmd {
	return nil
}

func (m *TUI) View() string {
	var (
		sections     []string
		availHeight  = m.Height()
		widgetHeight = m.Style.Widget.Height()
	)

	m.SetSize(m.Width(), availHeight)

	var widget string
	if m.showMenu {
		widget = m.CurrentMenu.View()
		availHeight -= widgetHeight
	}

	if m.showInfo {
		widget = m.info.View()
		availHeight -= widgetHeight
	}

	var help string
	if m.showHelp {
		help = m.ShortHelp.View()
		availHeight -= lipgloss.Height(help)
	}

	content := m.Main.View()
	sections = append(sections, content)

	if m.showMenu {
		sections = append(sections, widget)
	}

	if m.showInfo {
		sections = append(sections, widget)
	}

	if m.showHelp {
		sections = append(sections, help)
	}

	if min := m.Main.Frame.MinHeight; min > availHeight {
		if m.showMenu || m.showInfo {
			return lipgloss.NewStyle().Height(availHeight).Render(widget)
		}
		return lipgloss.NewStyle().Height(availHeight).Render(content)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

const (
	main state = iota
	form
	view
	edit
)

type state int

func (s state) String() string {
	switch s {
	case form:
		return "form"
	case view:
		return "view"
	case edit:
		return "edit"
	default:
		return "main"
	}
}
