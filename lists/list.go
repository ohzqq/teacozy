package lists

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cozykey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	Model            list.Model
	delegate         itemDelegate
	AllItems         []list.Item
	VisibleItems     []list.Item
	Items            Items
	area             textarea.Model
	input            textinput.Model
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

func NewSingleSelect(title string) List {
	m := List{
		IsMultiSelect: false,
		delegate:      NewItemDelegate(false),
		Items:         NewItems(),
	}
	m.Items.IsMultiSelect = false
	return m
}

func NewMultiSelect(title string) List {
	m := List{
		IsMultiSelect: true,
		delegate:      NewItemDelegate(true),
		Items:         NewItems(),
	}
	m.Items.IsMultiSelect = true
	return m
}

func (m *List) Build() {
	m.Model = list.New(m.VisibleItems, m.delegate, m.Width(), m.Height())
	m.Model.Title = m.Title
	m.Model.Styles = ListStyles()
	m.Model.SetShowStatusBar(false)
	m.Model.SetShowHelp(false)
}

func (m *List) Add(i Item) {
	m.appendItem(i)
	if i.HasList() {
		for _, item := range i.AllItems {
			li := item.(Item)
			li.IsHidden = true
			m.Add(li)
		}
	}
}

func (l List) HasSelections() bool {
	return len(l.Selections) > 0
}

func (l *List) appendItem(item Item) {
	item.Idx = len(l.AllItems)
	l.AllItems = append(l.AllItems, item)
}

func (l List) Width() int {
	return util.TermWidth()
}

func (l List) Height() int {
	return util.TermHeight()
}

func (m *List) SetItem(modelIndex int, item Item) {
	m.Model.SetItem(modelIndex, item)
	m.Items.All[item.Idx] = item
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

	var field string
	if m.area.Focused() {
		field = m.area.View()
		availHeight -= lipgloss.Height(field)
	}

	m.Model.SetSize(m.width, availHeight)
	content := m.Model.View()
	sections = append(sections, content)

	if m.focusWidget {
		sections = append(sections, menu)
	}

	if m.area.Focused() {
		sections = append(sections, field)
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
