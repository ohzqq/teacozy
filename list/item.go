package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	cozykey "github.com/ohzqq/teacozy/key"
)

type Items struct {
	All         []list.Item
	Selected    map[int]list.Item
	Visible     []list.Item
	MultiSelect bool
}

func NewItems() Items {
	items := Items{
		Selected: make(map[int]list.Item),
	}
	return items
}

func (i Items) Get(idx int) Item {
	item := i.All[idx].(Item)
	return item
}

func (i Items) NewList(title string, multi bool) List {
	l := New(title, i, multi)
	return l
}

func (i *Items) Add(item Item) {
	i.appendItem(item)
	if item.HasList() {
		item.state = itemListClosed
		for _, l := range item.Items.All {
			li := l.(Item)
			i.Add(li)
		}
	}
}

func (i *Items) appendItem(item Item) {
	if i.MultiSelect {
		item.state = itemNotSelected
	}
	item.Idx = len(i.All)
	i.All = append(i.All, item)
}

func (i Items) HasSelections() bool {
	return i.Selected != nil && len(i.Selected) > 0
}

func (i *Items) ToggleSelected(idx int) list.Item {
	item := i.All[idx].(Item)
	item.IsSelected = !item.IsSelected
	i.All[idx] = item
	return i.All[idx]
}

func (i Items) Selections() []Item {
	var items []Item
	for idx, _ := range i.Selected {
		item := i.All[idx]
		items = append(items, item.(Item))
	}
	return items
}

type itemState int

const (
	itemNotSelected itemState = iota + 1
	itemSelected
	itemListOpen
	itemListClosed
	check    string = "[x] "
	uncheck  string = "[ ] "
	dash     string = "- "
	openSub  string = `[+] `
	closeSub string = `[-] `
)

func (s itemState) Prefix() string {
	switch s {
	case itemNotSelected:
		return uncheck
	case itemSelected:
		return check
	case itemListOpen:
		return closeSub
	case itemListClosed:
		return openSub
	default:
		return dash
	}
}

type Item struct {
	input         textarea.Model
	state         itemState
	Idx           int
	level         int
	label         string
	Content       string
	Items         Items
	IsVisible     bool
	IsSelected    bool
	IsOpen        bool
	IsSub         bool
	ListIsOpen    bool
	IsMultiSelect bool
}

func NewItem(content string) Item {
	return Item{
		Content: content,
	}
}

func (i *Item) Edit() textarea.Model {
	i.input = textarea.New()
	i.input.SetValue(i.Content)
	i.input.ShowLineNumbers = false
	return i.input
}

func (m Item) Update(list *List, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, cozykey.SaveAndExit) {
			m.SetContent(m.input.Value())
			m.Blur()
		}
		if m.Focused() {
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (i Item) Toggle() key.Binding { return cozykey.EditField }
func (i Item) Label() string       { return i.label }
func (i Item) Focused() bool       { return i.input.Focused() }

func (i Item) FilterValue() string {
	return i.Content
}

func (i *Item) Focus() tea.Cmd {
	i.input.Focus()
	return nil
}

func (i *Item) Blur() {
	i.input.Blur()
}

func (i *Item) SetContent(content string) {
	i.Content = content
	//i.input.SetValue(content)
}

func (m Item) View() string {
	return m.input.View()
}

func (i Item) HasList() bool {
	has := len(i.Items.All) > 0
	return has
}

func (i Item) Prefix() string {
	var state itemState
	if i.IsSelected {
		return itemSelected.Prefix()
	} else {
		return itemNotSelected.Prefix()
	}
	return state.Prefix()
}
