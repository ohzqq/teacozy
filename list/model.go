package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
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
	List
	//List             list.Model
	area             textarea.Model
	info             viewport.Model
	Selections       Items
	Keys             urkey.KeyMap
	Menus            Menus
	state            listType
	Title            string
	IsPrompt         bool
	IsMultiSelect    bool
	ShowSelectedOnly bool
	FocusedView      string
	width            int
	height           int
	ShowWidget       bool
	showMenu         bool
	showInfo         bool
	isFullScreen     bool
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
	return slices.IndexFunc(l.Items, fn)
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
	cur := l.List.Model.Height()

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

func (l *Model) ShowMenu() {
	l.showMenu = true
}

func (l *Model) HideMenu() {
	l.showMenu = false
}

func (l *Model) ShowInfo() {
	l.showInfo = true
}

func (l *Model) HideInfo() {
	l.showInfo = false
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

func (l *Model) SetShowHelp() *Model {
	l.showMenu = true
	return l
}

func (l *Model) SetItems(items Items) *Model {
	l.Items = items
	return l
}

func (m Model) View() string {
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
		info = m.info.View()
		availHeight -= lipgloss.Height(info)
	}

	var field string
	if m.area.Focused() {
		field = m.area.View()
		availHeight -= lipgloss.Height(field)
	}

	m.List.Model.SetSize(m.width, availHeight)
	content := m.List.Model.View()
	sections = append(sections, content)

	if m.showMenu {
		sections = append(sections, menu)
	}

	if m.showInfo {
		sections = append(sections, info)
	}

	if m.area.Focused() {
		sections = append(sections, field)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (l *Model) Init() tea.Cmd {
	return nil
	//return SetItemsCmd(l.Items)
}
