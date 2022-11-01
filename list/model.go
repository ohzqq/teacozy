package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/urbooks-core/bubbles/key"
	"github.com/ohzqq/urbooks-core/bubbles/style"
	"github.com/ohzqq/urbooks-core/bubbles/util"
	"golang.org/x/exp/slices"
)

type listType int

const (
	isPrompt listType = iota
	isSingle
	isMulti
	isSub
)

type List struct {
	Model            list.Model
	AllItems         Items
	Items            Items
	Selections       Items
	Keys             urkey.KeyMap
	area             textarea.Model
	Menus            Menus
	state            listType
	Title            string
	IsPrompt         bool
	IsMultiSelect    bool
	ShowSelectedOnly bool
	FocusedView      string
	width            int
	height           int
	ShowMenu         bool
	isFullScreen     bool
	showHelp         bool
	CurrentMenu      *Menu
	frame            lipgloss.Style
}

func New(title string) *List {
	l := &List{
		Title:       title,
		Keys:        urkey.DefaultKeys(),
		Menus:       make(Menus),
		FocusedView: "list",
	}
	l.frame = style.FrameStyle()
	return l
}

func (l *List) BuildModel() list.Model {
	l.processAllItems()
	items := l.DisplayItems("all")
	w := l.Width()
	h := l.GetHeight(items)
	del := NewItemDelegate(l.IsMulti())
	model := list.New(items, del, w, h)
	model.Title = l.Title
	model.Styles = style.ListStyles()
	model.SetShowStatusBar(false)
	model.SetShowHelp(false)

	return model
}

func (l List) GetAbsIndex(i list.Item) int {
	id := i.(Item).id
	fn := func(item list.Item) bool {
		return id == item.(Item).id
	}
	return slices.IndexFunc(l.AllItems, fn)
}

func (l List) GetSubList(i list.Item) Items {
	item := i.(Item)
	if item.HasList() {
		t := len(item.items)
		return l.AllItems[item.id+1 : item.id+t+1]
	}
	return Items{}
}

func (l *List) NewList(i Items) list.Model {
	list := &List{
		Title:       l.Title,
		Keys:        urkey.DefaultKeys(),
		Menus:       l.Menus,
		FocusedView: "list",
	}
	list.frame = style.FrameStyle()
	list.SetItems(i)
	list.processAllItems()

	return list.BuildModel()
}

func (l *List) NewMenu(label string, t key.Binding, keys []MenuItem) *Menu {
	cm := NewMenu(label, t).SetKeys(keys)
	cm.SetWidth(l.width)
	cm.BuildModel()
	l.Menus[label] = cm
	return cm
}

func (l *List) AddMenu(menu *Menu) {
	menu.SetWidth(l.width)
	menu.BuildModel()
	l.Menus[menu.Label] = menu
}

func (l List) GetHeight(items []list.Item) int {
	max := util.TermHeight()
	total := len(items)
	cur := l.Model.Height()

	switch {
	case l.isFullScreen:
		return max
	case cur > max:
		return max
	case total < max:
		return total + 6
	default:
		return max
	}
}

func (l List) Width() int {
	return util.TermWidth()
}

func (l *List) SetMulti() *List {
	l.IsMultiSelect = true
	return l
}

func (l *List) Prompt() *List {
	l.ShowSelectedOnly = true
	return l
}

func (m *List) SetItem(modelIndex int, item Item) {
	m.Model.SetItem(modelIndex, item)
	m.AllItems[item.id] = item
}

func (l *List) IsMulti() bool {
	return l.IsMultiSelect
}

func (l *List) SetShowHelp() *List {
	l.ShowMenu = true
	return l
}

func (l *List) SetItems(items Items) *List {
	l.Items = items
	return l
}

func (l *List) AppendItem(i ListItem) *List {
	item := NewListItem(i)
	if l.IsMulti() {
		item.isMulti = true
	}
	l.Items = append(l.Items, item)
	return l
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

	var field string
	if m.area.Focused() {
		field = m.area.View()
		availHeight -= lipgloss.Height(field)
	}

	m.Model.SetSize(m.width, availHeight)
	content := m.Model.View()
	sections = append(sections, content)

	if m.ShowMenu {
		sections = append(sections, menu)
	}

	if m.area.Focused() {
		sections = append(sections, field)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (l *List) Init() tea.Cmd {
	return UpdateDisplayedItemsCmd("all")
}
