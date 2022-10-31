package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cozykey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	Model            list.Model
	Items            Items
	Selections       []Item
	Keys             cozykey.KeyMap
	Title            string
	ShowSelectedOnly bool
	FocusedView      string
	IsMultiSelect    bool
	width            int
	height           int
	ShowWidget       bool
	Widgets          map[string]Widget
	CurrentWidget    Widget
	ShowMenu         bool
	Menus            Menus
	CurrentMenu      Menu
	Action           ListAction
}

func New(title string, items Items, multi bool) List {
	m := List{
		Items:   items,
		Keys:    cozykey.DefaultKeys(),
		Widgets: make(map[string]Widget),
	}
	m.Model = list.New(items.All, NewItemDelegate(multi), m.Width(), m.Height())
	m.Model.Title = title
	m.Model.Styles = ListStyles()
	m.Model.SetShowStatusBar(false)
	m.Model.SetShowHelp(false)
	m.Menus = make(Menus)
	return m
}

func (l List) Width() int {
	return util.TermWidth()
}

func (l List) Height() int {
	return util.TermHeight()
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.IsMultiSelect {
		} else {
			switch {
			case key.Matches(msg, cozykey.Enter):
				cmds = append(cmds, m.Action(m))
			}
		}
		switch {
		case key.Matches(msg, m.Keys.ExitScreen):
			cmds = append(cmds, tea.Quit)
		default:
			//for label, menu := range m.Menus {
			//  if key.Matches(msg, menu.Toggle()) {
			//    m.CurrentMenu = menu
			//    m.ShowMenu = !m.ShowMenu
			//    cmds = append(cmds, SetFocusedViewCmd(label))
			//  }
			//}

		}
	case ToggleItemMsg:
		cur := m.Model.SelectedItem().(Item)
		if m.IsMultiSelect {
			cur.IsSelected = !cur.IsSelected
		}
		m.SetItem(m.Model.Index(), cur)
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
	case UpdateVisibleItemsMsg:
		switch string(msg) {
		case "selected":
		}
	}

	switch focus := m.FocusedView; focus {
	case "list":
		switch msg := msg.(type) {
		//case UpdateDisplayedItemsMsg:
		//items := m.DisplayItems(string(msg))
		//m.Model.SetHeight(m.GetHeight(items))
		//cmds = append(cmds, m.Model.SetItems(items))
		case UpdateMenuContentMsg:
			m.CurrentMenu.Model.SetContent(string(msg))
			m.ShowMenu = false
		}

		m.Model, cmd = m.Model.Update(msg)
		cmds = append(cmds, cmd)
	default:
		//for label, _ := range m.Menus {
		//  if focus == label {
		//    cmds = append(cmds, m.CurrentMenu.Update(&m, msg))
		//    //cmds = append(cmds, m.CurrentMenu.Update(&m, msg))
		//  }
		//}

	}

	return m, tea.Batch(cmds...)
}

func UpdateWidget(m *List, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.CurrentWidget.Toggle()):
			m.ShowWidget = false
			cmds = append(cmds, SetFocusedViewCmd("list"))
		default:
			for _, item := range m.CurrentMenu.Keys {
				if key.Matches(msg, item.Key) {
					cmds = append(cmds, item.Cmd(m))
					m.ShowWidget = false
				}
			}
			m.ShowWidget = false
			cmds = append(cmds, SetFocusedViewCmd("list"))
		}
	}
	m.CurrentMenu.Model, cmd = m.CurrentMenu.Model.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *List) SetItem(modelIndex int, item Item) {
	m.Model.SetItem(modelIndex, item)
	m.Items.All[item.Idx] = item
}

func (l *List) NewMenu(label string, t key.Binding, keys []MenuItem) Menu {
	cm := NewMenu(label, t)
	cm.SetKeys(keys)
	cm.SetWidth(l.width)
	cm.BuildModel()
	l.Menus[label] = cm
	//l.Widgets[label] = &cm
	return cm
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	var (
		sections    []string
		availHeight = m.Model.Height()
	)

	var menu string
	if m.ShowMenu {
		menu = m.CurrentMenu.Model.View()
		availHeight -= lipgloss.Height(menu)
	}

	m.Model.SetSize(m.width, availHeight)
	content := m.Model.View()
	sections = append(sections, content)

	if m.ShowMenu {
		sections = append(sections, menu)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

//func (l List) GetHeight(items []list.Item) int {
//  max := util.TermHeight()
//  total := len(items)
//  cur := l.Model.Height()

//  switch {
//  case l.isFullScreen:
//    return max
//  case cur > max:
//    return max
//  case total < max:
//    return total + 6
//  default:
//    return max
//  }
//}
