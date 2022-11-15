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
	Main        *List
	Alt         *List
	Input       textarea.Model
	view        viewport.Model
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
		Main:        NewList(title, items),
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
				sel := m.Main.Model.SelectedItem()
				cur := m.Main.Items.Get(sel)
				field := cur.Item.(FieldData)
				val := m.Input.Value()
				if original := field.Value(); original != val {
					field.Set(val)
					item := NewItem().SetData(field)
					item.Changed = true
					m.Main.Items.Set(cur.Index(), item)
				}
				m.Input.Blur()
				cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			if m.Main.isForm {
				switch {
				case key.Matches(msg, Keys.ExitScreen):
					m.Main = m.Alt
				case key.Matches(msg, Keys.SaveAndExit):
					cur := m.Main.Model.SelectedItem()
					i := m.Main.Items.Get(cur)
					if i.Changed {
						cmds = append(cmds, ItemChangedCmd())
					}
					m.Main = m.Alt
					m.HideInfo()
				}
			}
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
	case ItemChangedMsg:
		sel := m.Main.Model.SelectedItem()
		cur := m.Main.Items.Get(sel)
		cur.Changed = true
	case EditFormItemMsg:
		if m.Main.isForm {
			m.Input = textarea.New()
			m.Input.SetValue(msg.Value())
			m.Input.ShowLineNumbers = false
			m.Input.Focus()
		}
	case EditInfoMsg:
		sel := m.Main.Model.SelectedItem()
		cur := m.Main.Items.Get(sel)
		m.Alt = m.Main
		m.Main = cur.Fields.Edit()
		m.HideInfo()
		cmds = append(cmds, SetFocusedViewCmd("list"))
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
	case "input":
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
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
	return m.Main.Init()
}

func (m *TUI) View() string {
	var (
		sections    []string
		availHeight = m.Main.Model.Height()
		field       string
	)

	//switch m.FocusedView {
	//case "form":
	if m.Input.Focused() {
		field = m.Input.View()
		availHeight -= lipgloss.Height(field)
	}

	//m.form.SetSize(m.width, availHeight)
	//content := m.form.View()
	//sections = append(sections, content)

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

	m.Main.SetSize(m.width, availHeight)
	content := m.Main.Model.View()
	sections = append(sections, content)

	if m.showMenu {
		sections = append(sections, menu)
	}

	if m.showInfo {
		sections = append(sections, info)
	}

	if m.Input.Focused() {
		sections = append(sections, field)
	}
	//default:
	//}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
