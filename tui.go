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
	Main             *List
	Alt              *List
	Input            textarea.Model
	view             viewport.Model
	info             *Info
	Title            string
	FocusedView      string
	ShowWidget       bool
	showMenu         bool
	showInfo         bool
	currentModelItem int
	widgetStyle      WidgetStyle
	width            int
	height           int
	state            state
	Hash             map[string]string
	Menus            Menus
	CurrentMenu      *Menu
}

func New(title string, items Items) TUI {
	return TUI{
		Main:        NewList(title, items),
		Title:       title,
		Menus:       make(Menus),
		FocusedView: "list",
		widgetStyle: WidgetStyle{
			MaxHeight: TermHeight() / 3,
		},
	}
}

func (ui *TUI) SetSize(w, h int) *TUI {
	ui.width = w
	ui.height = h
	ui.Main.SetSize(w, h)
	return ui
}

func (l *TUI) AddMenu(menu *Menu) {
	l.Menus[menu.Label] = menu
}

func (l *TUI) ShowMenu() {
	l.showMenu = true
}

func (l *TUI) HideMenu() {
	l.showMenu = false
}

func (l *TUI) ShowInfo() {
	l.showInfo = true
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
				//m.Main = m.Alt
				m.HideInfo()
			}
		}
		switch {
		case key.Matches(msg, Keys.Info):
			cmds = append(cmds, HideInfoCmd())
		case key.Matches(msg, Keys.Quit):
			cmds = append(cmds, tea.Quit)
		default:
			for label, menu := range m.Menus {
				if key.Matches(msg, menu.Toggle) {
					m.CurrentMenu = menu
					m.ShowMenu()
					m.HideInfo()
					cmds = append(cmds, SetFocusedViewCmd(label))
				}
			}
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width-1, msg.Height-2)
	case EditInfoMsg:
		cur := m.Main.SelectedItem()
		m.Alt = m.Main
		m.Main = cur.Fields.Edit()
		cmds = append(cmds, HideInfoCmd())
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case ShowItemInfoMsg:
		m.info = NewInfo(msg.Fields).SetSize(m.widgetStyle.Width(), m.widgetStyle.Height())
		m.currentModelItem = m.Main.Model.Index()
		cmds = append(cmds, ShowInfoCmd())
	case HideInfoMsg:
		m.HideInfo()
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case ShowInfoMsg:
		m.ShowInfo()
		m.HideMenu()
		cmds = append(cmds, SetFocusedViewCmd("info"))
	case HideMenuMsg:
		m.HideMenu()
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case ShowMenuMsg:
		m.ShowMenu()
		m.HideInfo()
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
	case FormChangedMsg:
		m.Main = m.Alt
		m.Main.Model.Select(m.currentModelItem)
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
				cmds = append(cmds, UpdateMenu(m, msg))
			}
		}
	}

	//cmds = append(cmds, UpdateStatusCmd(m.FocusedView))
	return m, tea.Batch(cmds...)
}

func (m *TUI) Init() tea.Cmd {
	return nil
}

func (m *TUI) View() string {
	var (
		sections    []string
		availHeight = m.height
	)

	var menu string
	if m.showMenu {
		menu = m.CurrentMenu.Model.View()
		availHeight -= lipgloss.Height(menu)
	}

	var info string
	if m.showInfo {
		info = m.info.View()
		availHeight -= lipgloss.Height(info)
	}

	m.SetSize(m.width, availHeight)
	content := m.Main.View()
	sections = append(sections, content)

	if m.showMenu {
		sections = append(sections, menu)
	}

	if m.showInfo {
		sections = append(sections, info)
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
