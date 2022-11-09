package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/item"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/prompt"
)

type UI struct {
	*prompt.Model
	info             *info.Fields
	Keys             urkey.KeyMap
	Title            string
	IsMultiSelect    bool
	ShowSelectedOnly bool
	FocusedView      string
	ShowWidget       bool
	showMenu         bool
	showInfo         bool
	width            int
	height           int
	Menus            Menus
	CurrentMenu      *Menu
}

func NewUI(title string) UI {
	return UI{
		Model:       prompt.New(),
		Title:       title,
		Menus:       make(Menus),
		FocusedView: "list",
	}
}

func (l *UI) AddMenu(menu *Menu) {
	l.Menus[menu.Label] = menu
}

func (l *UI) ShowMenu() {
	l.showMenu = true
}

func (l *UI) HideMenu() {
	l.showMenu = false
}

func (l *UI) ShowInfo() {
	l.showInfo = true
}

func (l *UI) HideInfo() {
	l.showInfo = false
}

func (m *UI) Start() *UI {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		switch focus {
		case "info":
			if key.Matches(msg, urkey.Quit) {
				m.HideInfo()
				cmds = append(cmds, SetFocusedViewCmd("list"))
			}
		}
		switch {
		case key.Matches(msg, urkey.Quit):
			cmds = append(cmds, tea.Quit)
		default:
			for label, menu := range m.Menus {
				if key.Matches(msg, menu.Toggle) {
					m.CurrentMenu = menu
					m.ShowMenu()
					cmds = append(cmds, SetFocusedViewCmd(label))
				}
			}
		}
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
	case item.ShowInfoMsg:
		m.info = msg.Info
		m.ShowInfo()
		cmds = append(cmds, SetFocusedViewCmd("info"))
	}

	switch focus {
	case "info":
		m.info, cmd = m.info.Update(msg)
		cmds = append(cmds, cmd)
	case "list":
		//l, cmd := m.List.Update(msg)
		//m.List = l.(List)
		//cmds = append(cmds, cmd)
		m.Model, cmd = m.Model.Update(msg)
		cmds = append(cmds, cmd)
		//cmds = append(cmds, UpdateList(m, msg))
	default:
		for label, _ := range m.Menus {
			if focus == label {
				cmds = append(cmds, UpdateMenu(m, msg))
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *UI) View() string {
	var (
		sections    []string
		availHeight = m.Model.List.Height()
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

	m.Model.SetSize(m.width, availHeight)
	content := m.Model.List.View()
	sections = append(sections, content)

	if m.showMenu {
		sections = append(sections, menu)
	}

	if m.showInfo {
		sections = append(sections, info)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
