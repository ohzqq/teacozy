package list

import "github.com/charmbracelet/bubbles/list"

type Items struct {
	All      []list.Item
	Selected map[int]list.Item
	Visible  []list.Item
}

func (i Items) NewList(title string) List {
	return New(title, i)
}

func (i *Items) Add(item Item) {
	i.appendItem(item)
	if item.HasList {
		for _, l := range item.Items.All {
			li := l.(Item)
			i.Add(li)
		}
	}
}

func (i *Items) appendItem(item Item) {
	item.idx = len(i.All)
	i.All = append(i.All, item)
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
	idx           int
	level         int
	Content       string
	Items         Items
	IsVisible     bool
	IsSelected    bool
	IsSub         bool
	HasList       bool
	ListIsOpen    bool
	IsMultiSelect bool
}

func NewItem(content string) Item {
	return Item{Content: content}
}

func (i Item) FilterValue() string {
	return i.Content
}

func (i Item) Prefix() string {
	return i.state.Prefix()
}
