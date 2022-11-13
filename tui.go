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
	*List
	Input       textarea.Model
	view        viewport.Model
	form        *List
	info        *Info
	Title       string
	FocusedView string
	ShowWidget  bool
	showMenu    bool
	showInfo    bool
	width       int
	height      int
	state       state
	Hash        map[string]string
	Menus       Menus
	CurrentMenu *Menu
}

func New(title string, items Items) TUI {
	return TUI{
		List:        NewList(title, items),
		Title:       title,
		Menus:       make(Menus),
		FocusedView: "list",
		state:       main,
	}
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
		if m.Input.Focused() {
			if key.Matches(msg, Keys.SaveAndExit) {
				cur := m.List.Model.SelectedItem()
				i := m.Items.Get(cur)
				field := i.Data.(FieldData)
				val := m.Input.Value()
				field.Set(val)
				m.Items.Set(i.Index(), NewItem(field))
				m.Input.Blur()
				m.Render()
				cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
			}
			m.Input, cmd = m.Input.Update(msg)
		} else {
			switch {
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
		}
	case EditInfoMsg:
		cur := m.Items.Get(m.List.Model.SelectedItem())
		m.form = cur.Fields.Edit()
		//fields := f.Start()
		//cur.SetFields(fields)
	case HideInfoMsg:
		m.HideInfo()
		cmds = append(cmds, SetFocusedViewCmd("list"))
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
	case ShowItemInfoMsg:
		m.view = msg.Fields.Display()
		m.info = msg.Fields.Info()
		m.HideMenu()
		m.ShowInfo()
		cmds = append(cmds, SetFocusedViewCmd("info"))
	}

	switch focus {
	case "info":
		m.info, cmd = m.info.Update(msg)
		cmds = append(cmds, cmd)
	case "list":
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	default:
		for label, _ := range m.Menus {
			if focus == label {
				cmds = append(cmds, UpdateMenu(m, msg))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *TUI) View() string {
	var (
		sections    []string
		availHeight = m.List.Model.Height()
	)

	var menu string
	if m.showMenu {
		menu = m.CurrentMenu.Model.View()
		availHeight -= lipgloss.Height(menu)
	}

	var info string
	if m.showInfo {
		//m.info.Render()
		info = m.info.View()
		availHeight -= lipgloss.Height(info)
	}

	m.List.SetSize(m.width, availHeight)
	content := m.List.Model.View()
	sections = append(sections, content)

	if m.showMenu {
		sections = append(sections, menu)
	}

	if m.showInfo {
		sections = append(sections, info)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
