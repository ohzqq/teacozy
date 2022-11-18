package teacozy

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TUI struct {
	Main            *List
	Alt             *List
	Input           textarea.Model
	view            viewport.Model
	info            *Info
	Title           string
	FocusedView     string
	showMenu        bool
	showInfo        bool
	currentListItem int
	Style           TUIStyle
	width           int
	height          int
	state           state
	Hash            map[string]string
	HelpMenu        *Menu
	Menus           Menus
	CurrentMenu     *Menu
}

func New(title string, items Items) TUI {
	ui := TUI{
		Main:        NewList(title, items),
		Title:       title,
		Menus:       make(Menus),
		FocusedView: "list",
		Style:       DefaultTuiStyle(),
		HelpMenu:    DefaultMenu().SetToggle("?", "help").SetLabel("help"),
	}
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

func (l *TUI) AddMenu(menu *Menu) {
	l.HelpMenu.NewKey(
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
		if m.Main.isForm {
			switch {
			case key.Matches(msg, Keys.ExitScreen):
				m.Main = m.Alt
			case key.Matches(msg, Keys.SaveAndExit):
				m.HideInfo()
			}
		}
		if m.showMenu {
			cmds = append(cmds, m.UpdateMenu(msg))
		}
		switch {
		case key.Matches(msg, Keys.Info):
			cmds = append(cmds, HideInfoCmd())
		case key.Matches(msg, Keys.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, Keys.Help):
			if focus == "help" {
				cmds = append(cmds, HideMenuCmd())
			} else {
				cmds = append(cmds, ChangeMenuCmd(m.HelpMenu))
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
		m.info = NewInfo().SetData(msg.Fields)
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
	case FormChangedMsg:
		m.Main = m.Alt
		m.Main.Model.Select(m.currentListItem)
		cur := m.Main.SelectedItem()
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
				if key.Matches(msg, item.KeyBind) {
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

	content := m.Main.View()
	sections = append(sections, content)

	if m.showMenu {
		sections = append(sections, widget)
	}

	if m.showInfo {
		sections = append(sections, widget)
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
