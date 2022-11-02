package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/slices"
)

type listType int

const (
	isPrompt listType = iota
	isSingle
	isMulti
	isSub
)

type Model struct {
	List             list.Model
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

func New(title string) *Model {
	l := &Model{
		Title:       title,
		Keys:        urkey.DefaultKeys(),
		Menus:       make(Menus),
		FocusedView: "list",
	}
	l.frame = style.FrameStyle()
	return l
}

func (l *Model) NewItem(content string) {
	item := NewItem(Item{Content: content})
	l.AppendItem(item)
}

func (l *Model) AppendItem(item Item) *Model {
	//item := NewItem(i)
	if l.IsMulti() {
		item.IsMulti = true
	}
	l.Items = l.Items.Add(item)
	return l
}

func (l *Model) BuildModel() list.Model {
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

func (l Model) GetAbsIndex(i list.Item) int {
	id := i.(Item).id
	fn := func(item list.Item) bool {
		return id == item.(Item).id
	}
	return slices.IndexFunc(l.AllItems, fn)
}

func (l *Model) NewList(i Items) list.Model {
	list := &Model{
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

func (l *Model) NewMenu(label string, t key.Binding, keys []MenuItem) *Menu {
	cm := NewMenu(label, t).SetKeys(keys)
	cm.SetWidth(l.width)
	cm.BuildModel()
	l.Menus[label] = cm
	return cm
}

func (l *Model) AddMenu(menu *Menu) {
	menu.SetWidth(l.width)
	menu.BuildModel()
	l.Menus[menu.Label] = menu
}

func (l Model) GetHeight(items []list.Item) int {
	max := util.TermHeight()
	total := len(items)
	cur := l.List.Height()

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

func (l Model) Width() int {
	return util.TermWidth()
}

func (l *Model) SetMulti() *Model {
	l.IsMultiSelect = true
	return l
}

func (l *Model) Prompt() *Model {
	l.ShowSelectedOnly = true
	return l
}

func (m *Model) SetItem(modelIndex int, item Item) {
	m.List.SetItem(modelIndex, item)
	m.Items[item.id] = item
}

func (l *Model) IsMulti() bool {
	return l.IsMultiSelect
}

func (l *Model) SetShowHelp() *Model {
	l.ShowMenu = true
	return l
}

func (l *Model) SetItems(items Items) *Model {
	l.Items = items
	return l
}

func (m Model) View() string {
	var (
		sections    []string
		availHeight = m.List.Height()
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

	m.List.SetSize(m.width, availHeight)
	content := m.List.View()
	sections = append(sections, content)

	if m.ShowMenu {
		sections = append(sections, menu)
	}

	if m.area.Focused() {
		sections = append(sections, field)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (l *Model) Init() tea.Cmd {
	return UpdateVisibleItemsCmd("all")
}
