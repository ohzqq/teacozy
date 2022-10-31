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
	Widgets          map[string]Widget
	focusWidget      bool
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
			for label, widget := range m.Widgets {
				if key.Matches(msg, widget.Toggle()) {
					widget.Focus()
					m.ShowWidget()
					cmds = append(cmds, SetFocusedViewCmd(label))
				}
			}
		}
	case ToggleItemMsg:
		cur := m.Model.SelectedItem().(Item)
		if m.IsMultiSelect {
			cur.IsSelected = !cur.IsSelected
		}
		m.SetItem(m.Model.Index(), cur)
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
		if m.FocusedView == "list" && m.CurrentWidget() != nil {
			m.HideWidget()
		}
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
		case UpdateWidgetContentMsg:
			m.CurrentWidget().SetContent(string(msg))
			m.HideWidget()
		}

		m.Model, cmd = m.Model.Update(msg)
		cmds = append(cmds, cmd)
	default:
		if m.CurrentWidget() != nil {
			cmds = append(cmds, m.CurrentWidget().Update(&m, msg))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *List) SetItem(modelIndex int, item Item) {
	m.Model.SetItem(modelIndex, item)
	m.Items.All[item.Idx] = item
}

func (m List) CurrentWidget() Widget {
	for _, w := range m.Widgets {
		if w.Focused() {
			return w
		}
	}
	return nil
}

func (m *List) HideWidget() {
	m.focusWidget = false
	if m.CurrentWidget() != nil {
		m.CurrentWidget().Blur()
	}
}

func (m *List) ShowWidget() {
	m.focusWidget = true
	if m.CurrentWidget() != nil {
		m.CurrentWidget().Focus()
	}
}

func (l *List) NewWidget(widget Widget) {
	l.Widgets[widget.Label()] = widget
}

func (m List) Init() tea.Cmd {
	return SetFocusedViewCmd("list")
}

func (m List) View() string {
	var (
		sections    []string
		availHeight = m.Model.Height()
	)

	var menu string
	if m.focusWidget {
		menu = m.CurrentWidget().View()
		availHeight -= lipgloss.Height(menu)
	}

	m.Model.SetSize(m.width, availHeight)
	content := m.Model.View()
	sections = append(sections, content)

	if m.focusWidget {
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
