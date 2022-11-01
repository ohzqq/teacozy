package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
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
	input := textarea.New()
	input.SetValue(i.Content)
	input.ShowLineNumbers = false
	return input
}

func (i Item) Toggle() key.Binding { return cozykey.EditField }
func (i Item) Label() string       { return i.label }

func (i Item) FilterValue() string {
	return i.Content
}

func (i *Item) ToggleSelected() {
	i.IsSelected = !i.IsSelected
}

func (i *Item) SetContent(content string) {
	i.Content = content
}

func (i Item) HasList() bool {
	has := len(i.Items.All) > 0
	return has
}

func (i Item) Prefix() string {
	if i.HasList() {
		if i.IsOpen {
			return closeSub
		}
		return openSub
	} else {
		if i.IsSelected {
			return check
		}
		return uncheck
	}

	return dash
}
